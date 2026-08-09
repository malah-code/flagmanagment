package sdk

import (
	"testing"
	"github.com/open-feature/go-sdk/openfeature"
)

func TestBucketUser(t *testing.T) {
	// Simple test to check if bucketing is consistent
	bucket1 := bucketUser("user123")
	bucket2 := bucketUser("user123")
	if bucket1 != bucket2 {
		t.Errorf("Expected consistent bucketing for the same key, got %d and %d", bucket1, bucket2)
	}
}

func TestEvaluateFlag(t *testing.T) {
	flag := map[string]interface{}{
		"key":            "test-flag",
		"enabled":        true,
		"type":           "BOOLEAN",
		"defaultVariant": "off",
		"variants": map[string]interface{}{
			"on":  true,
			"off": false,
		},
		"rules": []interface{}{
			map[string]interface{}{
				"rollout": map[string]interface{}{
					"on":  float64(50),
					"off": float64(50),
				},
			},
		},
	}

	evalContext := openfeature.NewEvaluationContext("user123", nil)

	val, variant, reason, err := evaluateFlag(flag, evalContext)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if reason != openfeature.TargetingMatchReason {
		t.Errorf("Expected TargetingMatchReason, got %v", reason)
	}

	// Murmur3 hash of "user123" % 100 will fall either into "on" or "off"
	if variant != "on" && variant != "off" {
		t.Errorf("Expected variant 'on' or 'off', got %v", variant)
	}
	if _, ok := val.(bool); !ok {
		t.Errorf("Expected boolean value, got %T", val)
	}
}
