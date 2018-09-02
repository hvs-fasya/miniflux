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

// NewSecurityQueryBuilder returns a new SecurityQueryBuilder
func (s *Storage) NewSecurityQueryBuilder() *SecurityQueryBuilder {
	return NewSecurityQueryBuilder(s)
}

// AlertExists checks if a alert exists by using the country_id.
func (s *Storage) AlertExists(countryID int64) bool {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:AlertExists] countryID=%d", countryID))

	var result int
	s.db.QueryRow(`SELECT count(*) as c FROM alerts WHERE country_id=$1`, countryID).Scan(&result)
	return result >= 1
}

// CreateAlert creates a new alert.
func (s *Storage) CreateAlert(alert *model.Alert) (err error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:CreateAlert] country_id=%d", alert.CountryID))

	query := `INSERT INTO alerts
		(last_updated, still_valid, latest_updates, risk_level, risk_details, country_id)
		VALUES
		($1, $2, $3, $4, $5, $6)`

	_, err = s.db.Exec(query, alert.LastUpdated, alert.StillValid, alert.LatestUpdates, alert.RiskLevel, alert.RiskDetails, alert.CountryID)
	if err != nil {
		return fmt.Errorf("unable to create alert: %v", err)
	}

	return nil
}

// UpdateAlert updates a Alert.
func (s *Storage) UpdateAlert(alert *model.Alert) error {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:UpdateAlert] alertID=%d", alert.ID))

	query := `UPDATE alerts SET
			last_updated=$1,
			still_valid=$2,
			latest_updates=$3,
			risk_level=$4,
			risk_details=$5
			WHERE country_id=$6`

	_, err := s.db.Exec(
		query,
		&alert.LastUpdated,
		&alert.StillValid,
		&alert.LatestUpdates,
		&alert.RiskLevel,
		&alert.RiskDetails,
		alert.CountryID,
	)
	if err != nil {
		return fmt.Errorf("unable to update alert: %v", err)
	}

	return nil
}

// AlertByCountryID finds a Alert by the countryID.
func (s *Storage) AlertByCountryID(countryID int64) (*model.Alert, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:AlertByCountryID] countryID=%d", countryID))
	query := `SELECT
		id, last_updated, still_valid, latest_updates, risk_level, risk_details, country_id
		FROM alerts
		WHERE country_id = $1`

	alert := model.NewAlert()
	err := s.db.QueryRow(query, countryID).Scan(
		&alert.ID,
		&alert.LastUpdated,
		&alert.StillValid,
		&alert.LatestUpdates,
		&alert.RiskLevel,
		&alert.RiskDetails,
		&alert.CountryID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to fetch alert: %v", err)
	}
	return alert, nil
}
