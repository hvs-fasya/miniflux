// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/json"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/news"
)

type UpdateOutput struct {
	ID           int64     `json:"id"`
	PublishedAt  time.Time `json:"published_at"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Url          string    `json:"url"`
	CountryName  string    `json:"country_name"`
	CategoryName string    `json:"category_name"`
	VisaType     string    `json:"visatype"`
	Icon         *feedIcon `json:"icon"`
}

// CreateVisaupdate is the API handler to create a new visaupdate.
func (c *Controller) CreateVisaupdate(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	if !ctx.IsAdminUser() {
		json.Forbidden(w)
		return
	}

	visaupdate, err := decodeVisaupdatePayload(r.Body)
	if err != nil {
		logger.Debug("visaupdate payload decode error: %s", err)
		json.BadRequest(w, err)
		return
	}

	if err := visaupdate.ValidateVisaupdateCreation(); err != nil {
		json.BadRequest(w, err)
		return
	}

	if c.store.VisaupdateExists(visaupdate.Title, visaupdate.Content) {
		json.BadRequest(w, errors.New("This visaupdate already exists"))
		return
	}

	if !c.store.CountryExists(visaupdate.CountryID) {
		json.BadRequest(w, errors.New("Country id does not exist"))
		return
	}

	err = c.store.CreateVisaupdate(visaupdate)
	if err != nil {
		logger.Debug("create visaupdate error: %s", err)
		json.ServerError(w, errors.New("Unable to create this visaupdate"))
		return
	}

	json.Created(w, visaupdate)
}

// UpdateVisaupdate is the API handler to update the given visaupdate.
func (c *Controller) UpdateVisaupdate(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	if !ctx.IsAdminUser() {
		json.Forbidden(w)
		return
	}

	visaupdateID, err := request.IntParam(r, "visaupdateID")
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	visaupdate, err := decodeVisaupdatePayload(r.Body)
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	originalVisaupdate, err := c.store.VisaupdateByID(visaupdateID)
	if err != nil {
		json.BadRequest(w, errors.New("Unable to fetch this visaupdate from the database"))
		return
	}

	if originalVisaupdate == nil {
		json.NotFound(w, errors.New("Visa update not found"))
		return
	}

	originalVisaupdate.Merge(visaupdate)
	if err = c.store.UpdateVisaupdate(originalVisaupdate); err != nil {
		json.ServerError(w, errors.New("Unable to update this visaupdate"))
		return
	}

	json.Created(w, originalVisaupdate)
}

// Visaupdates is the API handler to get the list of manual visaupdates.
func (c *Controller) Visaupdates(w http.ResponseWriter, r *http.Request) {
	visaType := request.QueryParam(r, "visatype", "")
	countryName := request.QueryParam(r, "country", "")
	from := request.QueryParam(r, "from", "")
	var country *model.Country
	var err error
	if countryName != "" {
		country, err = c.store.CountryByName(countryName)
		if err != nil {
			json.ServerError(w, errors.New("Unable to fetch country by name"))
			return
		}
		if country == nil {
			json.BadRequest(w, errors.New("Wrong country name"))
			return
		}
	}

	visaupdatesBuilder := c.store.NewVisaupdatesQueryBuilder()
	visaupdatesBuilder.WithOrder(model.DefaultSortingOrder)
	visaupdatesBuilder.WithDirection(news.DefaultSortingDirection)
	if country != nil {
		visaupdatesBuilder.WithCountryID(country.ID)
	}
	if visaType != "" {
		visaupdatesBuilder.WithVisatype(visaType)
	}

	var startDate time.Time
	if from != "" {
		startDate, _ = time.Parse("2006-01-02", from)
	} else {
		startDate = time.Now().AddDate(0, -1, 0)
	}
	visaupdatesBuilder.After(&startDate)

	visaupdates, err := visaupdatesBuilder.GetVisaupdates()
	if err != nil {
		json.ServerError(w, errors.New("Unable to fetch visaupdates"))
		return
	}

	var updates = []UpdateOutput{}
	for _, e := range visaupdates {
		v := UpdateOutput{
			PublishedAt: e.PublishedAt,
			Title:       e.Title,
			Content:     e.Content,
			Url:         e.Authority,
			CountryName: e.Country.Name,
			VisaType:    e.VisaType,
		}

		updates = append(updates, v)
	}

	json.OK(w, updates)
}

// RemoveVisaupdate is the API handler to remove an existing visaupdate.
func (c *Controller) RemoveVisaupdate(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	if !ctx.IsAdminUser() {
		json.Forbidden(w)
		return
	}

	visaupdateID, err := request.IntParam(r, "visaupdateID")
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	visaupdate, err := c.store.VisaupdateByID(visaupdateID)
	if err != nil {
		json.ServerError(w, errors.New("Unable to fetch this visaupdate from the database"))
		return
	}

	if visaupdate == nil {
		json.NotFound(w, errors.New("Visaupdate not found"))
		return
	}

	if err := c.store.RemoveVisaupdate(visaupdate.ID); err != nil {
		json.BadRequest(w, errors.New("Unable to remove this visaupdate from the database"))
		return
	}

	json.NoContent(w)
}
