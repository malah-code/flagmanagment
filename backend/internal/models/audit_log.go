package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	ProjectID     *uuid.UUID `json:"project_id,omitempty" db:"project_id"`
	EnvironmentID *uuid.UUID `json:"environment_id,omitempty" db:"environment_id"`
	ActorID       uuid.UUID  `json:"actor_id" db:"actor_id"`
	Action        string     `json:"action" db:"action"`
	TargetType    string     `json:"target_type" db:"target_type"`
	TargetID      uuid.UUID  `json:"target_id" db:"target_id"`
	PreviousState JSONB      `json:"previous_state,omitempty" db:"previous_state"`
	NewState      JSONB      `json:"new_state,omitempty" db:"new_state"`
	ActorIP       string     `json:"actor_ip" db:"actor_ip"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}
