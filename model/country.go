// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model

// Country represents a country in the system.
type Country struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Alpha3      string `json:"alpha-3"`
	CountryCode string `json:"country-code"`
}

// NewCountry returns a new Country.
func NewCountry() *Country {
	return &Country{}
}

// Users represents a list of users.
type Countries []*Country
