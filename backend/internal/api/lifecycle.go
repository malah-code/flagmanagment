package api

import (
	"encoding/json"
	"net/http"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type LifecycleHandler struct {
	store        repository.Store
	rbac         *RBACMiddleware
	auditService *services.AuditService
}

func NewLifecycleHandler(store repository.Store, rbac *RBACMiddleware, auditService *services.AuditService) *LifecycleHandler {
	return &LifecycleHandler{
		store:        store,
		rbac:         rbac,
		auditService: auditService,
	}
}

type TransitionLifecycleRequest struct {
	Action string `json:"action"` // ARCHIVE, DEPRECATE, RESTORE, MARK_STALE
}

func (h *LifecycleHandler) RegisterRoutes(r chi.Router) {
	r.Post("/environments/{envId}/flags/{flagId}/lifecycle", h.TransitionLifecycle)
}

func (h *LifecycleHandler) TransitionLifecycle(w http.ResponseWriter, r *http.Request) {
	envIDStr := chi.URLParam(r, "envId")
	envID, err := uuid.Parse(envIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}

	flagIDStr := chi.URLParam(r, "flagId")
	flagID, err := uuid.Parse(flagIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid flag ID")
		return
	}

	var req TransitionLifecycleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var targetState models.FlagLifecycleState
	switch req.Action {
	case "ARCHIVE":
		targetState = models.LifecycleArchived
	case "DEPRECATE":
		targetState = models.LifecycleDeprecated
	case "RESTORE":
		targetState = models.LifecycleActive
	case "MARK_STALE":
		targetState = models.LifecycleStale
	default:
		RespondWithError(w, http.StatusBadRequest, "Invalid lifecycle action. Allowed: ARCHIVE, DEPRECATE, RESTORE, MARK_STALE")
		return
	}

	state, err := h.store.FlagStateRepo().GetByEnvAndFlag(r.Context(), envID, flagID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Flag state not found")
		return
	}

	if err := h.store.FlagStateRepo().UpdateLifecycleState(r.Context(), state.ID, targetState); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to update lifecycle state")
		return
	}

	// Audit logging
	actorID := models.SystemActor()
	_ = h.auditService.LogAction(r.Context(), &models.AuditLog{
		ID:            uuid.New(),
		ProjectID:     &envID, // Used envID or projectID if available
		EnvironmentID: &envID,
		ActorID:       actorID,
		Action:        "LIFECYCLE_TRANSITION",
		TargetType:    "FEATURE_FLAG",
		TargetID:      flagID,
	})

	state.LifecycleState = targetState
	RespondWithJSON(w, http.StatusOK, state)
}
