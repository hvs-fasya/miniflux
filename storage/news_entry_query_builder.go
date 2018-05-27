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
	"github.com/miniflux/miniflux/timezone"
)

// NewsEntryQueryBuilder builds a SQL query to fetch entries.
type NewsEntryQueryBuilder struct {
	store              *Storage
	feedID             int64
	categoryID         int64
	status             string
	notStatus          string
	order              string
	direction          string
	limit              int
	offset             int
	entryID            int64
	greaterThanEntryID int64
	entryIDs           []int64
	before             *time.Time
	filter             string
	after              *time.Time
}

// Before add condition base on the entry date.
func (e *NewsEntryQueryBuilder) Before(date *time.Time) *NewsEntryQueryBuilder {
	e.before = date
	return e
}

// After add condition base on the entry date.
func (e *NewsEntryQueryBuilder) After(date *time.Time) *NewsEntryQueryBuilder {
	e.after = date
	return e
}

// WithGreaterThanEntryID adds a condition > entryID.
func (e *NewsEntryQueryBuilder) WithGreaterThanEntryID(entryID int64) *NewsEntryQueryBuilder {
	e.greaterThanEntryID = entryID
	return e
}

// WithEntryIDs adds a condition to fetch only the given entry IDs.
func (e *NewsEntryQueryBuilder) WithEntryIDs(entryIDs []int64) *NewsEntryQueryBuilder {
	e.entryIDs = entryIDs
	return e
}

// WithEntryID set the entryID.
func (e *NewsEntryQueryBuilder) WithEntryID(entryID int64) *NewsEntryQueryBuilder {
	e.entryID = entryID
	return e
}

// WithFeedID set the feedID.
func (e *NewsEntryQueryBuilder) WithFeedID(feedID int64) *NewsEntryQueryBuilder {
	e.feedID = feedID
	return e
}

// WithCategoryID set the categoryID.
func (e *NewsEntryQueryBuilder) WithCategoryID(categoryID int64) *NewsEntryQueryBuilder {
	e.categoryID = categoryID
	return e
}

// WithStatus set the entry status.
func (e *NewsEntryQueryBuilder) WithStatus(status string) *NewsEntryQueryBuilder {
	e.status = status
	return e
}

// WithoutStatus set the entry status that should not be returned.
func (e *NewsEntryQueryBuilder) WithoutStatus(status string) *NewsEntryQueryBuilder {
	e.notStatus = status
	return e
}

// WithOrder set the sorting order.
func (e *NewsEntryQueryBuilder) WithOrder(order string) *NewsEntryQueryBuilder {
	e.order = order
	return e
}

// WithDirection set the sorting direction.
func (e *NewsEntryQueryBuilder) WithDirection(direction string) *NewsEntryQueryBuilder {
	e.direction = direction
	return e
}

// WithLimit set the limit.
func (e *NewsEntryQueryBuilder) WithLimit(limit int) *NewsEntryQueryBuilder {
	e.limit = limit
	return e
}

// WithOffset set the offset.
func (e *NewsEntryQueryBuilder) WithOffset(offset int) *NewsEntryQueryBuilder {
	e.offset = offset
	return e
}

// WithFilter adds content text filter.
func (e *NewsEntryQueryBuilder) WithFilter(filter []string) *NewsEntryQueryBuilder {
	e.filter = strings.Join(filter, "|")
	return e
}

// CountEntries count the number of entries that match the condition.
func (e *NewsEntryQueryBuilder) CountEntries() (count int, err error) {
	defer timer.ExecutionTime(
		time.Now(),
		fmt.Sprintf("[NewsEntryQueryBuilder:CountEntries] feedID=%d, status=%s", e.feedID, e.status),
	)

	query := `SELECT count(*) FROM entries e LEFT JOIN feeds f ON f.id=e.feed_id WHERE %s`
	args, condition := e.buildCondition()
	err = e.store.db.QueryRow(fmt.Sprintf(query, condition), args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("unable to count entries: %v", err)
	}

	return count, nil
}

// GetEntry returns a single entry that match the condition.
func (e *NewsEntryQueryBuilder) GetEntry() (*model.Entry, error) {
	e.limit = 1
	entries, err := e.GetEntries()
	if err != nil {
		return nil, err
	}

	if len(entries) != 1 {
		return nil, nil
	}

	entries[0].Enclosures, err = e.store.GetEnclosures(entries[0].ID)
	if err != nil {
		return nil, err
	}

	return entries[0], nil
}

