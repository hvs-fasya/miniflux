// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/timer"
)

// HeadlinesQueryBuilder builds a SQL query to fetch entries.
type SecurityQueryBuilder struct {
	store     *Storage
	countryID int64
	after     *time.Time
}

// After add condition base on the entry date.
func (e *SecurityQueryBuilder) After(date *time.Time) *SecurityQueryBuilder {
	e.after = date
	return e
}

// WithCountryID set the country
func (e *SecurityQueryBuilder) WithCountryID(countryID int64) *SecurityQueryBuilder {
	e.countryID = countryID
	return e
}

// GetAlert returns a single entry that match the condition.
func (e *SecurityQueryBuilder) GetAlert() (*model.Alert, error) {
	alerts, err := e.GetSecurityAlerts()
	if err != nil {
		return nil, err
	}

	if len(alerts) != 1 {
		return nil, nil
	}

	return alerts[0], nil
}

//  returns a list of entries that match the condition.
func (e *SecurityQueryBuilder) GetSecurityAlerts() (model.Alerts, error) {
	debugStr := "[SecurityQueryBuilder:GetSecurityAlerts] country_id=%d"
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf(debugStr, e.countryID))

	query := `
		SELECT
		e.id, e.hash, e.published_at, e.title,
		e.url, e.content,
		e.category_id, e.country_id, e.visatype, e.icon_id,
		c.title as category_title,
		co.name as country_title
		FROM alerts e
		LEFT JOIN categories c ON c.id=e.category_id
		LEFT JOIN countries co ON co.id=e.country_id
		WHERE %s %s
	`

	args, conditions := e.buildCondition()
	query = fmt.Sprintf(query, conditions)

	rows, err := e.store.db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unable to get alerts: %v", err)
	}
	defer rows.Close()

	alerts := make(model.Alerts, 0)
	for rows.Next() {
		var alert model.Alert

		alert.Country = &model.Country{}
		alert.Health = &model.Health{}

		err := rows.Scan(
			&alert.ID,
			&alert.Hash,
			&alert.PublishedAt,
			&alert.Title,
			&alert.Url,
			&alert.Content,
			&alert.CategoryID,
			&alert.CountryID,
			&alert.VisaType,
			&alert.IconID,
			&alert.Category.Title,
			&alert.Country.Name,
		)

		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("unable to fetch alert row: %v", err)
		}

		alerts = append(alerts, &alert)
	}

	return alerts, nil
}

func (e *SecurityQueryBuilder) buildCondition() ([]interface{}, string) {
	args := []interface{}{}
	conditions := []string{}

	if e.countryID != 0 {
		conditions = append(conditions, fmt.Sprintf("e.country_id=$%d", len(args)+1))
		args = append(args, e.countryID)
	}

	return args, strings.Join(conditions, " AND ")
}

// NewSecurityQueryBuilder returns a new SecurityQueryBuilder.
func NewSecurityQueryBuilder(store *Storage) *SecurityQueryBuilder {
	return &SecurityQueryBuilder{
		store: store,
	}
}
