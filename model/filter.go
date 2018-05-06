package model

import (
	"errors"
	"fmt"
	"strings"
)

// Filter represents a Filter in the system.
type Filter struct {
	ID         int64    `json:"id,omitempty"`
	UserID     int64    `json:"user_id,omitempty"`
	FilterName string   `json:"filter_name,omitempty"`
	Filters    []string `json:"filters,omitempty"`
}

func (c *Filter) String() string {
	return fmt.Sprintf("ID=%d, UserID=%d, FilterName=%s, Filters=%s", c.ID, c.UserID, c.FilterName, strings.Join(c.Filters, ", "))
}

// ValidateFilterCreation validates a filter during the creation.
func (c Filter) ValidateFilterCreation() error {
	if c.FilterName == "" {
		return errors.New("The filter name is mandatory")
	}

	if c.UserID == 0 {
		return errors.New("The userID is mandatory")
	}

	return nil
}

// Filters represents a list of filters.
type Filters []Filter
