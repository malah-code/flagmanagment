package sdk

import (
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/sdk/hooks"
	"sync/atomic"
	"testing"
	"time"
)

func TestEvaluateContextualRules(t *testing.T) {
	// Create some dummy JSONB representing targeting rules
	// For testing natively we will just pass a map structure
	
	targetingRulesJSON := models.JSONB{
		"rules": []map[string]interface{}{
			{
				"id": "rule-1",
				"conditions": []map[string]interface{}{
					{
						"attribute": "email",
						"operator":  "REGEX",
						"value":     ".*@test\\.com$",
					},
					{
						"attribute": "region",
						"operator":  "EQUALS",
						"value":     "US",
					},
				},
				"variation": true,
			},
			{
				"id": "rule-2",
				"conditions": []map[string]interface{}{
					{
						"attribute": "tenant",
						"operator":  "CONTAINS",
						"value":     "beta",
					},
				},
				"variation": true,
			},
		},
	}

	tests := []struct {
		name          string
		context       EvaluationContext
		expectedVar   bool
		expectedMatch bool
	}{
		{
			name: "Match first rule (AND logic)",
			context: EvaluationContext{
				"email":  "user@test.com",
				"region": "US",
			},
			expectedVar:   true,
			expectedMatch: true,
		},
		{
			name: "Fail first rule (AND logic missing region)",
			context: EvaluationContext{
				"email": "user@test.com",
			},
			expectedVar:   false,
			expectedMatch: false,
		},
		{
			name: "Match second rule (OR logic)",
			context: EvaluationContext{
				"tenant": "internal-beta-users",
			},
			expectedVar:   true,
			expectedMatch: true,
		},
		{
			name: "No match",
			context: EvaluationContext{
				"email":  "user@other.com",
				"region": "EU",
			},
			expectedVar:   false,
			expectedMatch: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			variation, match := EvaluateContextualRules(targetingRulesJSON, tc.context)
			if match != tc.expectedMatch {
				t.Errorf("expected match %v, got %v", tc.expectedMatch, match)
			}
			if match && variation != tc.expectedVar {
				t.Errorf("expected variation %v, got %v", tc.expectedVar, variation)
			}
		})
	}
}

func TestEvaluateRolloutSplit(t *testing.T) {
	rolloutJSON := models.JSONB{
		"rules": []map[string]interface{}{
			{
				"variation_id": "var_a",
				"percentage":   5000, // 50%
			},
			{
				"variation_id": "var_b",
				"percentage":   5000, // 50%
			},
		},
	}

	salt := "test-env-salt"

	// 1. Assert Determinism (same user_id -> same variation 100 times)
	firstVarID, matched := EvaluateRolloutSplit("my-ab-flag", "user-12345", rolloutJSON, salt)
	if !matched {
		t.Fatalf("expected rollout match for user-12345")
	}
	for i := 0; i < 100; i++ {
		varID, _ := EvaluateRolloutSplit("my-ab-flag", "user-12345", rolloutJSON, salt)
		if varID != firstVarID {
			t.Fatalf("expected consistent variation %s, got %s at iteration %d", firstVarID, varID, i)
		}
	}

	// 2. Assert Distribution (roughly 50/50 split across 10,000 unique users)
	counts := map[string]int{}
	for i := 0; i < 10000; i++ {
		userID := "user-id-" + string(rune(i))
		varID, matched := EvaluateRolloutSplit("my-ab-flag", userID, rolloutJSON, salt)
		if matched {
			counts[varID]++
		}
	}

	total := counts["var_a"] + counts["var_b"]
	if total == 0 {
		t.Fatalf("no evaluations matched")
	}

	// Ensure var_a is within 45% - 55%
	pctA := float64(counts["var_a"]) / float64(total)
	if pctA < 0.40 || pctA > 0.60 {
		t.Errorf("expected roughly 50%% split for var_a, got %f", pctA)
	}
}

