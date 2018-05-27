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
	"github.com/miniflux/miniflux/version"
)

// About shows the about page.
func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)
	view.Set("version", version.Version)
	view.Set("build_date", version.BuildDate)
	view.Set("menu", "settings")

	builder := c.store.NewNewsEntryQueryBuilder()
	offset := request.QueryIntParam(r, "offset", 0)
	builder.WithoutStatus(model.EntryStatusRemoved)
	builder.WithOrder(model.DefaultSortingOrder)
	builder.WithDirection(model.DefaultSortingDirection)
	builder.WithOffset(offset)
	builder.WithLimit(nbItemsPerPage)

	monthBefore := time.Now().AddDate(0, -1, 0)
	builder.After(&monthBefore)

	entries, err := builder.GetEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}

	count, err := builder.CountEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}

	//view.Set("current", current)
	view.Set("entries", entries)
	view.Set("total", count)

	html.OK(w, view.NewsRender("news_home"))
}
