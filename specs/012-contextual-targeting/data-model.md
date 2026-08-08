# Data Model & Contracts: Contextual Targeting Engine

## Data Model Updates

No new database tables are required. We will define Go structs that map to the `TargetingRules` JSONB column in `environment_flag_states`.

### `TargetingRule` Structs (Go)

```go
package models

// Operator defines the type of comparison for a targeting condition
type Operator string

const (
    OperatorEquals   Operator = "EQUALS"
    OperatorContains Operator = "CONTAINS"
    OperatorRegex    Operator = "REGEX"
)

// TargetingCondition represents a single AND condition within a rule
type TargetingCondition struct {
    Attribute string   `json:"attribute"`
    Operator  Operator `json:"operator"`
    Value     string   `json:"value"`
}

// TargetingRule represents a set of conditions that must ALL be true (AND).
// If the rule evaluates to true, the specified Variation is served.
type TargetingRule struct {
    ID         string               `json:"id"`
    Conditions []TargetingCondition `json:"conditions"`
    Variation  bool                 `json:"variation"` // Currently boolean, can be extended for multivariate
}

// TargetingRulesPayload represents the entire targeting rules payload for a flag state
type TargetingRulesPayload struct {
    Rules []TargetingRule `json:"rules"` // Evaluated sequentially (OR logic)
}
```

## API Contracts

### Update Environment Flag State (Existing Endpoint)
`PUT /api/v1/projects/{projectId}/environments/{environmentId}/flags/{flagId}`

The `targeting_rules` field in the payload will now accept a structured JSON object matching the `TargetingRulesPayload` struct.

**Example Payload:**
```json
{
  "enabled": true,
  "targeting_rules": {
    "rules": [
      {
        "id": "rule-1",
        "conditions": [
          {
            "attribute": "email",
            "operator": "REGEX",
            "value": ".*@internal\\.com$"
          },
          {
            "attribute": "tenant",
            "operator": "EQUALS",
            "value": "tenant_abc"
          }
        ],
        "variation": true
      }
    ]
  }
}
```

## SDK Evaluation Interface

The SDK will require a context object to be passed during evaluation:

```go
// EvaluationContext represents the user context passed to the SDK
type EvaluationContext map[string]string

// SDK Interface
func (s *SDKClient) EvaluateBooleanFlag(ctx context.Context, flagKey string, evalCtx EvaluationContext, defaultValue bool) bool
```
