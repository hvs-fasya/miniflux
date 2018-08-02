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
type HeadlinesQueryBuilder struct {
	store       *Storage
	categoryID  int64
	countryID   int64
	order       string
	direction   string
	headlineID  int64
	headlineIDs []int64
	before      *time.Time
	after       *time.Time
	visatype    string
}

// Before add condition base on the entry date.
func (e *HeadlinesQueryBuilder) Before(date *time.Time) *HeadlinesQueryBuilder {
	e.before = date
	return e
}

// After add condition base on the entry date.
func (e *HeadlinesQueryBuilder) After(date *time.Time) *HeadlinesQueryBuilder {
	e.after = date
	return e
}

// WithHeadlineIDs adds a condition to fetch only the given entry IDs.
func (e *HeadlinesQueryBuilder) WithHeadlineIDs(headlineIDs []int64) *HeadlinesQueryBuilder {
	e.headlineIDs = headlineIDs
	return e
}

// WithHeadlineID set the entryID.
func (e *HeadlinesQueryBuilder) WithHeadlineID(headlineID int64) *HeadlinesQueryBuilder {
	e.headlineID = headlineID
	return e
}

// WithCategoryID set the categoryID.
func (e *HeadlinesQueryBuilder) WithCategoryID(categoryID int64) *HeadlinesQueryBuilder {
	e.categoryID = categoryID
	return e
}

// WithCountryID set the country
func (e *HeadlinesQueryBuilder) WithCountryID(countryID int64) *HeadlinesQueryBuilder {
	e.countryID = countryID
	return e
}

// WithVisatype set the visatype
func (e *HeadlinesQueryBuilder) WithVisatype(visatype string) *HeadlinesQueryBuilder {
	e.visatype = visatype
	return e
}

// WithOrder set the sorting order.
func (e *HeadlinesQueryBuilder) WithOrder(order string) *HeadlinesQueryBuilder {
	e.order = order
	return e
}

// WithDirection set the sorting direction.
func (e *HeadlinesQueryBuilder) WithDirection(direction string) *HeadlinesQueryBuilder {
	e.direction = direction
	return e
}

// CountEntries count the number of entries that match the condition.
func (e *HeadlinesQueryBuilder) CountEntries() (count int, err error) {
	defer timer.ExecutionTime(
		time.Now(),
		fmt.Sprintf("[HeadlinesQueryBuilder:CountEntries]"),
	)

	query := `SELECT count(*) FROM headlines WHERE %s`
	args, condition := e.buildCondition()
	err = e.store.db.QueryRow(fmt.Sprintf(query, condition), args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("unable to count headlines: %v", err)
	}

	return count, nil
}

// GetHeadline returns a single entry that match the condition.
func (e *HeadlinesQueryBuilder) GetHeadline() (*model.Headline, error) {
	entries, err := e.GetHeadlines()
	if err != nil {
		return nil, err
	}

	if len(entries) != 1 {
		return nil, nil
	}

	return entries[0], nil
}

// GetHeadlines returns a list of entries that match the condition.
func (e *HeadlinesQueryBuilder) GetHeadlines() (model.Headlines, error) {
	debugStr := "[HeadlinesQueryBuilder:GetHeadlines] category_id=%d, country_id=%d, order=%s, direction=%s, visatype=%s"
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf(debugStr, e.categoryID, e.countryID, e.order, e.direction, e.visatype))

	query := `
		SELECT
		e.id, e.hash, e.published_at, e.title,
		e.url, e.content,
		e.category_id, e.country_id, e.visatype, e.icon_id,
		c.title as category_title,
		co.name as country_title
		FROM headlines e
		LEFT JOIN categories c ON c.id=e.category_id
		LEFT JOIN countries co ON co.id=e.country_id
		WHERE %s %s
	`

	args, conditions := e.buildCondition()
	query = fmt.Sprintf(query, conditions, e.buildSorting())
	fmt.Println(conditions)

	rows, err := e.store.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to get headlines: %v", err)
	}
	defer rows.Close()

	headlines := make(model.Headlines, 0)
	for rows.Next() {
		var headline model.Headline

		headline.Category = &model.Category{}
		headline.Country = &model.Country{}

		err := rows.Scan(
			&headline.ID,
			&headline.Hash,
			&headline.PublishedAt,
			&headline.Title,
			&headline.Url,
			&headline.Content,
			&headline.CategoryID,
			&headline.CountryID,
			&headline.VisaType,
			&headline.IconID,
			&headline.Category.Title,
			&headline.Country.Name,
		)

		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("unable to fetch headline row: %v", err)
		}

		headlines = append(headlines, &headline)
	}

	return headlines, nil
}

// GetHeadlineIDs returns a list of entry IDs that match the condition.
func (e *HeadlinesQueryBuilder) GetHeadlineIDs() ([]int64, error) {
	debugStr := "[HeadlinesQueryBuilder:GetHeadlineIDs] categoryID=%d,order=%s, direction=%s, visatype=%s"
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf(debugStr, e.categoryID, e.order, e.direction, e.visatype))

	query := `
		SELECT
		e.id
		FROM headlines e
		WHERE %s %s
	`

	args, conditions := e.buildCondition()
	query = fmt.Sprintf(query, conditions, e.buildSorting())

	rows, err := e.store.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to get headlines: %v", err)
	}
	defer rows.Close()

	var headlineIDs []int64
	for rows.Next() {
		var headlineID int64

		err := rows.Scan(&headlineID)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch headline row: %v", err)
		}

		headlineIDs = append(headlineIDs, headlineID)
	}

	return headlineIDs, nil
}

func (e *HeadlinesQueryBuilder) buildCondition() ([]interface{}, string) {
	args := []interface{}{}
	conditions := []string{}

	if e.categoryID != 0 {
		conditions = append(conditions, fmt.Sprintf("e.category_id=$%d", len(args)+1))
		args = append(args, e.categoryID)
	}

	if e.countryID != 0 {
		conditions = append(conditions, fmt.Sprintf("e.country_id=$%d", len(args)+1))
		args = append(args, e.countryID)
	}

	if e.headlineID != 0 {
		conditions = append(conditions, fmt.Sprintf("e.id=$%d", len(args)+1))
		args = append(args, e.headlineID)
	}

	if e.headlineIDs != nil {
		conditions = append(conditions, fmt.Sprintf("e.id=ANY($%d)", len(args)+1))
		args = append(args, pq.Array(e.headlineIDs))
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

func (e *HeadlinesQueryBuilder) buildSorting() string {
	var queries []string

	if e.order != "" {
		queries = append(queries, fmt.Sprintf(`ORDER BY "%s"`, e.order))
	}

	if e.direction != "" {
		queries = append(queries, fmt.Sprintf(`%s`, e.direction))
	}

	return strings.Join(queries, " ")
}

// NewHeadlinesQueryBuilder returns a new HeadlinesQueryBuilder.
func NewHeadlinesQueryBuilder(store *Storage) *HeadlinesQueryBuilder {
	return &HeadlinesQueryBuilder{
		store: store,
	}
}
