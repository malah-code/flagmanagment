package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
)

type RBACMiddleware struct {
	store repository.Store
}

func NewRBACMiddleware(store repository.Store) *RBACMiddleware {
	return &RBACMiddleware{store: store}
}

func (m *RBACMiddleware) RequireRole(requiredRoleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			saIDStr := r.Context().Value(ServiceAccountIDKey)
			var userRoles []*models.UserRole // We can reuse UserRole struct for checking permissions
			var err error

			projectIDStr := r.PathValue("project_id")
			if projectIDStr == "" {
				projectIDStr = chi.URLParam(r, "id")
			}

			var projectID *uuid.UUID
			if projectIDStr != "" {
				pid, pidErr := uuid.Parse(projectIDStr)
				if pidErr == nil {
					projectID = &pid
				}
			}

			if saIDStr != nil {
				saID, _ := saIDStr.(uuid.UUID)
				// We need a way to get SA roles. For now, we will add GetServiceAccountRoles to RoleRepo
				saRoles, repoErr := m.store.RoleRepo().GetServiceAccountRoles(r.Context(), saID, projectID)
				if repoErr != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				for _, sr := range saRoles {
					userRoles = append(userRoles, &models.UserRole{Role: sr.Role})
				}
			} else {
				userIDStr := r.Context().Value(UserIDKey)
				if userIDStr == nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				userID, parseErr := uuid.Parse(userIDStr.(string))
				if parseErr != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				userRoles, err = m.store.RoleRepo().GetUserRoles(r.Context(), userID, projectID)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			hasAccess := false
			for _, ur := range userRoles {
				roleName := ""
				if ur.Role != nil {
					roleName = ur.Role.Name
				}
				
				if hasPermission(roleName, requiredRoleName) {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func hasPermission(userRole, requiredRole string) bool {
	if userRole == "ADMIN" {
		return true
	}
	if userRole == "RELEASE_MANAGER" && (requiredRole == "RELEASE_MANAGER" || requiredRole == "EDITOR" || requiredRole == "VIEWER") {
		return true
	}
	if userRole == "EDITOR" && (requiredRole == "EDITOR" || requiredRole == "VIEWER") {
		return true
	}
	if userRole == "VIEWER" && requiredRole == "VIEWER" {
		return true
	}
	return false
}
