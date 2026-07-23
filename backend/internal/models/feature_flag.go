package models

import (
	"time"

	"github.com/google/uuid"
)

type FeatureFlag struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	ProjectID       uuid.UUID  `json:"project_id" db:"project_id"`
	Key             string     `json:"key" db:"key"`
	Name            string     `json:"name" db:"name"`
	Description     string     `json:"description" db:"description"`
	Type            FlagType   `json:"type" db:"type"`
	ParentFlagID    *uuid.UUID `json:"parent_flag_id,omitempty" db:"parent_flag_id"`
	LastEvaluatedAt *time.Time `json:"last_evaluated_at,omitempty" db:"last_evaluated_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}
