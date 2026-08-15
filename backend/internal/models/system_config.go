package models

import "time"

type SystemConfig struct {
	Key       string    `json:"key" db:"key"`
	Value     JSONB     `json:"value" db:"value"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
