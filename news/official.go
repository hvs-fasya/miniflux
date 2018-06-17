package news

import (
	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/ui/session"
	"github.com/miniflux/miniflux/ui/view"
	"net/http"
	"time"
)

const DefaultCountry = "WORLDWIDE"

var (
	OfficialCategoriesExcluded = []string{
		"Media News",
		"Travel Alerts",
		"travel alerts",
		"Security Alerts",
	}
	CategoriesSuffixes = []string{
		"",
		" Tourism",
		" Immigration",
	}
)

// Official shows the Official template
func (c *Controller) Official(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)
	offset := request.QueryIntParam(r, "offset", 0)
	limit := request.QueryIntParam(r, "limit", NewsEntriesLimit)
	country := request.QueryParam(r, "country", DefaultCountry)

	//official tab
	var excludeCategoryIDs []int64
	for _, t := range OfficialCategoriesExcluded {
		c, _ := c.store.CategoryByTitleWOUserID(t)
		excludeCategoryIDs = append(excludeCategoryIDs, c.ID)
	}

	officialBuilder := c.store.NewNewsEntryQueryBuilder()
	officialBuilder.WithoutStatus(model.EntryStatusRemoved)
	officialBuilder.WithOrder(model.DefaultSortingOrder)
	officialBuilder.WithDirection(DefaultSortingDirection)
	var cat_ids []int64
	officialBuilder.WithOffset(offset)
	if country != DefaultCountry {

		for _, suff := range CategoriesSuffixes {
			if cat, err := c.store.CategoryByTitleWOUserID(country + suff); err == nil && cat != nil {
				logger.Debug("category id: ", cat.ID)
				cat_ids = append(cat_ids, cat.ID)
				officialBuilder.WithCategoryID(cat.ID)
			}
		}
	} else {
		officialBuilder.WithoutCategoryIDs(excludeCategoryIDs)
	}
	officialBuilder.WithLimit(limit)

	officialStartDate := time.Now().AddDate(0, -1, 0)
	officialBuilder.After(&officialStartDate)

	officialEntries, err := officialBuilder.GetEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}

	officialCount, err := officialBuilder.CountEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}
	if country != DefaultCountry && len(cat_ids) == 0 {
		officialCount = 0
		officialEntries = nil
	}

	view.Set("officialentries", officialEntries)
	view.Set("officialtotal", officialCount)
	view.Set("officialoffset", offset)
	view.Set("limit", NewsEntriesLimit)

	var hasNext bool
	hasNext = (officialCount - offset) > limit
	view.Set("officialHasNext", hasNext)

	html.OK(w, view.NewsAjaxRender("news_official"))
}
