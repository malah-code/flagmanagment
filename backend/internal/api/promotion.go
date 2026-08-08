package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/flagmanagment/backend/internal/cache"
	"github.com/flagmanagment/backend/internal/dto"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PromoteFlagRequest struct {
	SourceEnvID string `json:"source_env_id" validate:"required,uuid"`
	TargetEnvID string `json:"target_env_id" validate:"required,uuid"`
}

type PromotionHandler struct {
	promotionService *services.PromotionService
	cacheClient      *cache.Client
	validate         *validator.Validate
	rbac             *RBACMiddleware
}

func NewPromotionHandler(promotionService *services.PromotionService, cacheClient *cache.Client, rbac *RBACMiddleware) *PromotionHandler {
	return &PromotionHandler{
		promotionService: promotionService,
		cacheClient:      cacheClient,
		validate:         validator.New(),
		rbac:             rbac,
	}
}

func (h *PromotionHandler) RegisterRoutes(r chi.Router) {
	r.With(h.rbac.RequireRole("EDITOR")).Post("/projects/{projectId}/flags/{flagId}/promote", h.PromoteFlag)
}

func (h *PromotionHandler) PromoteFlag(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	flagIDStr := chi.URLParam(r, "flagId")

	_, err1 := uuid.Parse(projectIDStr)
	flagID, err2 := uuid.Parse(flagIDStr)
	if err1 != nil || err2 != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req PromoteFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	sourceEnvID, _ := uuid.Parse(req.SourceEnvID)
	targetEnvID, _ := uuid.Parse(req.TargetEnvID)

	actorIDStr, _ := r.Context().Value(UserIDKey).(string)
	actorID, _ := uuid.Parse(actorIDStr)

	result, err := h.promotionService.PromoteFlag(r.Context(), flagID, sourceEnvID, targetEnvID, actorID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to promote flag")
		return
	}

	switch v := result.(type) {
	case *models.ChangeRequest:
		RespondWithJSON(w, http.StatusAccepted, v)
		return
	case *models.EnvironmentFlagState:
		if h.cacheClient != nil {
			_ = h.cacheClient.PublishRulesetUpdate(r.Context(), targetEnvID.String(), time.Now().UTC().String())
		}
		RespondWithJSON(w, http.StatusOK, dto.FlagStateResponse{
			EnvironmentID:  v.EnvironmentID.String(),
			FeatureFlagID:  v.FeatureFlagID.String(),
			Enabled:        v.Enabled,
			TargetingRules: v.TargetingRules,
			RemoteConfig:   v.RemoteConfig,
			UpdatedAt:      v.UpdatedAt,
		})
		return
	}

	RespondWithError(w, http.StatusInternalServerError, "Unknown result type")
}
