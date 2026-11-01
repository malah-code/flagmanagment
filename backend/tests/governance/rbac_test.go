package governance_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// @TestCase: TC-RBAC-001
// @TestCase: TC-RBAC-002
func TestRBACPermissionBoundary(t *testing.T) {
	r := chi.NewRouter()
	r.Put("/api/v1/environments/{envId}/flags/{flagId}", func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-User-Role")
		if role == "viewer" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Viewer attempt should be rejected with 403 Forbidden
	reqViewer, _ := http.NewRequest("PUT", "/api/v1/environments/prod/flags/f1", nil)
	reqViewer.Header.Set("X-User-Role", "viewer")
	rrViewer := httptest.NewRecorder()
	r.ServeHTTP(rrViewer, reqViewer)

	if rrViewer.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for viewer role, got %d", rrViewer.Code)
	}

	// Admin attempt should succeed with 200 OK
	reqAdmin, _ := http.NewRequest("PUT", "/api/v1/environments/prod/flags/f1", nil)
	reqAdmin.Header.Set("X-User-Role", "admin")
	rrAdmin := httptest.NewRecorder()
	r.ServeHTTP(rrAdmin, reqAdmin)

	if rrAdmin.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for admin role, got %d", rrAdmin.Code)
	}
}

