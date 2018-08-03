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

// NewVisaupdatesQueryBuilder returns a new VisaupdatesQueryBuilder
func (s *Storage) NewVisaupdatesQueryBuilder() *VisaupdatesQueryBuilder {
	return NewVisaupdatesQueryBuilder(s)
}

// VisaupdateExists checks if a visaupdate exists by using the given title + content.
func (s *Storage) VisaupdateExists(title string, content string) bool {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:VisaupdateExists] title=%s", title))

	hash := hashVisaupdate(title + content)
	var result int
	s.db.QueryRow(`SELECT count(*) as c FROM visaupdates WHERE hash=$1`, hash).Scan(&result)
	return result >= 1
}

// CreateVisaupdate creates a new visaupdate.
func (s *Storage) CreateVisaupdate(visaupdate *model.Visaupdate) (err error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:CreateVisaupdate] title=%s", visaupdate.Title))
	visaupdate.Hash = hashHeadline(visaupdate.Title + visaupdate.Content)
	visaupdate.PublishedAt = time.Now()
	if err != nil {
		return err
	}
	query := `INSERT INTO visaupdates
		(hash, published_at, title, content, url, country_id, visatype)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, hash, published_at, title, content, authority, country_id, visatype,`

	err = s.db.QueryRow(query, visaupdate.Hash, visaupdate.PublishedAt, visaupdate.Title, visaupdate.Content, visaupdate.Authority, visaupdate.CountryID, visaupdate.VisaType).Scan(
		&visaupdate.ID,
		&visaupdate.Hash,
		&visaupdate.PublishedAt,
		&visaupdate.Title,
		&visaupdate.Content,
		&visaupdate.Authority,
		&visaupdate.CountryID,
		&visaupdate.VisaType,
	)
	if err != nil {
		return fmt.Errorf("unable to create Visaupdate: %v", err)
	}

	return nil
}

// UpdateVisaupdate updates a visaupdate.
func (s *Storage) UpdateVisaupdate(visaupdate *model.Visaupdate) error {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:UpdateVisaupdate] visaupdateID=%d", visaupdate.ID))

	visaupdate.Hash = hashHeadline(visaupdate.Title + visaupdate.Content)

	query := `UPDATE visaupdates SET
			hash=$1,
			title=$2,
			content=$3,
			url=$4,
			country_id=$5,
			visatype=$6,
			WHERE id=$7`

	_, err := s.db.Exec(
		query,
		&visaupdate.Hash,
		&visaupdate.Title,
		&visaupdate.Content,
		&visaupdate.Authority,
		&visaupdate.CountryID,
		&visaupdate.VisaType,
		visaupdate.ID,
	)
	if err != nil {
		return fmt.Errorf("unable to update visaupdate: %v", err)
	}

	return nil
}

// VisaupdateByID finds a Visaupdate by the ID.
func (s *Storage) VisaupdateByID(visaupdateID int64) (*model.Visaupdate, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:UserByID] visaupdateID=%d", visaupdateID))
	query := `SELECT
		id, hash, published_at, title, content, url, country_id, visatype
		FROM visaupdates
		WHERE id = $1`

	return s.fetchVisaupdate(query, visaupdateID)
}

func (s *Storage) fetchVisaupdate(query string, args ...interface{}) (*model.Visaupdate, error) {

	visaupdate := model.NewVisaupdate()
	err := s.db.QueryRow(query, args...).Scan(
		&visaupdate.ID,
		&visaupdate.Hash,
		&visaupdate.PublishedAt,
		&visaupdate.Title,
		&visaupdate.Content,
		&visaupdate.Authority,
		&visaupdate.CountryID,
		&visaupdate.VisaType,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to fetch visaupdate: %v", err)
	}

	return visaupdate, nil
}

// RemoveVisaupdate deletes a visaupdate.
func (s *Storage) RemoveVisaupdate(visaupdateID int64) error {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:RemoveVisaupdate] visaupdateID=%d", visaupdateID))

	result, err := s.db.Exec(`DELETE FROM visaupdates WHERE id = $1`, visaupdateID)
	if err != nil {
		return fmt.Errorf("unable to remove this visaupdate: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to remove this visaupdate: %v", err)
	}

	if count == 0 {
		return errors.New("nothing has been removed")
	}

	return nil
}

// Visaupdates returns all Visaupdates.
func (s *Storage) Visaupdates() (model.Visaupdates, error) {
	defer timer.ExecutionTime(time.Now(), "[Storage:Visaupdates]")
	query := `
		SELECT
			id, hash, published_at, title, content, authority, country_id, visatype
		FROM visaupdates
		ORDER BY published_at ASC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch visaupdates: %v", err)
	}
	defer rows.Close()

	var visaupdates model.Visaupdates
	for rows.Next() {
		visaupdate := model.NewVisaupdate()
		err := rows.Scan(
			&visaupdate.ID,
			&visaupdate.Hash,
			&visaupdate.PublishedAt,
			&visaupdate.Title,
			&visaupdate.Content,
			&visaupdate.Authority,
			&visaupdate.CountryID,
			&visaupdate.VisaType,
		)

		if err != nil {
			return nil, fmt.Errorf("unable to fetch visaupdate row: %v", err)
		}

		visaupdates = append(visaupdates, visaupdate)
	}

	return visaupdates, nil
}

func hashVisaupdate(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
