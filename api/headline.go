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
	"github.com/miniflux/miniflux/reader/icon"
	"net/url"
	"strings"
)

var (
	keywordsOfficial = []string{
		"visa",
		"student",
		"work",
		"residence",
		"passport",
		"citizenship",
		"immigration",
		"startup",
		"business",
	}
)

type EntryOutput struct {
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Url          string    `json:"url"`
	PublishedAt  time.Time `json:"published_at"`
	CategoryName string    `json:"category_name"`
	Icon         *feedIcon `json:"icon"`
}

type HeadlineOutput struct {
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

// CreateHeadline is the API handler to create a new headline.
func (c *Controller) CreateHeadline(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	if !ctx.IsAdminUser() {
		json.Forbidden(w)
		return
	}

	headline, err := decodeHeadlinePayload(r.Body)
	if err != nil {
		logger.Debug("headline payload decode error: %s", err)
		json.BadRequest(w, err)
		return
	}

	if err := headline.ValidateHeadlineCreation(); err != nil {
		json.BadRequest(w, err)
		return
	}

	if c.store.HeadlineExists(headline.Title, headline.Content) {
		json.BadRequest(w, errors.New("This headline already exists"))
		return
	}

	if !c.store.CategoryWOUserIDExists(headline.CategoryID) {
		json.BadRequest(w, errors.New("Category id does not exist"))
		return
	}

	if !c.store.CountryExists(headline.CountryID) {
		json.BadRequest(w, errors.New("Country id does not exist"))
		return
	}

	err = c.store.CreateHeadline(headline)
	if err != nil {
		logger.Debug("create headline error: %s", err)
		json.ServerError(w, errors.New("Unable to create this headline"))
		return
	}

	hicon, err := icon.FindIcon(headline.Url)
	if err != nil {
		logger.Error("[Handler:CreateFeed] %v", err)
	} else if hicon == nil {
		logger.Info("No icon found for headlineID=%d", headline.ID)
	} else {
		c.store.CreateHeadlineIcon(headline, hicon)
	}

	json.Created(w, headline)
}

// UpdateHeadline is the API handler to update the given headline.
func (c *Controller) UpdateHeadline(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	if !ctx.IsAdminUser() {
		json.Forbidden(w)
		return
	}

	headlineID, err := request.IntParam(r, "headlineID")
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	headline, err := decodeHeadlinePayload(r.Body)
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	originalHeadline, err := c.store.HeadlineByID(headlineID)
	if err != nil {
		json.BadRequest(w, errors.New("Unable to fetch this headline from the database"))
		return
	}

	if originalHeadline == nil {
		json.NotFound(w, errors.New("Headline not found"))
		return
	}

	originalHeadline.Merge(headline)
	if err = c.store.UpdateHeadline(originalHeadline); err != nil {
		json.ServerError(w, errors.New("Unable to update this headline"))
		return
	}

	json.Created(w, originalHeadline)
}

// HeadlinesFull is the API handler to get the list of Headlines + Official filtered + Media.
func (c *Controller) HeadlinesFull(w http.ResponseWriter, r *http.Request) {
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

	//officials
	var excludeCategoryIDs []int64
	for _, t := range news.OfficialCategoriesExcluded {
		c, _ := c.store.CategoryByTitleWOUserID(t)
		excludeCategoryIDs = append(excludeCategoryIDs, c.ID)
	}

	officialBuilder := c.store.NewHeadlinesOfficialEntryQueryBuilder()
	officialBuilder.WithoutStatus(model.EntryStatusRemoved)
	officialBuilder.WithOrder(model.DefaultSortingOrder)
	officialBuilder.WithDirection(news.DefaultSortingDirection)
	officialBuilder.WithFilter(keywordsOfficial)
	var cat_ids []int64
	if countryName != "" {
		for _, suff := range news.CategoriesSuffixes {
			if cat, err := c.store.CategoryByTitleWOUserID(countryName + suff); err == nil && cat != nil {
				logger.Debug("category id: ", cat.ID)
				cat_ids = append(cat_ids, cat.ID)
				officialBuilder.WithCategoryID(cat.ID)
			}
		}
	} else {
		officialBuilder.WithoutCategoryIDs(excludeCategoryIDs)
	}

	var startDate time.Time
	if from != "" {
		startDate, _ = time.Parse("2006-01-02", from)
	} else {
		startDate = time.Now().AddDate(0, -1, 0)
	}
	officialBuilder.After(&startDate)

	officialEntries, err := officialBuilder.GetEntries()
	if err != nil {
		json.ServerError(w, errors.New("Unable to fetch official entries"))
		return
	}

	var officials = []EntryOutput{}
	if countryName != "" && len(cat_ids) == 0 {
		officialEntries = nil
	}

	for _, e := range officialEntries {
		var fIcon *model.Icon
		if c.store.HasIcon(e.Feed.ID) {
			fIcon, _ = c.store.IconByID(e.Feed.Icon.IconID)
		}

		var entryUrl = e.URL
		if strings.Contains(e.URL, `https://www.google.com/url?`) {
			u, _ := url.Parse(e.URL)
			entryUrl = u.Query().Get("url")
		}

		of := EntryOutput{
			Title:        e.Title,
			Content:      e.Content,
			Url:          entryUrl,
			PublishedAt:  e.Date,
			CategoryName: e.Feed.Category.Title,
		}
		if fIcon != nil {
			of.Icon = &feedIcon{
				ID:       fIcon.ID,
				MimeType: fIcon.MimeType,
				Data:     fIcon.DataURL(),
			}
		}
		officials = append(officials, of)
	}

	//media
	mediaCategory, err := c.store.CategoryByTitleWOUserID(news.MediaNewsCategoryTitle)
	mediaCategoryID := mediaCategory.ID

	mediaBuilder := c.store.NewNewsEntryQueryBuilder()
	mediaBuilder.WithoutStatus(model.EntryStatusRemoved)
	mediaBuilder.WithOrder(model.DefaultSortingOrder)
	mediaBuilder.WithDirection(news.DefaultSortingDirection)

	mediaBuilder.WithCategoryID(mediaCategoryID)
	if countryName != "" {
		mediaBuilder.WithCountry(countryName)
	}

	mediaBuilder.After(&startDate)

	mediaEntries, err := mediaBuilder.GetEntries()
	if err != nil {
		json.ServerError(w, errors.New("Unable to fetch media entries"))
		return
	}
	var medias = []EntryOutput{}
	for _, e := range mediaEntries {
		var fIcon *model.Icon
		if c.store.HasIcon(e.Feed.ID) {
			fIcon, _ = c.store.IconByID(e.Feed.Icon.IconID)
		}

		var entryUrl = e.URL
		if strings.Contains(e.URL, `https://www.google.com/url?`) {
			u, _ := url.Parse(e.URL)
			entryUrl = u.Query().Get("url")
		}

		of := EntryOutput{
			Title:        e.Title,
			Content:      e.Content,
			Url:          entryUrl,
			PublishedAt:  e.Date,
			CategoryName: e.Feed.Category.Title,
		}
		if fIcon != nil {
			of.Icon = &feedIcon{
				ID:       fIcon.ID,
				MimeType: fIcon.MimeType,
				Data:     fIcon.DataURL(),
			}
		}
		medias = append(medias, of)
	}

	//headlines
	headlinesBuilder := c.store.NewHeadlinesQueryBuilder()
	headlinesBuilder.WithOrder(model.DefaultSortingOrder)
	headlinesBuilder.WithDirection(news.DefaultSortingDirection)
	if country != nil {
		headlinesBuilder.WithCountryID(country.ID)
	}
	if visaType != "" {
		headlinesBuilder.WithVisatype(visaType)
	}

	headlinesBuilder.After(&startDate)

	headlines, err := headlinesBuilder.GetHeadlines()
	if err != nil {
		json.ServerError(w, errors.New("Unable to fetch headlines"))
		return
	}

	var hlines = []HeadlineOutput{}
	for _, e := range headlines {
		h := HeadlineOutput{
			PublishedAt:  e.PublishedAt,
			Title:        e.Title,
			Content:      e.Content,
			Url:          e.Url,
			CountryName:  e.Country.Name,
			CategoryName: e.Category.Title,
			VisaType:     e.VisaType,
		}
		var fIcon *model.Icon
		if e.IconID.Valid {
			fIcon, _ = c.store.IconByID(e.IconID.Int64)
		}
		if fIcon != nil {
			h.Icon = &feedIcon{
				ID:       fIcon.ID,
				MimeType: fIcon.MimeType,
				Data:     fIcon.DataURL(),
			}
		}
		hlines = append(hlines, h)
	}

	publishObj := struct {
		Officials      []EntryOutput    `json:"officials"`
		Headlines      []HeadlineOutput `json:"headlines"`
		Media          []EntryOutput    `json:"media"`
		OfficialCount  int              `json:"official_count"`
		MediaCount     int              `json:"media_count"`
		HeadlinesCount int              `json:"headlines_count"`
	}{
		officials,
		hlines,
		medias,
		len(officialEntries),
		len(mediaEntries),
		len(hlines),
	}

	json.OK(w, publishObj)
}

// Headlines is the API handler to get the list of Headlines.
func (c *Controller) Headlines(w http.ResponseWriter, r *http.Request) {
	headlines, err := c.store.Headlines()
	if err != nil {
		json.ServerError(w, errors.New("Unable to fetch the list of Headlines"))
		return
	}

	json.OK(w, headlines)
}

// HeadlineByID is the API handler to fetch the given Headline by the ID.
func (c *Controller) HeadlineByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	if !ctx.IsAdminUser() {
		json.Forbidden(w)
		return
	}

	headlineID, err := request.IntParam(r, "headlineID")
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	headline, err := c.store.HeadlineByID(headlineID)
	if err != nil {
		json.BadRequest(w, errors.New("Unable to fetch this headline from the database"))
		return
	}

	if headline == nil {
		json.NotFound(w, errors.New("Headline not found"))
		return
	}

	json.OK(w, headline)
}

// RemoveHeadline is the API handler to remove an existing headline.
func (c *Controller) RemoveHeadline(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	if !ctx.IsAdminUser() {
		json.Forbidden(w)
		return
	}

	headlineID, err := request.IntParam(r, "headlineID")
	if err != nil {
		json.BadRequest(w, err)
		return
	}

	headline, err := c.store.HeadlineByID(headlineID)
	if err != nil {
		json.ServerError(w, errors.New("Unable to fetch this headline from the database"))
		return
	}

	if headline == nil {
		json.NotFound(w, errors.New("Headline not found"))
		return
	}

	if err := c.store.RemoveHeadline(headline.ID); err != nil {
		json.BadRequest(w, errors.New("Unable to remove this headline from the database"))
		return
	}

	json.NoContent(w)
}
