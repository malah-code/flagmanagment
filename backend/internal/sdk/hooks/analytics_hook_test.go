package hooks

import (
	"bytes"
	"strings"
	"testing"

	"github.com/flagmanagment/backend/internal/models"
)

func TestAnalyticsHook_After(t *testing.T) {
	var buf bytes.Buffer
	hook := NewAnalyticsHook("TEST", &buf)

	ctx := HookContext{
		FlagKey:  "my-flag",
		FlagType: "boolean",
		EvaluationContext: &models.EvaluationContext{
			EntityKey: "user-99",
			Attributes: map[string]interface{}{
				"plan": "pro",
			},
		},
	}

	details := EvaluationDetails{
		FlagKey: "my-flag",
		Value:   "true",
		Reason:  "TARGETING_MATCH",
	}

	hook.After(ctx, details)

	output := buf.String()

	if !strings.Contains(output, "[TEST] ANALYTICS_EVENT:") {
		t.Errorf("Expected output to contain prefix, got: %s", output)
	}
	if !strings.Contains(output, `"flag_key":"my-flag"`) {
		t.Errorf("Expected output to contain flag_key, got: %s", output)
	}
	if !strings.Contains(output, `"entity_key":"user-99"`) {
		t.Errorf("Expected output to contain entity_key, got: %s", output)
	}
}

func TestAnalyticsHook_Error(t *testing.T) {
	var buf bytes.Buffer
	hook := NewAnalyticsHook("TEST", &buf)

	ctx := HookContext{
		FlagKey: "error-flag",
	}

	hook.Error(ctx, nil)

	output := buf.String()
	if !strings.Contains(output, "[TEST] ERROR_EVENT: flag_key=error-flag") {
		t.Errorf("Expected error log output, got: %s", output)
	}
}
