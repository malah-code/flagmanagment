package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
)

type SlackConfigHandler struct {
	store repository.Store
	rbac  *RBACMiddleware
}

func NewSlackConfigHandler(store repository.Store, rbac *RBACMiddleware) *SlackConfigHandler {
	return &SlackConfigHandler{
		store: store,
		rbac:  rbac,
	}
}

func (h *SlackConfigHandler) RegisterRoutes(r chi.Router) {
	r.With(h.rbac.RequireRole("VIEWER")).Get("/environments/{envId}/slack", h.GetConfig)
	r.With(h.rbac.RequireRole("EDITOR")).Post("/environments/{envId}/slack", h.UpsertConfig)
	r.With(h.rbac.RequireRole("EDITOR")).Delete("/environments/{envId}/slack", h.DeleteConfig)
}

type SlackConfigRequest struct {
	WebhookURL string `json:"webhook_url"`
	Enabled    bool   `json:"enabled"`
}

func (h *SlackConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	envIDStr := chi.URLParam(r, "envId")
	envID, err := uuid.Parse(envIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}

	config, err := h.store.SlackConfigRepo().GetByEnvironmentID(r.Context(), envID)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithJSON(w, http.StatusOK, map[string]interface{}{
				"webhook_url": "",
				"enabled":     false,
			})
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch Slack configuration")
		return
	}

	RespondWithJSON(w, http.StatusOK, config)
}

func (h *SlackConfigHandler) UpsertConfig(w http.ResponseWriter, r *http.Request) {
	envIDStr := chi.URLParam(r, "envId")
	envID, err := uuid.Parse(envIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}

	var req SlackConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	config := &models.SlackWebhookConfig{
		EnvironmentID: envID,
		WebhookURL:    req.WebhookURL,
		Enabled:       req.Enabled,
	}

	if err := h.store.SlackConfigRepo().Upsert(r.Context(), config); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to save Slack configuration")
		return
	}

	RespondWithJSON(w, http.StatusOK, config)
}

func (h *SlackConfigHandler) DeleteConfig(w http.ResponseWriter, r *http.Request) {
	envIDStr := chi.URLParam(r, "envId")
	envID, err := uuid.Parse(envIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}

	if err := h.store.SlackConfigRepo().Delete(r.Context(), envID); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete Slack configuration")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
