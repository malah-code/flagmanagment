package models

import (
	"time"

	"github.com/google/uuid"
)

type Environment struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Key         string    `json:"key" db:"key"`
	APIKeyHash  string    `json:"-" db:"api_key_hash"`
	Salt        string    `json:"salt" db:"salt"`
	IsProtected bool      `json:"is_protected" db:"is_protected"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	SDKSettings JSONB     `json:"sdkSettings" db:"sdk_settings"`
}

type EnvironmentServerKey struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	EnvironmentID uuid.UUID  `json:"environment_id" db:"environment_id"`
	Name          string     `json:"name" db:"name"`
	KeyHash       string     `json:"-" db:"key_hash"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	LastUsedAt    *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
}

