package models

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	Email      string     `json:"email" db:"email"`
	TokenHash  string     `json:"-" db:"token_hash"`
	Role       string     `json:"role" db:"role"`
	ProjectIDs JSONB      `json:"project_ids" db:"project_ids"`
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	CreatedBy  *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}
