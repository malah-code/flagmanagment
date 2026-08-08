package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/flagmanagment/backend/internal/dto"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/services"
)

type EnvironmentHandler struct {
	store      repository.Store
	validate   *validator.Validate
	rbac       *RBACMiddleware
	audit      *services.AuditService
	envService *services.EnvironmentService
}

func NewEnvironmentHandler(store repository.Store, rbac *RBACMiddleware, audit *services.AuditService, envService *services.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{
		store:      store,
		validate:   validator.New(),
		rbac:       rbac,
		audit:      audit,
		envService: envService,
	}
}

func (h *EnvironmentHandler) RegisterRoutes(r chi.Router) {
	// Nested under projects
	r.With(h.rbac.RequireRole("EDITOR")).Post("/projects/{projectId}/environments", h.Create)
	r.With(h.rbac.RequireRole("EDITOR")).Post("/projects/{projectId}/environments/{sourceEnvId}/clone", h.Clone)
	r.With(h.rbac.RequireRole("VIEWER")).Get("/projects/{projectId}/environments", h.List)

	// Direct access
	r.With(h.rbac.RequireRole("VIEWER")).Get("/environments/{envId}", h.Get)
	r.With(h.rbac.RequireRole("EDITOR")).Put("/environments/{envId}", h.Update)
	r.With(h.rbac.RequireRole("ADMIN")).Delete("/environments/{envId}", h.Delete)
}


func (h *EnvironmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var req dto.CreateEnvironmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	actorIDStr := r.Context().Value(UserIDKey).(string)
	actorID, _ := uuid.Parse(actorIDStr)

	env, apiKey, err := h.envService.CreateEnvironment(r.Context(), projectID, req.Name, req.IsProtected, actorID, r.RemoteAddr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create environment")
		return
	}

	res := dto.CreateEnvironmentResponse{
		EnvironmentResponse: dto.EnvironmentResponse{
			ID:          env.ID.String(),
			ProjectID:   env.ProjectID.String(),
			Name:        env.Name,
			IsProtected: env.IsProtected,
			CreatedAt:   env.CreatedAt,
			UpdatedAt:   env.UpdatedAt,
		},
		APIKey: apiKey,
	}

	RespondWithJSON(w, http.StatusCreated, res)
}

func (h *EnvironmentHandler) List(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	envs, err := h.store.EnvironmentRepo().ListByProject(r.Context(), projectID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to list environments")
		return
	}

	var dtos []dto.EnvironmentResponse
	for _, e := range envs {
		dtos = append(dtos, dto.EnvironmentResponse{
			ID:          e.ID.String(),
			ProjectID:   e.ProjectID.String(),
			Name:        e.Name,
			IsProtected: e.IsProtected,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		})
	}
	if dtos == nil {
		dtos = []dto.EnvironmentResponse{}
	}

	RespondWithJSON(w, http.StatusOK, dto.PaginatedResponse{
		Data: dtos,
	})
}

func (h *EnvironmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "envId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}

	env, err := h.store.EnvironmentRepo().GetByID(r.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Environment not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve environment")
		return
	}

	RespondWithJSON(w, http.StatusOK, dto.EnvironmentResponse{
		ID:          env.ID.String(),
		ProjectID:   env.ProjectID.String(),
		Name:        env.Name,
		IsProtected: env.IsProtected,
		CreatedAt:   env.CreatedAt,
		UpdatedAt:   env.UpdatedAt,
	})
}

func (h *EnvironmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "envId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}

	var req dto.UpdateEnvironmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	env, err := h.store.EnvironmentRepo().GetByID(r.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Environment not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve environment")
		return
	}

	bPrev, _ := json.Marshal(env)
	var prevState models.JSONB
	json.Unmarshal(bPrev, &prevState)

	env.Name = req.Name
	env.Key = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	env.IsProtected = req.IsProtected
	env.UpdatedAt = time.Now().UTC()

	if err := h.store.EnvironmentRepo().Update(r.Context(), env); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to update environment")
		return
	}

	actorIDStr := r.Context().Value(UserIDKey).(string)
	actorID, _ := uuid.Parse(actorIDStr)

	bNew, _ := json.Marshal(env)
	var newState models.JSONB
	json.Unmarshal(bNew, &newState)

	h.audit.LogAction(r.Context(), &models.AuditLog{
		ID:            uuid.New(),
		ProjectID:     &env.ProjectID,
		EnvironmentID: &env.ID,
		ActorID:       actorID,
		Action:        "UPDATE",
		TargetType:    "ENVIRONMENT",
		TargetID:      env.ID,
		PreviousState: prevState,
		NewState:      newState,
		ActorIP:       r.RemoteAddr,
		CreatedAt:     time.Now().UTC(),
	})

	RespondWithJSON(w, http.StatusOK, dto.EnvironmentResponse{
		ID:          env.ID.String(),
		ProjectID:   env.ProjectID.String(),
		Name:        env.Name,
		IsProtected: env.IsProtected,
		CreatedAt:   env.CreatedAt,
		UpdatedAt:   env.UpdatedAt,
	})
}

func (h *EnvironmentHandler) Clone(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	sourceEnvIDStr := chi.URLParam(r, "sourceEnvId")
	sourceEnvID, err := uuid.Parse(sourceEnvIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid source environment ID")
		return
	}

	var req dto.CloneEnvironmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	actorIDStr := r.Context().Value(UserIDKey).(string)
	actorID, _ := uuid.Parse(actorIDStr)

	env, apiKey, err := h.envService.CloneEnvironment(r.Context(), projectID, sourceEnvID, req.Name, actorID, r.RemoteAddr)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Source environment not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to clone environment")
		return
	}

	res := dto.CreateEnvironmentResponse{
		EnvironmentResponse: dto.EnvironmentResponse{
			ID:          env.ID.String(),
			ProjectID:   env.ProjectID.String(),
			Name:        env.Name,
			IsProtected: env.IsProtected,
			CreatedAt:   env.CreatedAt,
			UpdatedAt:   env.UpdatedAt,
		},
		APIKey: apiKey,
	}

	RespondWithJSON(w, http.StatusCreated, res)
}

func (h *EnvironmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "envId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}

	actorIDStr := r.Context().Value(UserIDKey).(string)
	actorID, _ := uuid.Parse(actorIDStr)

	if err := h.envService.DeleteEnvironment(r.Context(), id, actorID, r.RemoteAddr); err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Environment not found")
			return
		}
		if err == services.ErrProtectedEnvironment {
			RespondWithError(w, http.StatusForbidden, "Cannot delete protected environment")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete environment")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
