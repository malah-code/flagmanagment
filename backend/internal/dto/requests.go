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
	Name        string `json:"name" validate:"required,min=1,max=100"`
	IsProtected bool   `json:"isProtected"`
}

type UpdateEnvironmentRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=100"`
	IsProtected bool                   `json:"isProtected"`
	SdkSettings map[string]interface{} `json:"sdkSettings,omitempty"`
}

type CloneEnvironmentRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type CreateServerKeyRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type VariationDTO struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Value       interface{} `json:"value"`
}

// Feature Flag Requests
type CreateFeatureFlagRequest struct {
	Key          string         `json:"key" validate:"required,min=3,max=100"`
	Name         string         `json:"name" validate:"required,min=3,max=100"`
	Description  string         `json:"description"`
	Type             string         `json:"type" validate:"required"`
	EnabledByDefault *bool          `json:"enabledByDefault,omitempty"`
	Variations       []VariationDTO `json:"variations,omitempty"`
	Tags             []string       `json:"tags,omitempty"`
	ParentFlagID     *string        `json:"parentFlagId" validate:"omitempty,uuid"`
}

type UpdateFeatureFlagRequest struct {
	Name         string         `json:"name" validate:"required,min=3,max=100"`
	Description  string         `json:"description"`
	Tags         []string       `json:"tags,omitempty"`
	ParentFlagID *string        `json:"parentFlagId" validate:"omitempty,uuid"`
}

type UpdateFlagStateRequest struct {
	Enabled          *bool                  `json:"enabled,omitempty"`
	DefaultVariation string                 `json:"defaultVariation,omitempty"`
	TargetingRules   map[string]interface{} `json:"targetingRules,omitempty"`
	RemoteConfig     map[string]interface{} `json:"remoteConfig,omitempty"`
	RolloutRules     map[string]interface{} `json:"rolloutRules,omitempty"`
}

// Change Request Requests
type ApproveChangeRequestRequest struct {
	Comment string `json:"comment"`
}

type RejectChangeRequestRequest struct {
	Reason string `json:"reason"`
}
