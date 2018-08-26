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
	Health        `json:"health"`
}

// Health represents a Health in the system.
type Health struct {
	HealthTitle   string    `json:"title"`
	HealthDate    time.Time `json:"date"`
	HealthContent string    `json:"content"`
}

// Risks represents a Risks in the system.
type Risks struct {
	RiskLevel   string `json:"level"`
	RiskDetails string `json:"details"`
}

// NewAlert returns a new Alert.
func NewAlert() *Alert {
	return &Alert{}
}

// Alerts represents a list of Alerts.
type Alerts []*Alert
