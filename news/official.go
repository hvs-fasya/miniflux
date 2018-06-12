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

// Official shows the Official template
func (c *Controller) Official(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)
	offset := request.QueryIntParam(r, "offset", 0)
	limit := request.QueryIntParam(r, "limit", NewsEntriesLimit)

	//official tab
	//officialCategory, err := c.store.CategoryByTitleWOUserID(MediaNewsCategoryTitle)
	//mediaCategoryID := mediaCategory.ID
	officialBuilder := c.store.NewNewsEntryQueryBuilder()
	officialBuilder.WithoutStatus(model.EntryStatusRemoved)
	officialBuilder.WithOrder(model.DefaultSortingOrder)
	officialBuilder.WithDirection(DefaultSortingDirection)
	officialBuilder.WithOffset(offset)
	//officialBuilder.WithCategoryID(mediaCategoryID)
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
	view.Set("officialentries", officialEntries)
	view.Set("officialtotal", officialCount)
	view.Set("officialoffset", offset)
	view.Set("limit", NewsEntriesLimit)

	var hasNext bool
	hasNext = (officialCount - offset) > limit
	view.Set("officialHasNext", hasNext)

	html.OK(w, view.NewsAjaxRender("news_official"))
}
