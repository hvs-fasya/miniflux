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

	//media tab
	mediaCategory, err := c.store.CategoryByTitleWOUserID(MediaNewsCategoryTitle)
	mediaCategoryID := mediaCategory.ID
	mediaBuilder := c.store.NewNewsEntryQueryBuilder()
	mediaBuilder.WithoutStatus(model.EntryStatusRemoved)
	mediaBuilder.WithOrder(model.DefaultSortingOrder)
	mediaBuilder.WithDirection(DefaultSortingDirection)
	mediaBuilder.WithOffset(offset)
	mediaBuilder.WithCategoryID(mediaCategoryID)
	mediaBuilder.WithLimit(limit)

	mediaStartDate := time.Now().AddDate(0, -1, 0)
	mediaBuilder.After(&mediaStartDate)

	mediaEntries, err := mediaBuilder.GetEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}
	mediaCount, err := mediaBuilder.CountEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}
	view.Set("mediaentries", mediaEntries)
	view.Set("mediatotal", mediaCount)
	view.Set("offset", offset)
	view.Set("limit", NewsEntriesLimit)

	var hasNext bool
	hasNext = (mediaCount - offset) > limit
	view.Set("hasNext", hasNext)

	html.OK(w, view.NewsAjaxRender("news_media"))
}
