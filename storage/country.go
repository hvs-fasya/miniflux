// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package storage

import (
	"fmt"
	"time"

	"database/sql"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/timer"
)

// SeedCountries seeds countries slice.
func (s *Storage) SeedCountries(countries model.Countries) (err error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:SeedCountries]"))
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	query := `INSERT INTO countries
		(name, alpha3)
		VALUES
		($1, $2)`

	for _, country := range countries {
		_, err = tx.Exec(query, country.Name, country.Alpha3)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// CreateCountry creates a new country.
func (s *Storage) CreateCountry(country *model.Country) (err error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:CreaCreateCountryteUser] name=%s", country.Name))

	query := `INSERT INTO countries
		(name, alpha3)
		VALUES
		($1, $2)`

	_, err = s.db.Exec(query, country.Name, country.Alpha3)
	if err != nil {
		return fmt.Errorf("unable to create country: %v", err)
	}

	return nil
}

// CountryByID finds a country by the ID.
func (s *Storage) CountryByID(countryID int64) (*model.Country, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:CountryByID] countryID=%d", countryID))
	query := `SELECT
		id, name, alpha3
		FROM countries
		WHERE id = $1`

	return s.fetchCountry(query, countryID)
}

// CountryByName finds a country by the name.
func (s *Storage) CountryByName(name string) (*model.Country, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:CountryByName] name=%s", name))
	query := `SELECT
		id, name, alpha3
		FROM countries
		WHERE name=LOWER($1)`

	return s.fetchCountry(query, name)
}

//// RemoveUser deletes a user.
//func (s *Storage) RemoveUser(userID int64) error {
//	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:RemoveUser] userID=%d", userID))
//
//	result, err := s.db.Exec("DELETE FROM users WHERE id = $1", userID)
//	if err != nil {
//		return fmt.Errorf("unable to remove this user: %v", err)
//	}
//
//	count, err := result.RowsAffected()
//	if err != nil {
//		return fmt.Errorf("unable to remove this user: %v", err)
//	}
//
//	if count == 0 {
//		return errors.New("nothing has been removed")
//	}
//
//	return nil
//}

func (s *Storage) fetchCountry(query string, args ...interface{}) (*model.Country, error) {
	country := model.NewCountry()
	err := s.db.QueryRow(query, args...).Scan(
		&country.ID,
		&country.Name,
		&country.Alpha3,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to fetch country: %v", err)
	}

	return country, nil
}

// Countries returns all Countries.
func (s *Storage) Countries() (model.Countries, error) {
	defer timer.ExecutionTime(time.Now(), "[Storage:Countries]")
	query := `
		SELECT
			id, name, alpha3
		FROM countries
		ORDER BY name ASC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch Countries: %v", err)
	}
	defer rows.Close()

	var countries model.Countries
	for rows.Next() {
		country := model.NewCountry()
		err := rows.Scan(
			&country.ID,
			&country.Name,
			&country.Alpha3,
		)

		if err != nil {
			return nil, fmt.Errorf("unable to fetch country row: %v", err)
		}

		countries = append(countries, country)
	}

	return countries, nil
}
