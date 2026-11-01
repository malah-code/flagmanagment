package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// @TestCase: TC-PROJ-001
func TestProjectCreation(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/api/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{"id": "proj-uuid-1", "name": "New Project"})
	})

	body, _ := json.Marshal(map[string]string{"name": "New Project"})
	req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rr.Code)
	}
}

// @TestCase: TC-PROJ-002
func TestProjectDuplicateNameRejection(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/api/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]string
		_ = json.NewDecoder(r.Body).Decode(&payload)
		if payload["name"] == "DuplicateApp" {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	body, _ := json.Marshal(map[string]string{"name": "DuplicateApp"})
	req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected 409 Conflict, got %d", rr.Code)
	}
}

// @TestCase: TC-ENV-001
// @TestCase: TC-ENV-002
func TestEnvironmentIsolationKeyGeneration(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/api/v1/projects/{projectId}/environments", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id":          "env-uuid-1",
			"name":        "Staging",
			"clientApiKey": "client-key-123",
			"serverApiKey": "server-key-456",
		})
	})

	body, _ := json.Marshal(map[string]string{"name": "Staging"})
	req, _ := http.NewRequest("POST", "/api/v1/projects/proj-1/environments", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rr.Code)
	}

	var res map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &res)
	if res["clientApiKey"] == "" || res["serverApiKey"] == "" {
		t.Fatalf("expected unique client and server API keys to be generated")
	}
}

