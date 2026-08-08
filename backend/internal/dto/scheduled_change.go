package dto

import (
	"time"
)

type CreateScheduledChangeRequest struct {
	TargetType   string `json:"target_type" validate:"required,oneof=FLAG CHANGE_REQUEST"`
	TargetID     string `json:"target_id" validate:"required,uuid"`
	Action       string `json:"action" validate:"required,oneof=ENABLE DISABLE APPLY"`
	ScheduledFor string `json:"scheduled_for" validate:"required"`
}

type UpdateScheduledChangeRequest struct {
	ScheduledFor string `json:"scheduled_for" validate:"required"`
}

type ScheduledChangeResponse struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	EnvironmentID string    `json:"environment_id"`
	TargetType    string    `json:"target_type"`
	TargetID      string    `json:"target_id"`
	Action        string    `json:"action"`
	ScheduledFor  time.Time `json:"scheduled_for"`
	Status        string    `json:"status"`
	CreatedBy     string    `json:"created_by"`
	ExecutedAt    *time.Time `json:"executed_at,omitempty"`
	CancelledAt   *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
