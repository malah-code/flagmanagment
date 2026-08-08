package api

import (
	"encoding/json"
	"net/http"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type StalePolicyHandler struct {
	store repository.Store
	rbac  *RBACMiddleware
}

func NewStalePolicyHandler(store repository.Store, rbac *RBACMiddleware) *StalePolicyHandler {
	return &StalePolicyHandler{store: store, rbac: rbac}
}

type SetStalePolicyRequest struct {
	StaleAfterDays int `json:"stale_after_days"`
}

func (h *StalePolicyHandler) RegisterRoutes(r chi.Router) {
	r.Get("/projects/{projectId}/stale-policy", h.GetPolicy)
	r.Put("/projects/{projectId}/stale-policy", h.SetPolicy)
}

func (h *StalePolicyHandler) GetPolicy(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	policy, err := h.store.StalePolicyRepo().GetByEnvironment(r.Context(), projectID, uuid.Nil)
	if err != nil {
		if err == repository.ErrNotFound {
			// Return default 30 days
			RespondWithJSON(w, http.StatusOK, map[string]interface{}{
				"project_id":       projectID,
				"stale_after_days": 30,
			})
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch stale policy")
		return
	}
	RespondWithJSON(w, http.StatusOK, policy)
}

func (h *StalePolicyHandler) SetPolicy(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var req SetStalePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.StaleAfterDays <= 0 {
		RespondWithError(w, http.StatusBadRequest, "Invalid stale_after_days (must be > 0)")
		return
	}

	policy := &models.StaleFlagPolicy{
		ProjectID:      projectID,
		StaleAfterDays: req.StaleAfterDays,
	}

	if err := h.store.StalePolicyRepo().Upsert(r.Context(), policy); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to save policy")
		return
	}

	RespondWithJSON(w, http.StatusOK, policy)
}
