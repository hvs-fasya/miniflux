// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/timer"
	"github.com/miniflux/miniflux/timezone"
)

// NewsEntryQueryBuilder builds a SQL query to fetch entries.
type HeadlinesOfficialEntryQueryBuilder struct {
	store              *Storage
	feedID             int64
	categoryID         int64
	notCategoryIDs     []int64
	country            string
	nocountry          string
	status             string
	notStatus          string
	order              string
	direction          string
	entryID            int64
	greaterThanEntryID int64
	entryIDs           []int64
	before             *time.Time
	filter             string
	after              *time.Time
}

// Before add condition base on the entry date.
func (e *HeadlinesOfficialEntryQueryBuilder) Before(date *time.Time) *HeadlinesOfficialEntryQueryBuilder {
	e.before = date
	return e
}

// After add condition base on the entry date.
func (e *HeadlinesOfficialEntryQueryBuilder) After(date *time.Time) *HeadlinesOfficialEntryQueryBuilder {
	e.after = date
	return e
}

// WithGreaterThanEntryID adds a condition > entryID.
func (e *HeadlinesOfficialEntryQueryBuilder) WithGreaterThanEntryID(entryID int64) *HeadlinesOfficialEntryQueryBuilder {
	e.greaterThanEntryID = entryID
	return e
}

// WithEntryIDs adds a condition to fetch only the given entry IDs.
func (e *HeadlinesOfficialEntryQueryBuilder) WithEntryIDs(entryIDs []int64) *HeadlinesOfficialEntryQueryBuilder {
	e.entryIDs = entryIDs
	return e
}

// WithEntryID set the entryID.
func (e *HeadlinesOfficialEntryQueryBuilder) WithEntryID(entryID int64) *HeadlinesOfficialEntryQueryBuilder {
	e.entryID = entryID
	return e
}

// WithFeedID set the feedID.
func (e *HeadlinesOfficialEntryQueryBuilder) WithFeedID(feedID int64) *HeadlinesOfficialEntryQueryBuilder {
	e.feedID = feedID
	return e
}

// WithCategoryID set the categoryID.
func (e *HeadlinesOfficialEntryQueryBuilder) WithCategoryID(categoryID int64) *HeadlinesOfficialEntryQueryBuilder {
	e.categoryID = categoryID
	return e
}

// WithoutCategoryIDs set the categoryIDs excluded.
func (e *HeadlinesOfficialEntryQueryBuilder) WithoutCategoryIDs(categoryIDs []int64) *HeadlinesOfficialEntryQueryBuilder {
	e.notCategoryIDs = categoryIDs
	return e
}

// WithCountry set the country
func (e *HeadlinesOfficialEntryQueryBuilder) WithCountry(country string) *HeadlinesOfficialEntryQueryBuilder {
	e.country = country
	return e
}

// WithoutCountry set the country
func (e *HeadlinesOfficialEntryQueryBuilder) WithoutCountry(country string) *HeadlinesOfficialEntryQueryBuilder {
	e.nocountry = country
	return e
}

// WithStatus set the entry status.
func (e *HeadlinesOfficialEntryQueryBuilder) WithStatus(status string) *HeadlinesOfficialEntryQueryBuilder {
	e.status = status
	return e
}

// WithoutStatus set the entry status that should not be returned.
func (e *HeadlinesOfficialEntryQueryBuilder) WithoutStatus(status string) *HeadlinesOfficialEntryQueryBuilder {
	e.notStatus = status
	return e
}

// WithOrder set the sorting order.
func (e *HeadlinesOfficialEntryQueryBuilder) WithOrder(order string) *HeadlinesOfficialEntryQueryBuilder {
	e.order = order
	return e
}

// WithDirection set the sorting direction.
func (e *HeadlinesOfficialEntryQueryBuilder) WithDirection(direction string) *HeadlinesOfficialEntryQueryBuilder {
	e.direction = direction
	return e
}

// WithFilter adds content text filter.
func (e *HeadlinesOfficialEntryQueryBuilder) WithFilter(filter []string) *HeadlinesOfficialEntryQueryBuilder {
	e.filter = strings.Join(filter, "|")
	return e
}

