package sdk

import (
	"testing"
	"time"

	"github.com/flagmanagment/backend/internal/models"
)

// @TestCase: TC-PERF-001
// @TestCase: TC-SDK-001
func BenchmarkEvaluateFlag(b *testing.B) {
	rulesJSON := models.JSONB{
		"rules": []map[string]interface{}{
			{
				"id": "perf-rule-1",
				"conditions": []map[string]interface{}{
					{
						"attribute": "email",
						"operator":  "CONTAINS",
						"value":     "@example.com",
					},
					{
						"attribute": "tier",
						"operator":  "EQUALS",
						"value":     "enterprise",
					},
				},
				"variation": true,
			},
		},
	}

	evalContext := EvaluationContext{
		"email": "lead-dev@example.com",
		"tier":  "enterprise",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = EvaluateContextualRules(rulesJSON, evalContext)
	}
}

// TestSubMillisecondEvaluationLatency enforces Constitution Principle IV (< 1.0 ms)
func TestSubMillisecondEvaluationLatency(t *testing.T) {
	rulesJSON := models.JSONB{
		"rules": []map[string]interface{}{
			{
				"id": "rule-1",
				"conditions": []map[string]interface{}{
					{"attribute": "country", "operator": "EQUALS", "value": "US"},
				},
				"variation": true,
			},
		},
	}
	ctx := EvaluationContext{"country": "US"}

	// Warmup
	for i := 0; i < 1000; i++ {
		_, _ = EvaluateContextualRules(rulesJSON, ctx)
	}

	// Measure 10,000 evaluations
	const iterations = 10000
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, _ = EvaluateContextualRules(rulesJSON, ctx)
	}
	elapsed := time.Since(start)

	avgDuration := elapsed / iterations
	t.Logf("Average evaluation duration over %d runs: %v", iterations, avgDuration)

	// Constitution Principle IV: MUST be under 1 millisecond (1,000,000 ns)
	if avgDuration > time.Millisecond {
		t.Fatalf("Constitution Principle IV violated: evaluation took %v, expected < 1ms", avgDuration)
	}
}
