package models

import "encoding/json"

// RulesetSnapshot represents a complete set of flags and rules for an environment.
type RulesetSnapshot struct {
	Version string     `json:"version"`
	Flags   []FlagRule `json:"flags"`
}

// FlagRule defines a single feature flag's evaluation configuration.
type FlagRule struct {
	Key              string          `json:"key"`
	Type             string          `json:"type"` // BOOLEAN, MULTIVARIATE, STRING, NUMBER, JSON
	Enabled          bool            `json:"enabled"`
	DefaultVariation string          `json:"defaultVariation"`
	TargetingRules   json.RawMessage `json:"targetingRules,omitempty"`
	RolloutRules     json.RawMessage `json:"rolloutRules,omitempty"`
	Variations       json.RawMessage `json:"variations,omitempty"`
	ParentFlagKey    string          `json:"parentFlagKey,omitempty"`
}

// EvaluationContext represents the user attributes passed for flag evaluation.
type EvaluationContext struct {
	EntityKey  string                 `json:"entityKey"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// EvaluationResult represents the output of a flag evaluation.
type EvaluationResult struct {
	Value  interface{} `json:"value"`
	Reason string      `json:"reason"`
}
