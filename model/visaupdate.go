// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model

import (
	"errors"
	"time"
)

// Visaupdate represents a visaupdate in the system.
type Visaupdate struct {
	ID          int64     `json:"id"`
	Hash        string    `json:"-"`
	PublishedAt time.Time `json:"published_at"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Authority   string    `json:"authority"`
	CountryID   int64     `json:"country_id"`
	VisaType    string    `json:"visatype"`
	Country     *Country  `json:"country,omitempty"`
}

// NewVisaupdate returns a new Visaupdate.
func NewVisaupdate() *Visaupdate {
	return &Visaupdate{}
}

// Merge update the current Visaupdate with another Visaupdate.
func (h *Visaupdate) Merge(override *Visaupdate) {
	if override.Title != "" && h.Title != override.Title {
		h.Title = override.Title
	}

	if override.Content != "" && h.Content != override.Content {
		h.Content = override.Content
	}

	if override.Hash != "" && h.Hash != override.Hash {
		h.Hash = override.Hash
	}

	if override.Authority != "" && h.Authority != override.Authority {
		h.Authority = override.Authority
	}

	if override.VisaType != "" && h.VisaType != override.VisaType {
		h.VisaType = override.VisaType
	}

	if override.CountryID != 0 && h.CountryID != override.CountryID {
		h.CountryID = override.CountryID
	}
}

// Visaupdates represents a list of Visaupdates.
type Visaupdates []*Visaupdate

// ValidateVisaupdateCreation validates requirements.
func (h Visaupdate) ValidateVisaupdateCreation() error {
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

	return nil
}
