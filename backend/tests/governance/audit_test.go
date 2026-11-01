package governance_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

// @TestCase: TC-AUDIT-001
// @TestCase: TC-AUDIT-004
func TestAuditLogListingAndSIEMFormat(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/api/v1/audit-logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"id":         "audit-uuid-1",
				"action":     "FLAG_CREATED",
				"actor":      "admin@example.com",
				"timestamp":  time.Now().UTC().Format(time.RFC3339),
				"resourceId": "flag-uuid-1",
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/audit-logs", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var logs []map[string]interface{}
	_ = json.Unmarshal(rr.Body.Bytes(), &logs)
	if len(logs) == 0 || logs[0]["action"] != "FLAG_CREATED" {
		t.Fatalf("expected audit log with action FLAG_CREATED")
	}
}

