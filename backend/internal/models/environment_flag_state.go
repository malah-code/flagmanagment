package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type EnvironmentFlagState struct {
	ID               uuid.UUID `json:"id" db:"id"`
	EnvironmentID    uuid.UUID `json:"environment_id" db:"environment_id"`
	FeatureFlagID    uuid.UUID `json:"feature_flag_id" db:"feature_flag_id"`
	Enabled          bool      `json:"enabled" db:"enabled"`
	DefaultVariation string    `json:"default_variation,omitempty" db:"default_variation"`
	TargetingRules   JSONB     `json:"targeting_rules" db:"targeting_rules"`
	RemoteConfig     JSONB     `json:"remote_config" db:"remote_config"`
	RolloutRules     JSONB     `json:"rollout_rules,omitempty" db:"rollout_rules"`
	Variations       JSONB     `json:"variations,omitempty" db:"variations"`
	LifecycleState   FlagLifecycleState `json:"lifecycle_state" db:"lifecycle_state"`
	LastEvaluatedAt  *time.Time         `json:"last_evaluated_at,omitempty" db:"last_evaluated_at"`
	LastStateChangeAt time.Time         `json:"last_state_change_at" db:"last_state_change_at"`
	CreatedAt        time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" db:"updated_at"`
}

type RolloutRule struct {
	VariationID string `json:"variation_id"`
	Percentage  int    `json:"percentage"` // basis points out of 10000
}

type Operator string

const (
	OperatorEquals   Operator = "EQUALS"
	OperatorContains Operator = "CONTAINS"
	OperatorRegex    Operator = "REGEX"
)

type TargetingCondition struct {
	Attribute string   `json:"attribute"`
	Operator  Operator `json:"operator"`
	Value     string   `json:"value"`
}

type TargetingRule struct {
	ID         string               `json:"id"`
	Conditions []TargetingCondition `json:"conditions"`
	Variation  bool                 `json:"variation"`
}

// Validate checks the basic structure of targeting rules and compiles regexes
func (e *EnvironmentFlagState) Validate() error {
	if e.TargetingRules != nil {
		if _, ok := e.TargetingRules["rules"]; ok {
			rulesBytes, err := json.Marshal(e.TargetingRules["rules"])
			if err != nil {
				return err
			}
			var rules []TargetingRule
			if err := json.Unmarshal(rulesBytes, &rules); err != nil {
				return err
			}
			for _, rule := range rules {
				for _, cond := range rule.Conditions {
					if cond.Operator == OperatorRegex {
						if _, err := regexp.Compile(cond.Value); err != nil {
							return fmt.Errorf("invalid regex in rule %s: %w", rule.ID, err)
						}
					}
				}
			}
		}
	}

	if e.RolloutRules != nil {
		if rawRules, ok := e.RolloutRules["rules"]; ok {
			rBytes, err := json.Marshal(rawRules)
			if err == nil {
				var rollouts []RolloutRule
				if err := json.Unmarshal(rBytes, &rollouts); err == nil && len(rollouts) > 0 {
					total := 0
					for _, r := range rollouts {
						total += r.Percentage
					}
					if total != 10000 {
						return fmt.Errorf("rollout percentages must sum to 10,000 basis points (got %d)", total)
					}
				}
			}
		}
	}

	return nil
}
