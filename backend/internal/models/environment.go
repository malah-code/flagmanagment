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
}

