package news

import (
	"net/http"

	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/response/html"
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

	html.OK(w, view.NewsRender("news_home"))
}