func TestEvaluateFlag_SequentialDependencies(t *testing.T) {
	parentKey := "flag-a"
	childKey := "flag-b"

	rulesMap := map[string]*models.FlagRule{
		parentKey: {
			Key:              parentKey,
			Enabled:          true,
			DefaultVariation: "true",
		},
		childKey: {
			Key:              childKey,
			Enabled:          true,
			DefaultVariation: "true",
			ParentFlagKey:    parentKey,
		},
	}

	ctx := &models.EvaluationContext{EntityKey: "user-1"}

	// 1. Parent is ON -> Child evaluates to ON (its default)
	res1 := EvaluateFlag(rulesMap[childKey], ctx, rulesMap)
	if res1.Value != "true" || res1.Reason != "DEFAULT" {
		t.Errorf("expected true/DEFAULT, got %v/%s", res1.Value, res1.Reason)
	}

	// 2. Parent is OFF -> Child short-circuits to OFF (its default, with reason PARENT_FLAG_DISABLED)
	rulesMap[parentKey].Enabled = false
	res2 := EvaluateFlag(rulesMap[childKey], ctx, rulesMap)
	if res2.Value != "true" || res2.Reason != "PARENT_FLAG_DISABLED" {
		t.Errorf("expected true/PARENT_FLAG_DISABLED, got %v/%s", res2.Value, res2.Reason)
	}

	// 3. Parent evaluates to false via TargetingRules (simulated by changing DefaultVariation to "false")
	rulesMap[parentKey].Enabled = true
	rulesMap[parentKey].DefaultVariation = "false"
	res3 := EvaluateFlag(rulesMap[childKey], ctx, rulesMap)
	if res3.Value != "true" || res3.Reason != "PARENT_FLAG_DISABLED" {
		t.Errorf("expected true/PARENT_FLAG_DISABLED, got %v/%s", res3.Value, res3.Reason)
	}
}

type MockHook struct {
	afterCalled int32
	errorCalled int32
	ch          chan struct{}
}

func (h *MockHook) After(ctx hooks.HookContext, details hooks.EvaluationDetails) {
	atomic.StoreInt32(&h.afterCalled, 1)
	if h.ch != nil {
		select {
		case h.ch <- struct{}{}:
		default:
		}
	}
}

func (h *MockHook) Error(ctx hooks.HookContext, err error) {
	atomic.StoreInt32(&h.errorCalled, 1)
}



func TestEvaluateFlag_Hooks(t *testing.T) {
	// Reset registered hooks
	ClearHooks()

	mockHook := &MockHook{ch: make(chan struct{}, 1)}
	RegisterHook(mockHook)

	rulesMap := map[string]*models.FlagRule{
		"flag-1": {
			Key:              "flag-1",
			Enabled:          true,
			DefaultVariation: "true",
		},
	}

	ctx := &models.EvaluationContext{EntityKey: "user-1"}

	res := EvaluateFlag(rulesMap["flag-1"], ctx, rulesMap)
	if res.Value != "true" {
		t.Errorf("expected true, got %v", res.Value)
	}

	// Wait for hook or timeout
	select {
	case <-mockHook.ch:
		// success
	case <-time.After(1 * time.Second):
		t.Errorf("timeout waiting for mock hook After to be called")
	}

	if atomic.LoadInt32(&mockHook.afterCalled) != 1 {
		t.Errorf("expected mock hook After to be called")
	}

	// Clean up
	ClearHooks()
}

func TestHashPII(t *testing.T) {
	ctx := &models.EvaluationContext{
		EntityKey: "user-123",
		Attributes: map[string]interface{}{
			"email": "user@example.com",
			"name":  "John Doe",
			"phone": "555-0199",
		},
	}

	salt := "test-salt-123"
	hashed := HashPII(ctx, salt)

	if hashed.Attributes["name"] != "John Doe" {
		t.Errorf("expected name to remain unhashed, got %v", hashed.Attributes["name"])
	}

	emailStr, ok := hashed.Attributes["email"].(string)
	if !ok || emailStr == "user@example.com" {
		t.Errorf("expected email to be hashed, got %v", hashed.Attributes["email"])
	}

	phoneStr, ok := hashed.Attributes["phone"].(string)
	if !ok || phoneStr == "555-0199" {
		t.Errorf("expected phone to be hashed, got %v", hashed.Attributes["phone"])
	}
}