// CountEntries count the number of entries that match the condition.
func (e *HeadlinesOfficialEntryQueryBuilder) CountEntries() (count int, err error) {
	defer timer.ExecutionTime(
		time.Now(),
		fmt.Sprintf("[HeadlinesOfficialEntryQueryBuilder:CountEntries] feedID=%d, status=%s", e.feedID, e.status),
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
func (e *HeadlinesOfficialEntryQueryBuilder) GetEntry() (*model.Entry, error) {
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
func (e *HeadlinesOfficialEntryQueryBuilder) GetEntries() (model.Entries, error) {
	debugStr := "[HeadlinesOfficialEntryQueryBuilder:GetEntries] feedID=%d, categoryID=%d, status=%s, order=%s, direction=%s, filter=%s"
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf(debugStr, e.feedID, e.categoryID, e.status, e.order, e.direction, e.filter))

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

// GetEntries returns a list of entries that match the condition.
func (e *HeadlinesOfficialEntryQueryBuilder) GetEntriesWithIcons() ([]*model.EntryWithIcon, error) {
	debugStr := "[HeadlinesOfficialEntryQueryBuilder:GetEntriesWithIcons] feedID=%d, categoryID=%d, status=%s, order=%s, direction=%s, filter=%s"
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf(debugStr, e.feedID, e.categoryID, e.status, e.order, e.direction, e.filter))

	query := `
		SELECT
		e.id, e.feed_id, e.hash, e.published_at, e.title,
		e.url, e.comments_url, e.author, e.content, e.status,
		f.title as feed_title, f.feed_url, f.site_url, f.checked_at,
		f.category_id, c.title as category_title, f.scraper_rules, f.rewrite_rules, f.crawler,
		i.id, i.mime_type, i.content
		FROM entries e
		LEFT JOIN feeds f ON f.id=e.feed_id
		LEFT JOIN categories c ON c.id=f.category_id
		LEFT JOIN feed_icons fi ON fi.feed_id=f.id
		LEFT JOIN icons i ON fi.icon_id=i.id
		WHERE %s %s
	`

	args, conditions := e.buildCondition()
	query = fmt.Sprintf(query, conditions, e.buildSorting())

	rows, err := e.store.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to get entries: %v", err)
	}
	defer rows.Close()

	entries := make([]*model.EntryWithIcon, 0)
	for rows.Next() {
		var entry model.EntryWithIcon
		var tz string
		var icon = struct {
			ID       sql.NullInt64
			MimeType sql.NullString
			Content  []byte
		}{}

		entry.Feed = &model.Feed{}
		entry.Feed.Category = &model.Category{}

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
			&icon.ID,
			&icon.MimeType,
			&icon.Content,
		)

		if err != nil {
			return nil, fmt.Errorf("unable to fetch entry row: %v", err)
		}

		entry.EntryIcon = &model.Icon{}
		if icon.ID.Valid {
			entry.EntryIcon.ID = icon.ID.Int64
			entry.EntryIcon.MimeType = icon.MimeType.String
			entry.EntryIcon.Content = icon.Content
		} else {
			entry.EntryIcon.ID = 0
		}

		// Make sure that timestamp fields contains timezone information (API)
		entry.Date = timezone.Convert(tz, entry.Date)
		entry.Feed.CheckedAt = timezone.Convert(tz, entry.Feed.CheckedAt)

		entry.Feed.ID = entry.FeedID
		entries = append(entries, &entry)
	}

	return entries, nil
}

// GetEntryIDs returns a list of entry IDs that match the condition.
func (e *HeadlinesOfficialEntryQueryBuilder) GetEntryIDs() ([]int64, error) {
	debugStr := "[EntryQueryBuilder:GetEntryIDs] feedID=%d, categoryID=%d, status=%s, order=%s, direction=%s, filter=%s"
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf(debugStr, e.feedID, e.categoryID, e.status, e.order, e.direction, e.filter))

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

func (e *HeadlinesOfficialEntryQueryBuilder) buildCondition() ([]interface{}, string) {
	args := []interface{}{}
	conditions := []string{}

	if e.categoryID != 0 {
		conditions = append(conditions, fmt.Sprintf("f.category_id=$%d", len(args)+1))
		args = append(args, e.categoryID)
	}

	if len(e.notCategoryIDs) > 0 {
		for _, id := range e.notCategoryIDs {
			conditions = append(conditions, fmt.Sprintf("f.category_id != $%d", len(args)+1))
			args = append(args, id)
		}
	}

	if e.country != "" {
		conditions = append(conditions, fmt.Sprintf("e.title SIMILAR TO '%%(' || $%d || ')%%'", len(args)+1))
		args = append(args, e.country)
	}

	if e.nocountry != "" {
		conditions = append(conditions, fmt.Sprintf("e.title NOT SIMILAR TO '%%(' || $%d || ')%%'", len(args)+1))
		args = append(args, e.nocountry)
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

func (e *HeadlinesOfficialEntryQueryBuilder) buildSorting() string {
	var queries []string

	if e.order != "" {
		queries = append(queries, fmt.Sprintf(`ORDER BY "%s"`, e.order))
	}

	if e.direction != "" {
		queries = append(queries, fmt.Sprintf(`%s`, e.direction))
	}

	return strings.Join(queries, " ")
}

// NewHeadlinesOfficialEntryQueryBuilder returns a new HeadlinesOfficialEntryQueryBuilder.
func NewHeadlinesOfficialEntryQueryBuilder(store *Storage) *HeadlinesOfficialEntryQueryBuilder {
	return &HeadlinesOfficialEntryQueryBuilder{
		store: store,
	}
}
