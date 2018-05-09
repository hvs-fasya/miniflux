package api

import (
	"errors"
	"net/http"

	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/json"
)

// CreateFilter is the API handler to create a new filter.
func (c *Controller) CreateFilter(w http.ResponseWriter, r *http.Request) {
	filter, err := decodeFilterPayload(r.Body)
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	ctx := context.New(r)
	userID := ctx.UserID()
	filter.UserID = userID
	if err := filter.ValidateFilterCreation(); err != nil {
		json.BadRequest(w, err)
		return
	}

	if c, err := c.store.FilterByTitle(userID, filter.FilterName); err != nil || c != nil {
		json.BadRequest(w, errors.New("This filter already exists"))
		return
	}

	err = c.store.CreateFilter(filter)
	if err != nil {
		json.ServerError(w, errors.New("Unable to create this filter"))
		return
	}

	json.Created(w, filter)
}

// GetFilters is the API handler to get a list of filters for a given user.
func (c *Controller) GetFilters(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	filters, err := c.store.Filters(ctx.UserID())
	if err != nil {
		json.ServerError(w, errors.New("Unable to fetch filters"))
		return
	}

	json.OK(w, filters)
}

// RemoveFilter is the API handler to remove a filter.
func (c *Controller) RemoveFilter(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	userID := ctx.UserID()
	filterID, err := request.IntParam(r, "filterID")
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	if !c.store.FilterExists(userID, filterID) {
		json.NotFound(w, errors.New("Filter not found"))
		return
	}

	if err := c.store.RemoveFilter(userID, filterID); err != nil {
		json.ServerError(w, errors.New("Unable to remove this filter"))
		return
	}

	json.NoContent(w)
}
