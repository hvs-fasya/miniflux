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

//
//// NewHeadlinesQueryBuilder returns a new HeadlinesQueryBuilder
//func (s *Storage) NewHeadlinesQueryBuilder() *HeadlinesQueryBuilder {
//	return NewHeadlinesQueryBuilder(s)
//}

// HealthExists checks if a health exists by using the health_title.
func (s *Storage) HealthExists(healthTitle string) bool {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:HealthExists] healthTitle=%s", healthTitle))

	var result int
	s.db.QueryRow(`SELECT count(*) as c FROM healths WHERE health_title=$1`, healthTitle).Scan(&result)
	return result >= 1
}

// CreateHealth creates a new Health.
func (s *Storage) CreateHealth(health *model.Health) (*model.Health, error) {
	var err error
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:CreateHealth] health_title=%s", health.HealthTitle))

	query := `INSERT INTO healths
		(health_link, health_title, health_content, last_updated)
		VALUES
		($1, $2, $3, $4)
		RETURNING id
		`

	err = s.db.QueryRow(query, health.HealthLink, health.HealthTitle, health.HealthContent, time.Now()).Scan(&health.ID)
	if err != nil {
		return health, fmt.Errorf("unable to create health: %v", err)
	}

	return health, nil
}

// UpdateHealth updates a Health.
func (s *Storage) UpdateHealth(health *model.Health) error {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:UpdateHealth] health_title=%s", health.HealthTitle))

	query := `UPDATE healths SET
			health_link=$1,
			health_content=$2,
			last_updated=$3
			WHERE health_title=$4`

	_, err := s.db.Exec(
		query,
		&health.HealthLink,
		&health.HealthContent,
		time.Now(),
		health.HealthTitle,
	)
	if err != nil {
		return fmt.Errorf("unable to update health: %v", err)
	}

	return nil
}

// HealthByTitle finds a Health by the health_title.
func (s *Storage) HealthByTitle(healthTitle string) (*model.Health, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:HealthByTitle] healthTitle=%s", healthTitle))
	query := `SELECT
		id, health_link, health_title, health_content, last_updated
		FROM healths
		WHERE health_title = $1`

	return s.fetchHealth(query, healthTitle)
}

func (s *Storage) fetchHealth(query string, args ...interface{}) (*model.Health, error) {

	health := model.NewHealth()
	err := s.db.QueryRow(query, args...).Scan(
		&health.ID,
		&health.HealthLink,
		&health.HealthTitle,
		&health.HealthContent,
		&health.LastUpdated,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to fetch health: %v", err)
	}

	return health, nil
}

//// RemoveHeadline deletes a headline.
//func (s *Storage) RemoveHeadline(headlineID int64) error {
//	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Storage:RemoveHeadline] headlineID=%d", headlineID))
//
//	result, err := s.db.Exec(`DELETE FROM headlines WHERE id = $1`, headlineID)
//	if err != nil {
//		return fmt.Errorf("unable to remove this headline: %v", err)
//	}
//
//	count, err := result.RowsAffected()
//	if err != nil {
//		return fmt.Errorf("unable to remove this headline: %v", err)
//	}
//
//	if count == 0 {
//		return errors.New("nothing has been removed")
//	}
//
//	return nil
//}
//
//// Headlines returns all Headlines.
//func (s *Storage) Headlines() (model.Headlines, error) {
//	defer timer.ExecutionTime(time.Now(), "[Storage:Headlines]")
//	query := `
//		SELECT
//			id, hash, published_at, title, content, url, country_id, visatype, category_id
//		FROM headlines
//		ORDER BY published_at ASC`
//
//	rows, err := s.db.Query(query)
//	if err != nil {
//		return nil, fmt.Errorf("unable to fetch headlines: %v", err)
//	}
//	defer rows.Close()
//
//	var headlines model.Headlines
//	for rows.Next() {
//		headline := model.NewHeadline()
//		err := rows.Scan(
//			&headline.ID,
//			&headline.Hash,
//			&headline.PublishedAt,
//			&headline.Title,
//			&headline.Content,
//			&headline.Url,
//			&headline.CountryID,
//			&headline.VisaType,
//			&headline.CategoryID,
//		)
//
//		if err != nil {
//			return nil, fmt.Errorf("unable to fetch headline row: %v", err)
//		}
//
//		headlines = append(headlines, headline)
//	}
//
//	return headlines, nil
//}
