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
	MediaNewsCategoryTitle = "Media News"
	NewsEntriesLimit       = 10
)

// Media shows the Media template
func (c *Controller) Media(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)
	offset := request.QueryIntParam(r, "offset", 0)
	limit := request.QueryIntParam(r, "limit", NewsEntriesLimit)
	country := request.QueryParam(r, "country", DefaultCountry)

	//media tab
	mediaCategory, err := c.store.CategoryByTitleWOUserID(MediaNewsCategoryTitle)
	mediaCategoryID := mediaCategory.ID

	countryBuilder := c.store.NewNewsEntryQueryBuilder()
	countryBuilder.WithoutStatus(model.EntryStatusRemoved)
	countryBuilder.WithOrder(model.DefaultSortingOrder)
	countryBuilder.WithDirection(DefaultSortingDirection)

	countryBuilder.WithOffset(offset)
	countryBuilder.WithCategoryID(mediaCategoryID)
	countryBuilder.WithLimit(limit)
	if country != DefaultCountry {
		countryBuilder.WithCountry(country)
	}

	countryStartDate := time.Now().AddDate(0, -1, 0)
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

	mediaEntries := countryEntries
	var allCount int

	if country != DefaultCountry && len(countryEntries) < limit {
		mediaBuilder := c.store.NewNewsEntryQueryBuilder()
		mediaBuilder.WithLimit(limit - len(countryEntries))
		mediaBuilder.WithCategoryID(mediaCategoryID)
		mediaBuilder.WithOffset(allOffset)
		mediaBuilder.WithDirection(DefaultSortingDirection)
		mediaBuilder.WithOrder(model.DefaultSortingOrder)
		mediaBuilder.WithoutStatus(model.EntryStatusRemoved)
		mediaBuilder.WithoutCountry(country)

		mediaStartDate := time.Now().AddDate(0, -1, 0)
		mediaBuilder.After(&mediaStartDate)

		allEntries, err := mediaBuilder.GetEntries()
		if err != nil {
			html.ServerError(w, err)
			return
		}
		mediaEntries = append(mediaEntries, allEntries...)
		allCount, err = mediaBuilder.CountEntries()
		if err != nil {
			html.ServerError(w, err)
			return
		}
	}

	view.Set("mediaentries", mediaEntries)
	view.Set("countrytotal", len(countryEntries))
	view.Set("mediatotal", allCount+countryCount)
	view.Set("offset", offset)
	view.Set("limit", NewsEntriesLimit)

	var hasNext bool
	hasNext = (allCount + countryCount - offset) > limit
	view.Set("hasNext", hasNext)

	html.OK(w, view.NewsAjaxRender("news_media"))
}
