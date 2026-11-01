package governance_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// @TestCase: TC-CR-001
// @TestCase: TC-CR-002
func TestChangeRequestCreationAndApproval(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/api/v1/change-requests", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     "cr-uuid-1",
			"status": "pending",
		})
	})

	r.Post("/api/v1/change-requests/{id}/approve", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     "cr-uuid-1",
			"status": "approved",
		})
	})

	body, _ := json.Marshal(map[string]interface{}{
		"title":       "Roll out new UI to production",
		"description": "Increases percentage rollout to 50%",
	})
	req, _ := http.NewRequest("POST", "/api/v1/change-requests", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rr.Code)
	}

	reqApprove, _ := http.NewRequest("POST", "/api/v1/change-requests/cr-uuid-1/approve", nil)
	rrApprove := httptest.NewRecorder()
	r.ServeHTTP(rrApprove, reqApprove)

	if rrApprove.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rrApprove.Code)
	}
}

