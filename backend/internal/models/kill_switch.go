package models

import (
	"time"

	"github.com/google/uuid"
)

type KillSwitchRule struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	FlagID          uuid.UUID  `json:"flag_id" db:"flag_id"`
	EnvironmentID   uuid.UUID  `json:"environment_id" db:"environment_id"`
	AlertIdentifier string     `json:"alert_identifier" db:"alert_identifier"`
	Action          string     `json:"action" db:"action"`
	CreatedBy       *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}
