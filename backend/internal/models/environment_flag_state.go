package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type EnvironmentFlagState struct {
	ID             uuid.UUID `json:"id" db:"id"`
	EnvironmentID  uuid.UUID `json:"environment_id" db:"environment_id"`
	FeatureFlagID  uuid.UUID `json:"feature_flag_id" db:"feature_flag_id"`
	Enabled        bool      `json:"enabled" db:"enabled"`
	TargetingRules JSONB     `json:"targeting_rules" db:"targeting_rules"`
	RemoteConfig   JSONB     `json:"remote_config" db:"remote_config"`
	Variations     JSONB     `json:"variations,omitempty" db:"variations"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Validate checks the basic structure of targeting rules
func (e *EnvironmentFlagState) Validate() error {
	// A simple mock validation for JSONB targeting rules per FR-014
	if e.TargetingRules != nil {
		if _, ok := e.TargetingRules["rules"]; !ok {
			return errors.New("targeting_rules must contain a 'rules' key")
		}
	}
	return nil
}
