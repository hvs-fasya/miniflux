package news

import (
	"net/http"
	"time"

	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/ui/session"
	"github.com/miniflux/miniflux/ui/view"
)

const (
	TravelNewsCategoryTitle = "Travel Alerts"
)

// Travel shows the Travel template
func (c *Controller) Travel(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)
	offset := request.QueryIntParam(r, "offset", 0)
	limit := request.QueryIntParam(r, "limit", NewsEntriesLimit)
	country := request.QueryParam(r, "country", DefaultCountry)

	//travel tab
	travelCategory, err := c.store.CategoryByTitleWOUserID(TravelNewsCategoryTitle)
	travelCategoryID := travelCategory.ID

	countryBuilder := c.store.NewNewsEntryQueryBuilder()
	countryBuilder.WithoutStatus(model.EntryStatusRemoved)
	countryBuilder.WithOrder(model.DefaultSortingOrder)
	countryBuilder.WithDirection(DefaultSortingDirection)

	countryBuilder.WithOffset(offset)
	countryBuilder.WithCategoryID(travelCategoryID)
	countryBuilder.WithLimit(limit)
	if country != DefaultCountry {
		countryBuilder.WithCountry(country)
	}

	countryStartDate := time.Now().AddDate(0, -3, 0)
	countryBuilder.After(&countryStartDate)

	countryEntries, err := countryBuilder.GetEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}
	countryCount, err := countryBuilder.CountEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}
	var allOffset int
	if offset == 0 {
		allOffset = 0
	} else {
		allOffset = offset - countryCount - 1
	}

	travelEntries := countryEntries
	var allCount int

	if country != DefaultCountry && len(countryEntries) < limit {
		travelBuilder := c.store.NewNewsEntryQueryBuilder()
		travelBuilder.WithLimit(limit - len(countryEntries))
		travelBuilder.WithCategoryID(travelCategoryID)
		travelBuilder.WithOffset(allOffset)
		travelBuilder.WithDirection(DefaultSortingDirection)
		travelBuilder.WithOrder(model.DefaultSortingOrder)
		travelBuilder.WithoutStatus(model.EntryStatusRemoved)
		travelBuilder.WithoutCountry(country)

		travelStartDate := time.Now().AddDate(0, -3, 0)
		travelBuilder.After(&travelStartDate)

		allEntries, err := travelBuilder.GetEntries()
		if err != nil {
			html.ServerError(w, err)
			return
		}
		travelEntries = append(travelEntries, allEntries...)
		allCount, err = travelBuilder.CountEntries()
		if err != nil {
			html.ServerError(w, err)
			return
		}
	}

	view.Set("travelentries", travelEntries)
	view.Set("countrytotal", len(countryEntries))
	view.Set("traveltotal", allCount+countryCount)
	view.Set("offset", offset)
	view.Set("limit", NewsEntriesLimit)

	var hasNext bool
	hasNext = (allCount + countryCount - offset) > limit
	view.Set("hasNext", hasNext)

	html.OK(w, view.NewsAjaxRender("news_travel"))
}
