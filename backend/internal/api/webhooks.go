package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/flagmanagment/backend/internal/services"
)

type WebhookHandler struct {
	webhookService services.WebhookService
}

func NewWebhookHandler(webhookService services.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
	}
}

func (h *WebhookHandler) RegisterRoutes(r chi.Router) {
	r.Post("/apm", h.HandleAPMWebhook)
}

type APMWebhookPayload struct {
	AlertIdentifier string `json:"alert_identifier"`
	Status          string `json:"status"`
	Description     string `json:"description"`
}

func (h *WebhookHandler) HandleAPMWebhook(w http.ResponseWriter, r *http.Request) {
	env := GetEnvironmentFromContext(r.Context())
	if env == nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized: missing or invalid environment context")
		return
	}

	var payload APMWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Bad Request: malformed JSON")
		return
	}

	if payload.AlertIdentifier == "" {
		RespondWithError(w, http.StatusBadRequest, "Bad Request: missing alert_identifier")
		return
	}

	flagsKilled, err := h.webhookService.ProcessAPMAlert(r.Context(), env.ID, payload.AlertIdentifier, payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "processed",
		"flags_killed": flagsKilled,
	})
}
