package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Permissions JSONB     `json:"permissions" db:"permissions"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserRole struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	RoleID        uuid.UUID  `json:"role_id" db:"role_id"`
	ProjectID     *uuid.UUID `json:"project_id,omitempty" db:"project_id"`
	EnvironmentID *uuid.UUID `json:"environment_id,omitempty" db:"environment_id"`
	Role          *Role      `json:"role,omitempty" db:"-"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}
