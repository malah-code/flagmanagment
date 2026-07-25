package models

import (
	"time"

	"github.com/google/uuid"
)

type ChangeRequest struct {
	ID              uuid.UUID           `json:"id" db:"id"`
	ProjectID       uuid.UUID           `json:"project_id" db:"project_id"`
	EnvironmentID   uuid.UUID           `json:"environment_id" db:"environment_id"`
	Title           string              `json:"title" db:"title"`
	Description     string              `json:"description" db:"description"`
	Status          ChangeRequestStatus `json:"status" db:"status"`
	ProposedChanges JSONB               `json:"proposed_changes" db:"proposed_changes"`
	CurrentState    JSONB               `json:"current_state" db:"current_state"`
	CreatedBy       uuid.UUID           `json:"created_by" db:"created_by"`
	AppliedBy       *uuid.UUID          `json:"applied_by,omitempty" db:"applied_by"`
	CreatedAt       time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at" db:"updated_at"`
}

type ChangeRequestApproval struct {
	ID              uuid.UUID        `json:"id" db:"id"`
	ChangeRequestID uuid.UUID        `json:"change_request_id" db:"change_request_id"`
	ApproverID      uuid.UUID        `json:"approver_id" db:"approver_id"`
	Decision        ApprovalDecision `json:"decision" db:"decision"`
	Comment         string           `json:"comment" db:"comment"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
}
