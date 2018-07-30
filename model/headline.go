// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model

import (
	"database/sql"
	"errors"
	"time"
)

// User represents a user in the system.
type Headline struct {
	ID          int64         `json:"id"`
	Hash        string        `json:"-"`
	PublishedAt time.Time     `json:"published_at"`
	Title       string        `json:"title"`
	Content     string        `json:"content"`
	Url         string        `json:"url"`
	CountryID   int64         `json:"country_id"`
	VisaType    string        `json:"visatype"`
	CategoryID  int64         `json:"category_id"`
	IconID      sql.NullInt64 `json:"icon_id"`
}

// NewHeadline returns a new Headline.
func NewHeadline() *Headline {
	return &Headline{}
}

// Merge update the current Headline with another Headline.
func (h *Headline) Merge(override *Headline) {
	if override.Title != "" && h.Title != override.Title {
		h.Title = override.Title
	}

	if override.Content != "" && h.Content != override.Content {
		h.Content = override.Content
	}

	if override.Hash != "" && h.Hash != override.Hash {
		h.Hash = override.Hash
	}

	if override.Url != "" && h.Url != override.Url {
		h.Url = override.Url
	}

	if override.VisaType != "" && h.VisaType != override.VisaType {
		h.VisaType = override.VisaType
	}

	if override.CategoryID != 0 && h.CategoryID != override.CategoryID {
		h.CategoryID = override.CategoryID
	}

	if override.CountryID != 0 && h.CountryID != override.CountryID {
		h.CountryID = override.CountryID
	}
}

// Headlines represents a list of Headlines.
type Headlines []*Headline

// ValidateUserLogin validates user credential requirements.
func (h Headline) ValidateHeadlineCreation() error {
	if h.Title == "" {
		return errors.New("The title is mandatory")
	}

	if h.Content == "" {
		return errors.New("The content is mandatory")
	}

	if h.VisaType == "" {
		return errors.New("The visatype is mandatory")
	}

	if h.CountryID == 0 {
		return errors.New("The country_id is mandatory")
	}

	if h.CategoryID == 0 {
		return errors.New("The category_id is mandatory")
	}

	return nil
}
