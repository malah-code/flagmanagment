package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMetricsRegistryAndHandler(t *testing.T) {
	m := Init()
	if m == nil {
		t.Fatal("expected non-nil metrics instance")
	}

	// Record sample operations
	m.RecordEvaluation("proj-test", "production", "new-checkout", "variant-a", "TARGETING_MATCH", 250*time.Microsecond)
	m.RecordCacheHit()
	m.RecordCacheMiss()
	m.ActiveFlags.WithLabelValues("proj-test").Set(42)
	m.StaleFlags.WithLabelValues("proj-test").Set(3)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	m.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()

	expectedSubstrings := []string{
		"flagmanagment_evaluator_evaluations_total",
		"flagmanagment_evaluator_evaluation_duration_seconds",
		"flagmanagment_cache_operations_total",
		"flagmanagment_flags_active_total",
		"flagmanagment_flags_stale_total",
		"new-checkout",
		"variant-a",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(body, sub) {
			t.Errorf("expected metrics output to contain %q, but was missing", sub)
		}
	}
}
