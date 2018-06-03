package news

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/ui/session"
	"github.com/miniflux/miniflux/ui/view"
	"github.com/miniflux/miniflux/version"
)

const CountriesFile = "./news/country-json-as-per-Canada.json"

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Home shows the about page.
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

	view.Set("entries", entries)
	view.Set("total", count)

	countries := c.getCountries()
	view.Set("countries", countries)

	ticker := ""
	for _, e := range entries {
		ticker = ticker + e.Feed.Title + ": \"" + e.Title + "\" "
	}
	view.Set("ticker", ticker)

	html.OK(w, view.NewsRender("news_home"))
}

// Sources shows the Data Sources page.
func (c *Controller) Sources(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)

	countries := c.getCountries()
	view.Set("countries", countries)

	html.OK(w, view.NewsRender("news_sources"))
}

func (c *Controller) getCountries() []Country {
	var countries []Country
	content, err := ioutil.ReadFile(CountriesFile)
	if err != nil {
		logger.Debug("Error: %s", err)
		return countries
	}
	err = json.Unmarshal(content, &countries)
	if err != nil {
		logger.Debug("Error: %s", err)
		return countries
	}
	return countries
}
