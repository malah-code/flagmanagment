package models

import (
	"time"

	"github.com/google/uuid"
)

type FlagLifecycleState string

const (
	LifecycleActive     FlagLifecycleState = "ACTIVE"
	LifecycleStale      FlagLifecycleState = "STALE"
	LifecycleDeprecated FlagLifecycleState = "DEPRECATED"
	LifecycleArchived   FlagLifecycleState = "ARCHIVED"
)

func (s FlagLifecycleState) IsValid() bool {
	switch s {
	case LifecycleActive, LifecycleStale, LifecycleDeprecated, LifecycleArchived:
		return true
	default:
		return false
	}
}

type StaleFlagPolicy struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	ProjectID      uuid.UUID  `json:"project_id" db:"project_id"`
	EnvironmentID  *uuid.UUID `json:"environment_id,omitempty" db:"environment_id"`
	StaleAfterDays int        `json:"stale_after_days" db:"stale_after_days"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}
