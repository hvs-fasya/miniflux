package news // Copyright 2018 Frédéric Guillot. All rights reserved.

import (
	"errors"
	"net/http"

	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/json"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/reader/sanitizer"
	"github.com/miniflux/miniflux/reader/scraper"
)

// FetchContent downloads the original HTML page and returns relevant contents.
func (c *Controller) FetchContent(w http.ResponseWriter, r *http.Request) {
	entryID, err := request.IntParam(r, "entryID")
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	builder := c.store.NewNewsEntryQueryBuilder()
	builder.WithEntryID(entryID)
	builder.WithoutStatus(model.EntryStatusRemoved)

	entry, err := builder.GetEntry()
	if err != nil {
		json.ServerError(w, err)
		return
	}

	if entry == nil {
		json.NotFound(w, errors.New("Entry not found"))
		return
	}

	content, err := scraper.Fetch(entry.URL, entry.Feed.ScraperRules)
	if err != nil {
		json.ServerError(w, err)
		return
	}

	entry.Content = sanitizer.Sanitize(entry.URL, content)
	c.store.UpdateEntryContent(entry)

	json.Created(w, map[string]string{"content": entry.Content})
}
