// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/timer"
)

// HeadlinesQueryBuilder builds a SQL query to fetch entries.
type VisaupdatesQueryBuilder struct {
	store         *Storage
	countryID     int64
	order         string
	direction     string
	visaupdateID  int64
	visaupdateIDs []int64
	before        *time.Time
	after         *time.Time
	visatype      string
}

// Before add condition base on the entry date.
func (e *VisaupdatesQueryBuilder) Before(date *time.Time) *VisaupdatesQueryBuilder {
	e.before = date
	return e
}

// After add condition base on the entry date.
func (e *VisaupdatesQueryBuilder) After(date *time.Time) *VisaupdatesQueryBuilder {
	e.after = date
	return e
}

// WithHeadlineIDs adds a condition to fetch only the given entry IDs.
func (e *VisaupdatesQueryBuilder) WithHeadlineIDs(visaupdateIDs []int64) *VisaupdatesQueryBuilder {
	e.visaupdateIDs = visaupdateIDs
	return e
}

// WithHeadlineID set the entryID.
func (e *VisaupdatesQueryBuilder) WithHeadlineID(visaupdateID int64) *VisaupdatesQueryBuilder {
	e.visaupdateID = visaupdateID
	return e
}

// WithCountryID set the country
func (e *VisaupdatesQueryBuilder) WithCountryID(countryID int64) *VisaupdatesQueryBuilder {
	e.countryID = countryID
	return e
}

// WithVisatype set the visatype
func (e *VisaupdatesQueryBuilder) WithVisatype(visatype string) *VisaupdatesQueryBuilder {
	e.visatype = visatype
	return e
}

// WithOrder set the sorting order.
func (e *VisaupdatesQueryBuilder) WithOrder(order string) *VisaupdatesQueryBuilder {
	e.order = order
	return e
}

// WithDirection set the sorting direction.
func (e *VisaupdatesQueryBuilder) WithDirection(direction string) *VisaupdatesQueryBuilder {
	e.direction = direction
	return e
}

// CountEntries count the number of entries that match the condition.
func (e *VisaupdatesQueryBuilder) CountEntries() (count int, err error) {
	defer timer.ExecutionTime(
		time.Now(),
		fmt.Sprintf("[VisaupdatesQueryBuilder:CountEntries]"),
	)

	query := `SELECT count(*) FROM visaupdates WHERE %s`
	args, condition := e.buildCondition()
	err = e.store.db.QueryRow(fmt.Sprintf(query, condition), args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("unable to count visaupdates: %v", err)
	}

	return count, nil
}

// GetVisaupdate returns a single entry that match the condition.
func (e *VisaupdatesQueryBuilder) GetVisaupdate() (*model.Visaupdate, error) {
	entries, err := e.GetVisaupdates()
	if err != nil {
		return nil, err
	}

	if len(entries) != 1 {
		return nil, nil
	}

	return entries[0], nil
}

// GetHeadlines returns a list of entries that match the condition.
func (e *VisaupdatesQueryBuilder) GetVisaupdates() (model.Visaupdates, error) {
	debugStr := "[VisaupdatesQueryBuilder:GetVisaupdates] country_id=%d, order=%s, direction=%s, visatype=%s"
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf(debugStr, e.countryID, e.order, e.direction, e.visatype))

	query := `
		SELECT
		e.id, e.hash, e.published_at, e.title,
		e.authority, e.content,
		e.country_id, e.visatype,
		co.name as country_title
		FROM visaupdates e
		LEFT JOIN countries co ON co.id=e.country_id
		WHERE %s %s
	`

	args, conditions := e.buildCondition()
	query = fmt.Sprintf(query, conditions, e.buildSorting())

	rows, err := e.store.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to get visaupdates: %v", err)
	}
	defer rows.Close()

	visaupdates := make(model.Visaupdates, 0)
	for rows.Next() {
		var visaupdate model.Visaupdate

		visaupdate.Country = &model.Country{}

		err := rows.Scan(
			&visaupdate.ID,
			&visaupdate.Hash,
			&visaupdate.PublishedAt,
			&visaupdate.Title,
			&visaupdate.Authority,
			&visaupdate.Content,
			&visaupdate.CountryID,
			&visaupdate.VisaType,
			&visaupdate.Country.Name,
		)

		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("unable to fetch visaupdate row: %v", err)
		}

		visaupdates = append(visaupdates, &visaupdate)
	}

	return visaupdates, nil
}

// GetVisaupdateIDs returns a list of entry IDs that match the condition.
func (e *VisaupdatesQueryBuilder) GetVisaupdateIDs() ([]int64, error) {
	debugStr := "[HeadlinesQueryBuilder:GetVisaupdateIDs] order=%s, direction=%s, visatype=%s"
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf(debugStr, e.order, e.direction, e.visatype))

	query := `
		SELECT
		e.id
		FROM visaupdates e
		WHERE %s %s
	`

	args, conditions := e.buildCondition()
	query = fmt.Sprintf(query, conditions, e.buildSorting())

	rows, err := e.store.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to get visaupdates: %v", err)
	}
	defer rows.Close()

	var visaupdateIDs []int64
	for rows.Next() {
		var visaupdateID int64

		err := rows.Scan(&visaupdateID)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch visaupdate row: %v", err)
		}

		visaupdateIDs = append(visaupdateIDs, visaupdateID)
	}

	return visaupdateIDs, nil
}

func (e *VisaupdatesQueryBuilder) buildCondition() ([]interface{}, string) {
	args := []interface{}{}
	conditions := []string{}

	if e.countryID != 0 {
		conditions = append(conditions, fmt.Sprintf("e.country_id=$%d", len(args)+1))
		args = append(args, e.countryID)
	}

	if e.visaupdateID != 0 {
		conditions = append(conditions, fmt.Sprintf("e.id=$%d", len(args)+1))
		args = append(args, e.visaupdateID)
	}

	if e.visaupdateIDs != nil {
		conditions = append(conditions, fmt.Sprintf("e.id=ANY($%d)", len(args)+1))
		args = append(args, pq.Array(e.visaupdateIDs))
	}

	if e.before != nil {
		conditions = append(conditions, fmt.Sprintf("e.published_at < $%d", len(args)+1))
		args = append(args, e.before)
	}

	if e.after != nil {
		conditions = append(conditions, fmt.Sprintf("e.published_at > $%d", len(args)+1))
		args = append(args, e.after)
	}

	if e.visatype != "" {
		conditions = append(conditions, fmt.Sprintf("e.visatype=$%d", len(args)+1))
		args = append(args, e.visatype)
	}

	return args, strings.Join(conditions, " AND ")
}

func (e *VisaupdatesQueryBuilder) buildSorting() string {
	var queries []string

	if e.order != "" {
		queries = append(queries, fmt.Sprintf(`ORDER BY "%s"`, e.order))
	}

	if e.direction != "" {
		queries = append(queries, fmt.Sprintf(`%s`, e.direction))
	}

	return strings.Join(queries, " ")
}

// NewVisaupdatesQueryBuilder returns a new VisaupdatesQueryBuilder.
func NewVisaupdatesQueryBuilder(store *Storage) *VisaupdatesQueryBuilder {
	return &VisaupdatesQueryBuilder{
		store: store,
	}
}
