package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/flagmanagment/backend/internal/cache"
	"github.com/flagmanagment/backend/internal/dto"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type FlagHandler struct {
	store       repository.Store
	cacheClient *cache.Client
	validate    *validator.Validate
	rbac        *RBACMiddleware
	audit       *services.AuditService
}

func NewFlagHandler(store repository.Store, cacheClient *cache.Client, rbac *RBACMiddleware, audit *services.AuditService) *FlagHandler {
	return &FlagHandler{
		store:       store,
		cacheClient: cacheClient,
		validate:    validator.New(),
		rbac:        rbac,
		audit:       audit,
	}
}

func (h *FlagHandler) RegisterRoutes(r chi.Router) {
	// Global flag definitions per project
	r.With(h.rbac.RequireRole("EDITOR")).Post("/projects/{projectId}/flags", h.CreateFlag)
	r.With(h.rbac.RequireRole("VIEWER")).Get("/projects/{projectId}/flags", h.ListFlags)

	// State specific to environment (envId route doesn't have projectId, might need to rely on env's projectId internally or apply general editor role)
	// For now, applying editor level at global scope
	r.With(h.rbac.RequireRole("VIEWER")).Get("/environments/{envId}/flags/{flagId}/state", h.GetFlagState)
	r.With(h.rbac.RequireRole("EDITOR")).Put("/environments/{envId}/flags/{flagId}/state", h.UpdateFlagState)
}

func (h *FlagHandler) CreateFlag(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var req dto.CreateFeatureFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var parentFlagID *uuid.UUID
	if req.ParentFlagID != nil {
		id, err := uuid.Parse(*req.ParentFlagID)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid parent flag ID")
			return
		}
		parentFlagID = &id
	}

	now := time.Now().UTC()
	flag := &models.FeatureFlag{
		ID:           uuid.New(),
		ProjectID:    projectID,
		Key:          req.Key,
		Name:         req.Name,
		Description:  req.Description,
		Type:         models.FlagType(req.Type),
		ParentFlagID: parentFlagID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.store.FlagRepo().Create(r.Context(), flag); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create flag")
		return
	}

	actorIDStr := r.Context().Value(UserIDKey).(string)
	actorID, _ := uuid.Parse(actorIDStr)

	bNew, _ := json.Marshal(flag)
	var newState models.JSONB
	json.Unmarshal(bNew, &newState)

	h.audit.LogAction(r.Context(), &models.AuditLog{
		ID:         uuid.New(),
		ProjectID:     &projectID,
		ActorID:    actorID,
		Action:     "CREATE",
		TargetType: "FLAG",
		TargetID:      flag.ID,
		NewState:   newState,
		ActorIP:    r.RemoteAddr,
		CreatedAt:  time.Now().UTC(),
	})

	var parentStr *string
	if flag.ParentFlagID != nil {
		s := flag.ParentFlagID.String()
		parentStr = &s
	}

	RespondWithJSON(w, http.StatusCreated, dto.FeatureFlagResponse{
		ID:           flag.ID.String(),
		ProjectID:    flag.ProjectID.String(),
		Key:          flag.Key,
		Name:         flag.Name,
		Description:  flag.Description,
		Type:         string(flag.Type),
		ParentFlagID: parentStr,
		CreatedAt:    flag.CreatedAt,
		UpdatedAt:    flag.UpdatedAt,
	})
}

func (h *FlagHandler) ListFlags(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	pagination := GetPagination(r)
	offset := TokenToOffset(pagination.PageToken)

	flags, total, err := h.store.FlagRepo().ListByProject(r.Context(), projectID, pagination.PageSize, offset)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to list flags")
		return
	}

	var dtos []dto.FeatureFlagResponse
	for _, f := range flags {
		var parentStr *string
		if f.ParentFlagID != nil {
			s := f.ParentFlagID.String()
			parentStr = &s
		}
		dtos = append(dtos, dto.FeatureFlagResponse{
			ID:           f.ID.String(),
			ProjectID:    f.ProjectID.String(),
			Key:          f.Key,
			Name:         f.Name,
			Description:  f.Description,
			Type:         string(f.Type),
			ParentFlagID: parentStr,
			CreatedAt:    f.CreatedAt,
			UpdatedAt:    f.UpdatedAt,
		})
	}
	if dtos == nil {
		dtos = []dto.FeatureFlagResponse{}
	}

	resp := dto.PaginatedResponse{Data: dtos}
	if offset+len(flags) < total {
		resp.NextPageToken = OffsetToToken(offset + len(flags))
	}

	RespondWithJSON(w, http.StatusOK, resp)
}

