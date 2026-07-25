package dto

// Project Requests
type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Key         string `json:"key" validate:"omitempty,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description"`
}

// Environment Requests
type CreateEnvironmentRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	IsProtected bool   `json:"isProtected"`
}

type UpdateEnvironmentRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	IsProtected bool   `json:"isProtected"`
}

// Feature Flag Requests
type CreateFeatureFlagRequest struct {
	Key          string  `json:"key" validate:"required,min=3,max=100"`
	Name         string  `json:"name" validate:"required,min=3,max=100"`
	Description  string  `json:"description"`
	Type         string  `json:"type" validate:"required,oneof=boolean string number json"`
	ParentFlagID *string `json:"parentFlagId" validate:"omitempty,uuid"`
}

type UpdateFlagStateRequest struct {
	Enabled        bool                   `json:"enabled"`
	TargetingRules map[string]interface{} `json:"targetingRules" validate:"required"`
	RemoteConfig   map[string]interface{} `json:"remoteConfig" validate:"required"`
}

// Change Request Requests
type ApproveChangeRequestRequest struct {
	Comment string `json:"comment"`
}

type RejectChangeRequestRequest struct {
	Reason string `json:"reason"`
}
