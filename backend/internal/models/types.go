package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// FlagType defines the type of feature flag
type FlagType string

const (
	FlagTypeBoolean      FlagType = "BOOLEAN"
	FlagTypeMultivariate FlagType = "MULTIVARIATE"
	FlagTypeJSON         FlagType = "JSON"
)

func (f FlagType) IsValid() bool {
	switch f {
	case FlagTypeBoolean, FlagTypeMultivariate, FlagTypeJSON:
		return true
	default:
		return false
	}
}

// ChangeRequestStatus defines the state of a change request
type ChangeRequestStatus string

const (
	StatusPending  ChangeRequestStatus = "PENDING"
	StatusApproved ChangeRequestStatus = "APPROVED"
	StatusRejected ChangeRequestStatus = "REJECTED"
	StatusApplied  ChangeRequestStatus = "APPLIED"
)

// ApprovalDecision defines an approver's decision
type ApprovalDecision string

const (
	DecisionApprove ApprovalDecision = "APPROVE"
	DecisionReject  ApprovalDecision = "REJECT"
)

// JSONB represents a PostgreSQL JSONB column
type JSONB map[string]interface{}

// Value implements driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}
