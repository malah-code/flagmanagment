package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
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

func (h *WebhookHandler) RegisterManagementRoutes(r chi.Router) {
	r.Post("/projects/{id}/webhooks", h.CreateWebhook)
	r.Get("/projects/{id}/webhooks", h.ListWebhooks)
}

func (h *WebhookHandler) RegisterAPMRoutes(r chi.Router) {
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

func (h *WebhookHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var payload struct {
		URL       string          `json:"url"`
		SecretKey *string         `json:"secret_key,omitempty"`
		Events    json.RawMessage `json:"events"`
		IsActive  bool            `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Bad Request: malformed JSON")
		return
	}

	wh := &models.WebhookIntegration{
		ProjectID: projectID,
		URL:       payload.URL,
		SecretKey: payload.SecretKey,
		Events:    payload.Events,
		IsActive:  payload.IsActive,
	}

	if err := h.webhookService.CreateWebhook(r.Context(), wh); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create webhook")
		return
	}

	RespondWithJSON(w, http.StatusCreated, wh)
}

func (h *WebhookHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	list, err := h.webhookService.ListWebhooks(r.Context(), projectID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to list webhooks")
		return
	}

	RespondWithJSON(w, http.StatusOK, list)
}
