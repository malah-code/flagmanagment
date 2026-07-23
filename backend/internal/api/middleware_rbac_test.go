package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flagmanagment/backend/internal/api"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// MockStore satisfies repository.Store for testing.
type mockStore struct {
	repository.Store
	mockRoleRepo repository.RoleRepository
	mockAuditRepo repository.AuditRepository
}

func (m *mockStore) RoleRepo() repository.RoleRepository {
	return m.mockRoleRepo
}

func (m *mockStore) AuditRepo() repository.AuditRepository {
	return m.mockAuditRepo
}

// MockRoleRepository satisfies repository.RoleRepository for testing.
type mockRoleRepository struct {
	repository.RoleRepository
	userRoles []*models.UserRole
	err       error
}

func (m *mockRoleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID, projectID *uuid.UUID) ([]*models.UserRole, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.userRoles, nil
}

func TestRBACMiddleware_RequireRole(t *testing.T) {
	testUserID := uuid.New()

	tests := []struct {
		name           string
		userRole       string
		requiredRole   string
		expectedStatus int
	}{
		{
			name:           "Admin allowed on Admin route",
			userRole:       "ADMIN",
			requiredRole:   "ADMIN",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Admin allowed on Editor route",
			userRole:       "ADMIN",
			requiredRole:   "EDITOR",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Editor allowed on Editor route",
			userRole:       "EDITOR",
			requiredRole:   "EDITOR",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Editor allowed on Viewer route",
			userRole:       "EDITOR",
			requiredRole:   "VIEWER",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Viewer forbidden on Editor route",
			userRole:       "VIEWER",
			requiredRole:   "EDITOR",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Viewer allowed on Viewer route",
			userRole:       "VIEWER",
			requiredRole:   "VIEWER",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockStore{
				mockRoleRepo: &mockRoleRepository{
					userRoles: []*models.UserRole{
						{
							UserID: testUserID,
							Role: &models.Role{
								Name: tt.userRole,
							},
						},
					},
				},
			}

			rbac := api.NewRBACMiddleware(store)

			r := chi.NewRouter()
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					ctx := context.WithValue(req.Context(), api.UserIDKey, testUserID.String())
					next.ServeHTTP(w, req.WithContext(ctx))
				})
			})

			r.With(rbac.RequireRole(tt.requiredRole)).Get("/test", func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestRBACMiddleware_Unauthorized(t *testing.T) {
	store := &mockStore{
		mockRoleRepo: &mockRoleRepository{},
	}
	rbac := api.NewRBACMiddleware(store)

	r := chi.NewRouter()
	r.With(rbac.RequireRole("VIEWER")).Get("/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized for missing user ID, got %d", rr.Code)
	}
}
