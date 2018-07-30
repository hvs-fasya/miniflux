// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/miniflux/miniflux/http/response/json"
)

// GetCountries is the API handler to get a list of countries.
func (c *Controller) GetCountries(w http.ResponseWriter, r *http.Request) {
	countries, err := c.store.Countries()
	if err != nil {
		json.ServerError(w, errors.New("Unable to fetch countries"))
		return
	}

	json.OK(w, countries)
}