// GetEntries returns a list of entries that match the condition.
func (e *NewsEntryQueryBuilder) GetEntries() (model.Entries, error) {
	debugStr := "[NewsEntryQueryBuilder:GetEntries] feedID=%d, categoryID=%d, status=%s, order=%s, direction=%s, offset=%d, limit=%d"
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf(debugStr, e.feedID, e.categoryID, e.status, e.order, e.direction, e.offset, e.limit))

	query := `
		SELECT
		e.id, e.feed_id, e.hash, e.published_at, e.title,
		e.url, e.comments_url, e.author, e.content, e.status,
		f.title as feed_title, f.feed_url, f.site_url, f.checked_at,
		f.category_id, c.title as category_title, f.scraper_rules, f.rewrite_rules, f.crawler,
		fi.icon_id
		FROM entries e
		LEFT JOIN feeds f ON f.id=e.feed_id
		LEFT JOIN categories c ON c.id=f.category_id
		LEFT JOIN feed_icons fi ON fi.feed_id=f.id
		WHERE %s %s
	`

	args, conditions := e.buildCondition()
	query = fmt.Sprintf(query, conditions, e.buildSorting())

	rows, err := e.store.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to get entries: %v", err)
	}
	defer rows.Close()

	entries := make(model.Entries, 0)
	for rows.Next() {
		var entry model.Entry
		var iconID interface{}
		var tz string

		entry.Feed = &model.Feed{}
		entry.Feed.Category = &model.Category{}
		entry.Feed.Icon = &model.FeedIcon{}

		err := rows.Scan(
			&entry.ID,
			&entry.FeedID,
			&entry.Hash,
			&entry.Date,
			&entry.Title,
			&entry.URL,
			&entry.CommentsURL,
			&entry.Author,
			&entry.Content,
			&entry.Status,
			//&entry.Starred,
			&entry.Feed.Title,
			&entry.Feed.FeedURL,
			&entry.Feed.SiteURL,
			&entry.Feed.CheckedAt,
			&entry.Feed.Category.ID,
			&entry.Feed.Category.Title,
			&entry.Feed.ScraperRules,
			&entry.Feed.RewriteRules,
			&entry.Feed.Crawler,
			&iconID,
			//&tz,
		)

		if err != nil {
			return nil, fmt.Errorf("unable to fetch entry row: %v", err)
		}

		if iconID == nil {
			entry.Feed.Icon.IconID = 0
		} else {
			entry.Feed.Icon.IconID = iconID.(int64)
		}

		// Make sure that timestamp fields contains timezone information (API)
		entry.Date = timezone.Convert(tz, entry.Date)
		entry.Feed.CheckedAt = timezone.Convert(tz, entry.Feed.CheckedAt)

		entry.Feed.ID = entry.FeedID
		entry.Feed.Icon.FeedID = entry.FeedID
		entries = append(entries, &entry)
	}

	return entries, nil
}

// GetEntryIDs returns a list of entry IDs that match the condition.
func (e *NewsEntryQueryBuilder) GetEntryIDs() ([]int64, error) {
	debugStr := "[EntryQueryBuilder:GetEntryIDs] feedID=%d, categoryID=%d, status=%s, order=%s, direction=%s, offset=%d, limit=%d, filter=%s"
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf(debugStr, e.feedID, e.categoryID, e.status, e.order, e.direction, e.offset, e.limit, e.filter))

	query := `
		SELECT
		e.id
		FROM entries e
		LEFT JOIN feeds f ON f.id=e.feed_id
		WHERE %s %s
	`

	args, conditions := e.buildCondition()
	query = fmt.Sprintf(query, conditions, e.buildSorting())

	rows, err := e.store.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to get entries: %v", err)
	}
	defer rows.Close()

	var entryIDs []int64
	for rows.Next() {
		var entryID int64

		err := rows.Scan(&entryID)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch entry row: %v", err)
		}

		entryIDs = append(entryIDs, entryID)
	}

	return entryIDs, nil
}

func (e *NewsEntryQueryBuilder) buildCondition() ([]interface{}, string) {
	args := []interface{}{}
	conditions := []string{}

	if e.categoryID != 0 {
		conditions = append(conditions, fmt.Sprintf("f.category_id=$%d", len(args)+1))
		args = append(args, e.categoryID)
	}

	if e.feedID != 0 {
		conditions = append(conditions, fmt.Sprintf("e.feed_id=$%d", len(args)+1))
		args = append(args, e.feedID)
	}

	if e.entryID != 0 {
		conditions = append(conditions, fmt.Sprintf("e.id=$%d", len(args)+1))
		args = append(args, e.entryID)
	}

	if e.greaterThanEntryID != 0 {
		conditions = append(conditions, fmt.Sprintf("e.id > $%d", len(args)+1))
		args = append(args, e.greaterThanEntryID)
	}

	if e.entryIDs != nil {
		conditions = append(conditions, fmt.Sprintf("e.id=ANY($%d)", len(args)+1))
		args = append(args, pq.Array(e.entryIDs))
	}

	if e.status != "" {
		conditions = append(conditions, fmt.Sprintf("e.status=$%d", len(args)+1))
		args = append(args, e.status)
	}

	if e.notStatus != "" {
		conditions = append(conditions, fmt.Sprintf("e.status != $%d", len(args)+1))
		args = append(args, e.notStatus)
	}

	if e.before != nil {
		conditions = append(conditions, fmt.Sprintf("e.published_at < $%d", len(args)+1))
		args = append(args, e.before)
	}

	if e.filter != "" {
		conditions = append(conditions, fmt.Sprintf("(e.content SIMILAR TO '%%(' || $%d || ')%%' OR e.title SIMILAR TO '%%(' || $%d || ')%%')", len(args)+1, len(args)+1))
		args = append(args, e.filter)
	}

	if e.after != nil {
		conditions = append(conditions, fmt.Sprintf("e.published_at > $%d", len(args)+1))
		args = append(args, e.after)
	}

	return args, strings.Join(conditions, " AND ")
}

func (e *NewsEntryQueryBuilder) buildSorting() string {
	var queries []string

	if e.order != "" {
		queries = append(queries, fmt.Sprintf(`ORDER BY "%s"`, e.order))
	}

	if e.direction != "" {
		queries = append(queries, fmt.Sprintf(`%s`, e.direction))
	}

	if e.limit != 0 {
		queries = append(queries, fmt.Sprintf(`LIMIT %d`, e.limit))
	}

	if e.offset != 0 {
		queries = append(queries, fmt.Sprintf(`OFFSET %d`, e.offset))
	}

	return strings.Join(queries, " ")
}

// NewNewsEntryQueryBuilder returns a new NewsEntryQueryBuilder.
func NewNewsEntryQueryBuilder(store *Storage) *NewsEntryQueryBuilder {
	return &NewsEntryQueryBuilder{
		store: store,
	}
}
