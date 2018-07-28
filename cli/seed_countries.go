// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/storage"
)

func seedCountries(store *storage.Storage) {
	var countries model.Countries
	file, err := ioutil.ReadFile("./news/3-codes-news.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_ = json.Unmarshal(file, &countries)
	if err := store.SeedCountries(countries); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
