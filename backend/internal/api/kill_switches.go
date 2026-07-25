package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/go-playground/validator/v10"
)

type KillSwitchHandler struct {
	store    repository.Store
	rbac     *RBACMiddleware
	validate *validator.Validate
}

func NewKillSwitchHandler(store repository.Store, rbac *RBACMiddleware) *KillSwitchHandler {
	return &KillSwitchHandler{
		store:    store,
		rbac:     rbac,
		validate: validator.New(),
	}
}

func (h *KillSwitchHandler) RegisterRoutes(r chi.Router) {
	r.With(h.rbac.RequireRole("VIEWER")).Get("/environments/{envId}/flags/{flagId}/kill-switches", h.List)
	r.With(h.rbac.RequireRole("EDITOR")).Post("/environments/{envId}/flags/{flagId}/kill-switches", h.Create)
	r.With(h.rbac.RequireRole("EDITOR")).Delete("/environments/{envId}/flags/{flagId}/kill-switches/{id}", h.Delete)
}

type CreateKillSwitchRequest struct {
	AlertIdentifier string `json:"alert_identifier" validate:"required"`
}

func (h *KillSwitchHandler) List(w http.ResponseWriter, r *http.Request) {
	envIDStr := chi.URLParam(r, "envId")
	flagIDStr := chi.URLParam(r, "flagId")

	envID, _ := uuid.Parse(envIDStr)
	flagID, _ := uuid.Parse(flagIDStr)

	rules, err := h.store.KillSwitchRepo().ListByEnvironmentAndFlag(r.Context(), envID, flagID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch kill switches")
		return
	}

	RespondWithJSON(w, http.StatusOK, rules)
}

func (h *KillSwitchHandler) Create(w http.ResponseWriter, r *http.Request) {
	envIDStr := chi.URLParam(r, "envId")
	flagIDStr := chi.URLParam(r, "flagId")

	envID, _ := uuid.Parse(envIDStr)
	flagID, _ := uuid.Parse(flagIDStr)

	var req CreateKillSwitchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ks := &models.KillSwitchRule{
		ID:              uuid.New(),
		FlagID:          flagID,
		EnvironmentID:   envID,
		AlertIdentifier: req.AlertIdentifier,
		Action:          "DISABLE",
	}

	if err := h.store.KillSwitchRepo().Create(r.Context(), ks); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create kill switch")
		return
	}

	RespondWithJSON(w, http.StatusCreated, ks)
}

func (h *KillSwitchHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := uuid.Parse(idStr)

	if err := h.store.KillSwitchRepo().Delete(r.Context(), id); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete kill switch")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
