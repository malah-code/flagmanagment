package models

import (
	"time"

	"github.com/google/uuid"
)

type ScheduledChangeTargetType string
type ScheduledChangeAction string
type ScheduledChangeStatus string

const (
	TargetTypeFlag          ScheduledChangeTargetType = "FLAG"
	TargetTypeChangeRequest ScheduledChangeTargetType = "CHANGE_REQUEST"

	ActionEnable  ScheduledChangeAction = "ENABLE"
	ActionDisable ScheduledChangeAction = "DISABLE"
	ActionApply   ScheduledChangeAction = "APPLY"

	ScheduleStatusPending   ScheduledChangeStatus = "PENDING"
	ScheduleStatusExecuted  ScheduledChangeStatus = "EXECUTED"
	ScheduleStatusCancelled ScheduledChangeStatus = "CANCELLED"
)

// systemActorUUID is the sentinel UUID used in audit logs for scheduler-driven actions.
// It is intentionally unexported; callers must use SystemActor().
var systemActorID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

// SystemActor returns the sentinel UUID for the automated scheduler actor.
// Returning by value prevents callers from mutating the sentinel.
func SystemActor() uuid.UUID { return systemActorID }

type ScheduledChange struct {
	ID            uuid.UUID                 `json:"id" db:"id"`
	ProjectID     uuid.UUID                 `json:"project_id" db:"project_id"`
	EnvironmentID uuid.UUID                 `json:"environment_id" db:"environment_id"`
	TargetType    ScheduledChangeTargetType `json:"target_type" db:"target_type"`
	TargetID      uuid.UUID                 `json:"target_id" db:"target_id"`
	Action        ScheduledChangeAction     `json:"action" db:"action"`
	ScheduledFor  time.Time                 `json:"scheduled_for" db:"scheduled_for"`
	Status        ScheduledChangeStatus     `json:"status" db:"status"`
	CreatedBy     uuid.UUID                 `json:"created_by" db:"created_by"`
	ExecutedAt    *time.Time                `json:"executed_at,omitempty" db:"executed_at"`
	CancelledAt   *time.Time                `json:"cancelled_at,omitempty" db:"cancelled_at"`
	CreatedAt     time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at" db:"updated_at"`
}
