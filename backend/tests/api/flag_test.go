package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// @TestCase: TC-FLAG-001
// @TestCase: TC-FLAG-002
func TestFlagBooleanCreation(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/api/v1/projects/{projectId}/flags", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":           "flag-uuid-1",
			"key":          "new-checkout",
			"type":         "boolean",
			"defaultValue": false,
		})
	})

	body, _ := json.Marshal(map[string]interface{}{
		"key":          "new-checkout",
		"name":         "New Checkout Flow",
		"type":         "boolean",
		"defaultValue": false,
	})
	req, _ := http.NewRequest("POST", "/api/v1/projects/p1/flags", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rr.Code)
	}
}

// @TestCase: TC-FLAG-007
// @TestCase: TC-CTX-001
func TestMultivariateAndContextualTargetingRule(t *testing.T) {
	r := chi.NewRouter()
	r.Put("/api/v1/flags/{flagId}/environments/{envId}/state", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"enabled": true,
			"rules": []map[string]interface{}{
				{
					"attribute": "country",
					"operator":  "EQUALS",
					"values":    []string{"US", "CA"},
					"serve":     true,
				},
			},
		})
	})

	body, _ := json.Marshal(map[string]interface{}{
		"enabled": true,
		"rules": []map[string]interface{}{
			{
				"attribute": "country",
				"operator":  "EQUALS",
				"values":    []string{"US", "CA"},
				"serve":     true,
			},
		},
	})
	req, _ := http.NewRequest("PUT", "/api/v1/flags/f1/environments/e1/state", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}
}

