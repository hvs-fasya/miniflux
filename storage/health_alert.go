// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package storage

import (
	"fmt"
	"time"

	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/timer"
)

//
//// NewHeadlinesQueryBuilder returns a new HeadlinesQueryBuilder
//func (s *Storage) NewHeadlinesQueryBuilder() *HeadlinesQueryBuilder {
//	return NewHeadlinesQueryBuilder(s)
//}

// CreateHealthAlert creates a new Health-Alert record.
func (s *Storage) CreateHealthAlert(ha *model.HealthAlert) (err error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:CreateHealthAlert] health_title=%s country_name=%s", ha.Health.HealthTitle, ha.Country.Name))

	query := `INSERT INTO alert_health
		(country_id, health_id, alert_health_date)
		VALUES
		($1, $2, $3)`

	_, err = s.db.Exec(query, ha.Country.ID, ha.Health.ID, ha.Date)
	if err != nil {
		return fmt.Errorf("unable to create health_alert: %v", err)
	}

	return nil
}

// UpdateHealthAlert updates a HealthAlert record.
func (s *Storage) UpdateHealthAlert(ha *model.HealthAlert) error {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:UpdateHealthAlert] health_title=%s", ha.Health.HealthTitle))

	query := `UPDATE alert_health SET
			date=$1
			WHERE health_id=$2 AND country_id=$3`

	_, err := s.db.Exec(
		query,
		&ha.Date,
		ha.Health.ID,
		ha.Country.ID,
	)
	if err != nil {
		return fmt.Errorf("unable to update health_alert: %v", err)
	}

	return nil
}

// HealthAlertsByCountryName finds a Health by country_name.
func (s *Storage) HealthAlertsByCountryName(countryName string) (model.HealthAlerts, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:HealthAlertsByCountryName] countryName=%s", countryName))
	query := `SELECT
		ah.health_id, ah.country_id, ah.alert_health_date,
		c.name, c.alpha3
		FROM alert_health ah
		JOIN countries c ON c.id=ah.country_id
		WHERE c.name = $1`

	rows, err := s.db.Query(query, countryName)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch health_alerts: %v", err)
	}
	defer rows.Close()

	healthAlerts := make(model.HealthAlerts, 0)
	for rows.Next() {
		var healthAlert model.HealthAlert
		if err := rows.Scan(&healthAlert.Health.ID, &healthAlert.Country.ID, &healthAlert.Date, &healthAlert.Country.Name, &healthAlert.Country.Alpha3); err != nil {
			return nil, fmt.Errorf("Unable to fetch health_alert row: %v", err)
		}
		healthAlerts = append(healthAlerts, &healthAlert)
	}

	return healthAlerts, nil
}

// HealthAlertsByCountryNameWithHealth finds a HealthAlert with health data by country_name.
func (s *Storage) HealthAlertsByCountryNameWithHealth(countryName string) (model.HealthAlerts, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:HealthAlertsByCountryNameWithHealth] countryName=%s", countryName))
	query := `SELECT
		ah.health_id, ah.country_id, ah.alert_health_date,
		c.name, c.alpha3,
		h.health_title, h.health_link, h.health_content
		FROM alert_health ah
		JOIN countries c ON c.id=ah.country_id
		JOIN healths h ON ah.health_id=h.id
		WHERE c.name = $1`

	rows, err := s.db.Query(query, countryName)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch health_alerts with health: %v", err)
	}
	defer rows.Close()

	healthAlerts := make(model.HealthAlerts, 0)
	for rows.Next() {
		var ha model.HealthAlert
		if err := rows.Scan(&ha.Health.ID, &ha.Country.ID, &ha.Date, &ha.Country.Name, &ha.Country.Alpha3, &ha.Health.HealthTitle, &ha.Health.HealthLink, &ha.Health.HealthContent); err != nil {
			return nil, fmt.Errorf("Unable to fetch health_alert with health row: %v", err)
		}
		healthAlerts = append(healthAlerts, &ha)
	}

	return healthAlerts, nil
}

// ClearHealthAlertsForCountry removes HealthAlert records for country.
func (s *Storage) ClearHealthAlertsForCountry(country *model.Country) (err error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:ClearHealthAlertsForCountry] for country_name=%s", country.Name))

	query := `DELETE FROM alert_health ah
		WHERE country_id IN 
		(SELECT id FROM countries WHERE name = $1)`

	_, err = s.db.Exec(query, country.Name)
	if err != nil {
		return fmt.Errorf("unable to clear health_alert records: %v", err)
	}

	return nil
}