func (h *FlagHandler) GetFlagState(w http.ResponseWriter, r *http.Request) {
	envIDStr := chi.URLParam(r, "envId")
	flagIDStr := chi.URLParam(r, "flagId")

	envID, err1 := uuid.Parse(envIDStr)
	flagID, err2 := uuid.Parse(flagIDStr)
	if err1 != nil || err2 != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	state, err := h.store.FlagStateRepo().GetByEnvAndFlag(r.Context(), envID, flagID)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Flag state not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve flag state")
		return
	}

	RespondWithJSON(w, http.StatusOK, dto.FlagStateResponse{
		EnvironmentID:  state.EnvironmentID.String(),
		FeatureFlagID:  state.FeatureFlagID.String(),
		Enabled:        state.Enabled,
		TargetingRules: state.TargetingRules,
		RemoteConfig:   state.RemoteConfig,
		UpdatedAt:      state.UpdatedAt,
	})
}

func (h *FlagHandler) UpdateFlagState(w http.ResponseWriter, r *http.Request) {
	envIDStr := chi.URLParam(r, "envId")
	flagIDStr := chi.URLParam(r, "flagId")

	envID, err1 := uuid.Parse(envIDStr)
	flagID, err2 := uuid.Parse(flagIDStr)
	if err1 != nil || err2 != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req dto.UpdateFlagStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	state, err := h.store.FlagStateRepo().GetByEnvAndFlag(r.Context(), envID, flagID)
	if err != nil {
		if err == repository.ErrNotFound {
			now := time.Now().UTC()
			state = &models.EnvironmentFlagState{
				ID:             uuid.New(),
				EnvironmentID:  envID,
				FeatureFlagID:  flagID,
				Enabled:        req.Enabled,
				TargetingRules: req.TargetingRules,
				RemoteConfig:   req.RemoteConfig,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			if err := h.store.FlagStateRepo().Create(r.Context(), state); err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Failed to create flag state")
				return
			}
			
			actorIDStr := r.Context().Value(UserIDKey).(string)
			actorID, _ := uuid.Parse(actorIDStr)

			bStateNew, _ := json.Marshal(state)
			var newState models.JSONB
			json.Unmarshal(bStateNew, &newState)

			h.audit.LogAction(r.Context(), &models.AuditLog{
				ID:            uuid.New(),
				EnvironmentID: &envID,
				ActorID:       actorID,
				Action:        "CREATE",
				TargetType:    "FLAG_STATE",
				TargetID:      state.ID,
				NewState:      newState,
				ActorIP:       r.RemoteAddr,
				CreatedAt:     time.Now().UTC(),
			})
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve flag state")
			return
		}
	} else {
		bPrev, _ := json.Marshal(state)
		var prevState models.JSONB
		json.Unmarshal(bPrev, &prevState)
		
		state.Enabled = req.Enabled
		state.TargetingRules = req.TargetingRules
		state.RemoteConfig = req.RemoteConfig
		state.UpdatedAt = time.Now().UTC()

		if err := h.store.FlagStateRepo().Update(r.Context(), state); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to update flag state")
			return
		}
		
		actorIDStr := r.Context().Value(UserIDKey).(string)
		actorID, _ := uuid.Parse(actorIDStr)

		bStateNew, _ := json.Marshal(state)
		var newState models.JSONB
		json.Unmarshal(bStateNew, &newState)

		h.audit.LogAction(r.Context(), &models.AuditLog{
			ID:            uuid.New(),
			EnvironmentID: &envID,
			ActorID:       actorID,
			Action:        "UPDATE",
			TargetType:    "FLAG_STATE",
			TargetID:      state.ID,
			PreviousState: prevState,
			NewState:      newState,
			ActorIP:       r.RemoteAddr,
			CreatedAt:     time.Now().UTC(),
		})
	}

	if h.cacheClient != nil {
		_ = h.cacheClient.PublishRulesetUpdate(r.Context(), envID.String(), time.Now().UTC().String())
	}

	RespondWithJSON(w, http.StatusOK, dto.FlagStateResponse{
		EnvironmentID:  state.EnvironmentID.String(),
		FeatureFlagID:  state.FeatureFlagID.String(),
		Enabled:        state.Enabled,
		TargetingRules: state.TargetingRules,
		RemoteConfig:   state.RemoteConfig,
		UpdatedAt:      state.UpdatedAt,
	})
}
