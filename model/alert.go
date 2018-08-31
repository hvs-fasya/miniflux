// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

// Alert represents a Alert in the system.
type Alert struct {
	ID            int64     `json:"id"`
	LastUpdated   time.Time `json:"last_updated"`
	StillValid    time.Time `json:"still_valid"`
	LatestUpdates string    `json:"latest_updates"`
	CountryID     int64     `json:"country_id"`
	Risks         `json:"risks"`
}

// HealthAlert represents a HealthAlert record in the system.
type HealthAlert struct {
	Date    time.Time `json:"health_date"`
	Country `json:"country"`
	Health  `json:"health"`
}

// Health represents a Health in the system.
type Health struct {
	ID            int64     `json:"-"`
	HealthTitle   string    `json:"title"`
	HealthLink    string    `json:"-"`
	HealthContent string    `json:"content"`
	LastUpdated   time.Time `json:"-"`
}

// Alerts represents a list of Alerts.
type Alerts []*Alert

// Healths represents a list of Health objects.
type Healths []*Health

// HealthAlerts represents a list of HealthAlert objects.
type HealthAlerts []*HealthAlert

// Risks represents a Risks in the system.
type Risks struct {
	RiskLevel   string `json:"level"`
	RiskDetails string `json:"details"`
}

// NewAlert returns a new Alert.
func NewAlert() *Alert {
	return &Alert{}
}

// NewHealth returns a new Health.
func NewHealth() *Health {
	return &Health{}
}

// NewHealthAlert returns a new HealthAlert.
func NewHealthAlert() *HealthAlert {
	return &HealthAlert{}
}
