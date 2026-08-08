package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type WebhookIntegration struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	ProjectID uuid.UUID       `json:"project_id" db:"project_id"`
	URL       string          `json:"url" db:"url"`
	SecretKey *string         `json:"secret_key,omitempty" db:"secret_key"`
	Events    json.RawMessage `json:"events" db:"events"`
	IsActive  bool            `json:"is_active" db:"is_active"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type WebhookPayload struct {
	EventID   string      `json:"event_id"`
	EventType string      `json:"event_type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}
