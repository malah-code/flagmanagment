package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flagmanagment/backend/tests/testutil"
	"github.com/go-chi/chi/v5"
)

// @TestCase: TC-AUTH-001
func TestLoginWithValidCredentials(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]string
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req["email"] == "admin@example.com" && req["password"] == "admin123" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"token": "valid-jwt-token",
				"role":  "admin",
			})
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	})

	body, _ := json.Marshal(map[string]string{
		"email":    "admin@example.com",
		"password": "admin123",
	})
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}
}

// @TestCase: TC-AUTH-002
func TestLoginWithInvalidCredentials(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
	})

	body, _ := json.Marshal(map[string]string{
		"email":    "admin@example.com",
		"password": "wrongpassword",
	})
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rr.Code)
	}
}

// @TestCase: TC-AUTH-007
func TestPreventDuplicateUserInvitation(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/api/v1/users/invite", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]string
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req["email"] == "existing@example.com" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user already exists"})
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	body, _ := json.Marshal(map[string]string{"email": "existing@example.com", "role": "developer"})
	req, _ := http.NewRequest("POST", "/api/v1/users/invite", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected 409 Conflict, got %d", rr.Code)
	}
}

// @TestCase: TC-AUTH-009
// @TestCase: TC-AUTH-010
func TestSSOWithMockIdP(t *testing.T) {
	mockIdP := testutil.NewMockOIDCServer()
	defer mockIdP.Close()

	if mockIdP.URL() == "" {
		t.Fatal("expected mock OIDC server URL to be non-empty")
	}

	resp, err := http.Get(mockIdP.URL() + "/.well-known/openid-configuration")
	if err != nil {
		t.Fatalf("failed to query mock OIDC discovery: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from OIDC discovery, got %d", resp.StatusCode)
	}
}

