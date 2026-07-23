package dto

import (
	"time"
)

// Base Responses
type PaginatedResponse struct {
	Data          any    `json:"data"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// Project Responses
type ProjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Environment Responses
type EnvironmentResponse struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"projectId"`
	Name        string    `json:"name"`
	IsProtected bool      `json:"isProtected"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateEnvironmentResponse struct {
	EnvironmentResponse
	APIKey string `json:"apiKey"` // Only returned on creation
}

// Feature Flag Responses
type FeatureFlagResponse struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"projectId"`
	Key          string    `json:"key"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Type         string    `json:"type"`
	ParentFlagID *string   `json:"parentFlagId,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Flag State Responses
type FlagStateResponse struct {
	EnvironmentID  string                 `json:"environmentId"`
	FeatureFlagID  string                 `json:"featureFlagId"`
	Enabled        bool                   `json:"enabled"`
	TargetingRules map[string]interface{} `json:"targetingRules"`
	RemoteConfig   map[string]interface{} `json:"remoteConfig"`
	UpdatedAt      time.Time              `json:"updatedAt"`
}

// SDK Responses
type SDKFlag struct {
	Enabled bool                   `json:"enabled"`
	Type    string                 `json:"type"`
	Rules   map[string]interface{} `json:"rules"`
	Value   map[string]interface{} `json:"value"`
}

type SDKEvaluationResponse struct {
	EnvironmentID string             `json:"environmentId"`
	Flags         map[string]SDKFlag `json:"flags"`
}
