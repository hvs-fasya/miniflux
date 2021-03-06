// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package storage

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/timer"
)

// NewHeadlinesQueryBuilder returns a new HeadlinesQueryBuilder
func (s *Storage) NewHeadlinesQueryBuilder() *HeadlinesQueryBuilder {
	return NewHeadlinesQueryBuilder(s)
}

// HeadlineExists checks if a headline exists by using the given title + content.
func (s *Storage) HeadlineExists(title string, content string) bool {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:HeadlineExists] title=%s", title))

	hash := hashHeadline(title + content)
	var result int
	s.db.QueryRow(`SELECT count(*) as c FROM headlines WHERE hash=$1`, hash).Scan(&result)
	return result >= 1
}

// CreateHeadline creates a new headline.
func (s *Storage) CreateHeadline(headline *model.Headline) (err error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:CreateHeadline] title=%s", headline.Title))
	headline.Hash = hashHeadline(headline.Title + headline.Content)
	headline.PublishedAt = time.Now()
	if headline.CategoryID == 0 {
		headline.CategoryID = 1
	}
	if err != nil {
		return err
	}
	query := `INSERT INTO headlines
		(hash, published_at, title, content, url, country_id, visatype, category_id)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, hash, published_at, title, content, url, country_id, visatype, category_id`

	err = s.db.QueryRow(query, headline.Hash, headline.PublishedAt, headline.Title, headline.Content, headline.Url, headline.CountryID, headline.VisaType, headline.CategoryID).Scan(
		&headline.ID,
		&headline.Hash,
		&headline.PublishedAt,
		&headline.Title,
		&headline.Content,
		&headline.Url,
		&headline.CountryID,
		&headline.VisaType,
		&headline.CategoryID,
	)
	if err != nil {
		return fmt.Errorf("unable to create headline: %v", err)
	}

	return nil
}

// UpdateHeadline updates a headline.
func (s *Storage) UpdateHeadline(headline *model.Headline) error {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:UpdateHeadline] headlineID=%d", headline.ID))

	headline.Hash = hashHeadline(headline.Title + headline.Content)

	query := `UPDATE headlines SET
			hash=$1,
			title=$2,
			content=$3,
			url=$4,
			country_id=$5,
			visatype=$6,
			category_id=$7
			WHERE id=$8`

	_, err := s.db.Exec(
		query,
		&headline.Hash,
		&headline.Title,
		&headline.Content,
		&headline.Url,
		&headline.CountryID,
		&headline.VisaType,
		&headline.CategoryID,
		headline.ID,
	)
	if err != nil {
		return fmt.Errorf("unable to update headline: %v", err)
	}

	return nil
}

// HeadlineByID finds a Headline by the ID.
func (s *Storage) HeadlineByID(headlineID int64) (*model.Headline, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:UserByID] headlineID=%d", headlineID))
	query := `SELECT
		id, hash, published_at, title, content, url, country_id, visatype, category_id
		FROM headlines
		WHERE id = $1`

	return s.fetchHeadline(query, headlineID)
}

func (s *Storage) fetchHeadline(query string, args ...interface{}) (*model.Headline, error) {

	headline := model.NewHeadline()
	err := s.db.QueryRow(query, args...).Scan(
		&headline.ID,
		&headline.Hash,
		&headline.PublishedAt,
		&headline.Title,
		&headline.Content,
		&headline.Url,
		&headline.CountryID,
		&headline.VisaType,
		&headline.CategoryID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to fetch headline: %v", err)
	}

	return headline, nil
}

// RemoveHeadline deletes a headline.
func (s *Storage) RemoveHeadline(headlineID int64) error {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:RemoveHeadline] headlineID=%d", headlineID))

	result, err := s.db.Exec(`DELETE FROM headlines WHERE id = $1`, headlineID)
	if err != nil {
		return fmt.Errorf("unable to remove this headline: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to remove this headline: %v", err)
	}

	if count == 0 {
		return errors.New("nothing has been removed")
	}

	return nil
}

// Headlines returns all Headlines.
func (s *Storage) Headlines() (model.Headlines, error) {
	defer timer.ExecutionTime(time.Now(), "[Storage:Headlines]")
	query := `
		SELECT
			id, hash, published_at, title, content, url, country_id, visatype, category_id
		FROM headlines
		ORDER BY published_at ASC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch headlines: %v", err)
	}
	defer rows.Close()

	var headlines model.Headlines
	for rows.Next() {
		headline := model.NewHeadline()
		err := rows.Scan(
			&headline.ID,
			&headline.Hash,
			&headline.PublishedAt,
			&headline.Title,
			&headline.Content,
			&headline.Url,
			&headline.CountryID,
			&headline.VisaType,
			&headline.CategoryID,
		)

		if err != nil {
			return nil, fmt.Errorf("unable to fetch headline row: %v", err)
		}

		headlines = append(headlines, headline)
	}

	return headlines, nil
}

func hashHeadline(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
