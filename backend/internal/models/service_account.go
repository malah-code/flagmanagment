package models

import (
	"time"

	"github.com/google/uuid"
)

type ServiceAccount struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ServiceAccountKey struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	ServiceAccountID uuid.UUID  `json:"service_account_id" db:"service_account_id"`
	KeyHash          string     `json:"-" db:"key_hash"`
	Name             string     `json:"name" db:"name"`
	ExpiresAt        *time.Time `json:"expires_at" db:"expires_at"`
	LastUsedAt       *time.Time `json:"last_used_at" db:"last_used_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type ServiceAccountRole struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	ServiceAccountID uuid.UUID  `json:"service_account_id" db:"service_account_id"`
	RoleID           uuid.UUID  `json:"role_id" db:"role_id"`
	ProjectID        *uuid.UUID `json:"project_id" db:"project_id"`
	EnvironmentID    *uuid.UUID `json:"environment_id" db:"environment_id"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	Role             *Role      `json:"role,omitempty" db:"-"`
}
