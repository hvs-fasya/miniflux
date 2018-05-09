// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui

import (
	"net/http"

	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/http/route"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/ui/form"
	"github.com/miniflux/miniflux/ui/session"
	"github.com/miniflux/miniflux/ui/view"
	"strconv"
	"time"
)

// ShowFeedFilters shows all entries for the given filter.
func (c *Controller) ShowFeedFilters(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	user, err := c.store.UserByID(ctx.UserID())
	if err != nil {
		html.ServerError(w, err)
		return
	}

	var filters []model.Filter
	filters, err = c.store.Filters(ctx.UserID())
	if err != nil {
		html.ServerError(w, err)
		return
	}

	var current model.Filter
	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)
	view.Set("menu", "filters")
	view.Set("user", user)
	view.Set("countUnread", c.store.CountUnreadEntries(user.ID))

	filterID, err := strconv.ParseInt(request.Param(r, "id", "0"), 10, 64)
	if filterID == 0 && len(filters) > 0 {
		current = filters[0]
	} else if filterID > 0 {
		found := false
		for _, f := range filters {
			if f.ID == filterID {
				current = f
				found = true
			}
		}
		if !found {
			view.Set("errorMessage", "Filter not found.")
		}
	}

	offset := request.QueryIntParam(r, "offset", 0)
	builder := c.store.NewEntryQueryBuilder(user.ID)
	builder.WithoutStatus(model.EntryStatusRemoved)
	builder.WithOrder(model.DefaultSortingOrder)
	builder.WithDirection(user.EntryDirection)
	builder.WithOffset(offset)
	builder.WithLimit(nbItemsPerPage)
	builder.WithFilter(current.Filters)

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

	view.Set("filters", filters)
	view.Set("current", current)
	view.Set("entries", entries)
	view.Set("total", count)

	//view.Set("pagination", c.getPagination(route.Path(c.router, "feedEntries", "feedID", feed.ID), count, offset))

	//view.Set("hasSaveEntry", c.store.HasSaveEntry(user.ID))

	html.OK(w, view.Render("feed_filters"))
}

// CreateFilter add new filter record.
func (c *Controller) CreateFilter(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	user, err := c.store.UserByID(ctx.UserID())
	if err != nil {
		html.ServerError(w, err)
		return
	}

	filterForm := form.NewFilterForm(r)

	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)
	view.Set("menu", "filters")
	view.Set("user", user)
	view.Set("countUnread", c.store.CountUnreadEntries(user.ID))

	if err := filterForm.Validate(); err != nil {
		view.Set("errorMessage", err.Error())
		html.OK(w, view.Render("feed_filters"))
		return
	}

	duplicateFilter, err := c.store.FilterByName(user.ID, filterForm.FilterName)
	if err != nil {
		html.ServerError(w, err)
		return
	}

	if duplicateFilter != nil {
		view.Set("errorMessage", "This filter already exists.")
		html.OK(w, view.Render("feed_filters"))
		return
	}

	filter := model.Filter{
		FilterName: filterForm.FilterName,
		UserID:     user.ID,
		Filters:    filterForm.Filters,
	}

	if err = c.store.CreateFilter(&filter); err != nil {
		logger.Error("[Controller:CreateFilter] %v", err)
		view.Set("errorMessage", "Unable to create this filter.")
		html.OK(w, view.Render("feed_filters"))
		return
	}

	filters, err := c.store.Filters(ctx.UserID())
	if err != nil {
		html.ServerError(w, err)
		return
	}

	current := filters[0]

	view.Set("filters", filters)
	view.Set("current", current)

	html.OK(w, view.Render("feed_filters"))
}

// RemoveFilter remove filter record.
func (c *Controller) RemoveFilter(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	user, err := c.store.UserByID(ctx.UserID())
	if err != nil {
		html.ServerError(w, err)
		return
	}

	filterID, err := request.IntParam(r, "filterID")
	if err != nil {
		html.BadRequest(w, err)
		return
	}

	filter, err := c.store.Filter(ctx.UserID(), filterID)
	if err != nil {
		html.ServerError(w, err)
		return
	}

	if filter == nil {
		html.NotFound(w)
		return
	}

	if err := c.store.RemoveFilter(user.ID, filter.ID); err != nil {
		html.ServerError(w, err)
		return
	}

	response.Redirect(w, r, route.Path(c.router, "filters"))
}
