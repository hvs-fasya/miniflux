package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/timer"
	"strings"
)

// Filter returns a Filter from the database.
func (s *Storage) Filter(userID, filterID int64) (*model.Filter, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:Filter] userID=%d, getFilter=%d", userID, filterID))
	var filter model.Filter
	var pgFilters string
	query := `SELECT id, user_id, filter_name, filters FROM filters WHERE user_id=$1 AND id=$2`
	err := s.db.QueryRow(query, userID, filterID).Scan(&filter.ID, &filter.UserID, &filter.FilterName, &pgFilters)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to fetch filter: %v", err)
	}
	filter.Filters = strings.Split(pgFilters, ",")
	return &filter, nil
}

// FilterByName finds a filter by the name.
func (s *Storage) FilterByName(userID int64, name string) (*model.Filter, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:FilterByName] userID=%d, name=%s", userID, name))
	var filter model.Filter
	var pgFilters string
	query := `SELECT id, user_id, filter_name, array_to_string(filters,',') FROM filters WHERE user_id=$1 AND filter_name=$2`
	err := s.db.QueryRow(query, userID, name).Scan(&filter.ID, &filter.UserID, &filter.FilterName, &pgFilters)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("Unable to fetch filter: %v", err)
	}
	filter.Filters = strings.Split(pgFilters, ",")
	return &filter, nil
}

// FilterByNameWOUserID finds a filter by the name.
func (s *Storage) FilterByNameWOUserID(name string) (*model.Filter, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:FilterByName] name=%s", name))
	var filter model.Filter
	var pgFilters string
	query := `SELECT id, filter_name, array_to_string(filters,',') FROM filters WHERE filter_name=$1`
	err := s.db.QueryRow(query, name).Scan(&filter.ID, &filter.FilterName, &pgFilters)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("Unable to fetch filter: %v", err)
	}
	filter.Filters = strings.Split(pgFilters, ",")
	return &filter, nil
}

// Filters returns all Filters that belongs to the given user.
func (s *Storage) Filters(userID int64) (model.Filters, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:Filters] userID=%d", userID))

	query := `SELECT id, user_id, filter_name, array_to_string(filters,',') FROM filters WHERE user_id=$1 ORDER BY created_at DESC`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch filters: %v", err)
	}
	defer rows.Close()

	filters := make(model.Filters, 0)
	var pgFilters string
	for rows.Next() {
		var filter model.Filter
		if err := rows.Scan(&filter.ID, &filter.UserID, &filter.FilterName, &pgFilters); err != nil {
			return nil, fmt.Errorf("Unable to fetch filters row: %v", err)
		}
		filter.Filters = strings.Split(pgFilters, ",")
		filters = append(filters, filter)
	}

	return filters, nil
}

// CreateFilter creates a new filter.
func (s *Storage) CreateFilter(filter *model.Filter) error {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:CreateFilter] title=%s", filter.FilterName))

	query := `
		INSERT INTO filters
		(user_id, filter_name, filters)
		VALUES
		($1, $2, string_to_array($3,','))
		RETURNING id
	`
	err := s.db.QueryRow(
		query,
		filter.UserID,
		filter.FilterName,
		strings.Join(filter.Filters, ",")).Scan(&filter.ID)

	if err != nil {
		return fmt.Errorf("Unable to create filter: %v", err)
	}

	return nil
}

// RemoveFilter deletes a filter.
func (s *Storage) RemoveFilter(userID, filterID int64) error {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:RemoveFilter] userID=%d, filterID=%d", userID, filterID))

	result, err := s.db.Exec("DELETE FROM filters WHERE id = $1 AND user_id = $2", filterID, userID)
	if err != nil {
		return fmt.Errorf("Unable to remove this filter: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Unable to remove this filter: %v", err)
	}

	if count == 0 {
		return errors.New("no filter has been removed")
	}

	return nil
}

// FilterExists checks if the given filter exists into the database.
func (s *Storage) FilterExists(userID, filterID int64) bool {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:FilterExists] userID=%d, filterID=%d", userID, filterID))

	var result int
	query := `SELECT count(*) as c FROM filters WHERE user_id=$1 AND id=$2`
	s.db.QueryRow(query, userID, filterID).Scan(&result)
	return result >= 1
}
