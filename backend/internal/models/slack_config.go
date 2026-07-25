package models

import (
	"time"

	"github.com/google/uuid"
)

type SlackWebhookConfig struct {
	ID            uuid.UUID `json:"id" db:"id"`
	EnvironmentID uuid.UUID `json:"environment_id" db:"environment_id"`
	WebhookURL    string    `json:"webhook_url" db:"webhook_url"`
	Enabled       bool      `json:"enabled" db:"enabled"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
