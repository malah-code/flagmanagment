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

type ScheduledChangeHandler struct {
	store       repository.Store
	scService   *services.ScheduledChangeService
	rbac        *RBACMiddleware
	cacheClient *cache.Client
	validate    *validator.Validate
}

func NewScheduledChangeHandler(store repository.Store, scService *services.ScheduledChangeService, rbac *RBACMiddleware, cacheClient *cache.Client) *ScheduledChangeHandler {
	return &ScheduledChangeHandler{
		store:       store,
		scService:   scService,
		rbac:        rbac,
		cacheClient: cacheClient,
		validate:    validator.New(),
	}
}

func (h *ScheduledChangeHandler) RegisterRoutes(r chi.Router) {
	r.With(h.rbac.RequireRole("RELEASE_MANAGER")).Post("/environments/{envId}/scheduled-changes", h.Create)
	r.With(h.rbac.RequireRole("VIEWER")).Get("/environments/{envId}/scheduled-changes", h.List)
	r.With(h.rbac.RequireRole("VIEWER")).Get("/scheduled-changes/{id}", h.GetByID)
	r.With(h.rbac.RequireRole("RELEASE_MANAGER")).Patch("/scheduled-changes/{id}", h.Update)
	r.With(h.rbac.RequireRole("RELEASE_MANAGER")).Delete("/scheduled-changes/{id}", h.Cancel)
}

// parseActorID extracts the authenticated user UUID from the request context.
// Returns uuid.Nil and false if the actor ID is missing or malformed.
func parseActorID(r *http.Request) (uuid.UUID, bool) {
	raw, ok := r.Context().Value(UserIDKey).(string)
	if !ok || raw == "" {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

// toScheduledChangeResponse maps a model to the public DTO, keeping internal
// fields (db tags, unexported fields) out of the HTTP response surface.
func toScheduledChangeResponse(sc *models.ScheduledChange) dto.ScheduledChangeResponse {
	return dto.ScheduledChangeResponse{
		ID:            sc.ID.String(),
		ProjectID:     sc.ProjectID.String(),
		EnvironmentID: sc.EnvironmentID.String(),
		TargetType:    string(sc.TargetType),
		TargetID:      sc.TargetID.String(),
		Action:        string(sc.Action),
		ScheduledFor:  sc.ScheduledFor,
		Status:        string(sc.Status),
		CreatedBy:     sc.CreatedBy.String(),
		ExecutedAt:    sc.ExecutedAt,
		CancelledAt:   sc.CancelledAt,
		CreatedAt:     sc.CreatedAt,
		UpdatedAt:     sc.UpdatedAt,
	}
}

func (h *ScheduledChangeHandler) Create(w http.ResponseWriter, r *http.Request) {
	envIDStr := chi.URLParam(r, "envId")
	envID, err := uuid.Parse(envIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}

	actorID, ok := parseActorID(r)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.CreateScheduledChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	scheduledFor, err := time.Parse(time.RFC3339, req.ScheduledFor)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid scheduled_for timestamp format (must be RFC3339/ISO-8601 UTC)")
		return
	}

	targetID, err := uuid.Parse(req.TargetID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid target_id")
		return
	}

	env, err := h.store.EnvironmentRepo().GetByID(r.Context(), envID)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Environment not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve environment")
		return
	}

	// Validate target_id entity existence
	if req.TargetType == "FLAG" {
		if _, err := h.store.FlagRepo().GetByID(r.Context(), targetID); err != nil {
			if err == repository.ErrNotFound {
				RespondWithError(w, http.StatusNotFound, "Target feature flag not found")
				return
			}
			RespondWithError(w, http.StatusInternalServerError, "Failed to verify feature flag existence")
			return
		}
	} else if req.TargetType == "CHANGE_REQUEST" {
		if _, err := h.store.ChangeRequestRepo().GetByID(r.Context(), targetID); err != nil {
			if err == repository.ErrNotFound {
				RespondWithError(w, http.StatusNotFound, "Target change request not found")
				return
			}
			RespondWithError(w, http.StatusInternalServerError, "Failed to verify change request existence")
			return
		}
	}

	sc := &models.ScheduledChange{
		ProjectID:     env.ProjectID,
		EnvironmentID: envID,
		TargetType:    models.ScheduledChangeTargetType(req.TargetType),
		TargetID:      targetID,
		Action:        models.ScheduledChangeAction(req.Action),
		ScheduledFor:  scheduledFor.UTC(),
		CreatedBy:     actorID,
	}

	if err := h.scService.Create(r.Context(), sc); err != nil {
		if err == repository.ErrPendingScheduleExists {
			RespondWithError(w, http.StatusConflict, err.Error())
			return
		}
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, toScheduledChangeResponse(sc))
}

func (h *ScheduledChangeHandler) List(w http.ResponseWriter, r *http.Request) {
	envIDStr := chi.URLParam(r, "envId")
	envID, err := uuid.Parse(envIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}

	status := r.URL.Query().Get("status")
	pagination := GetPagination(r)
	offset := TokenToOffset(pagination.PageToken)

	schedules, total, err := h.scService.ListByEnvironment(r.Context(), envID, status, pagination.PageSize, offset)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to list scheduled changes")
		return
	}

	responses := make([]dto.ScheduledChangeResponse, 0, len(schedules))
	for _, sc := range schedules {
		responses = append(responses, toScheduledChangeResponse(sc))
	}

	resp := dto.PaginatedResponse{Data: responses}
	if offset+len(schedules) < total {
		resp.NextPageToken = OffsetToToken(offset + len(schedules))
	}

	RespondWithJSON(w, http.StatusOK, resp)
}

func (h *ScheduledChangeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid scheduled change ID")
		return
	}

	sc, err := h.scService.GetByID(r.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Scheduled change not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve scheduled change")
		return
	}

	RespondWithJSON(w, http.StatusOK, toScheduledChangeResponse(sc))
}

func (h *ScheduledChangeHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid scheduled change ID")
		return
	}

	actorID, ok := parseActorID(r)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.UpdateScheduledChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	newTime, err := time.Parse(time.RFC3339, req.ScheduledFor)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid scheduled_for timestamp format (must be RFC3339/ISO-8601 UTC)")
		return
	}

	// Ownership check: verify the scheduled change belongs to a project the actor can access.
	existing, err := h.scService.GetByID(r.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Scheduled change not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve scheduled change")
		return
	}
	if err := h.requireProjectAccess(r, existing.ProjectID); err != nil {
		RespondWithError(w, http.StatusForbidden, "Access denied to this scheduled change")
		return
	}

	updated, err := h.scService.UpdateScheduledFor(r.Context(), id, newTime.UTC(), actorID)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Scheduled change not found")
			return
		}
		if err == services.ErrScheduleNotPending {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, toScheduledChangeResponse(updated))
}

func (h *ScheduledChangeHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid scheduled change ID")
		return
	}

	actorID, ok := parseActorID(r)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Ownership check: verify the scheduled change belongs to a project the actor can access.
	existing, err := h.scService.GetByID(r.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Scheduled change not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve scheduled change")
		return
	}
	if err := h.requireProjectAccess(r, existing.ProjectID); err != nil {
		RespondWithError(w, http.StatusForbidden, "Access denied to this scheduled change")
		return
	}

	cancelled, err := h.scService.Cancel(r.Context(), id, actorID)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Scheduled change not found")
			return
		}
		if err == services.ErrScheduleNotPending {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to cancel scheduled change")
		return
	}

	RespondWithJSON(w, http.StatusOK, toScheduledChangeResponse(cancelled))
}

// requireProjectAccess checks that the authenticated user has at least one role
// in the given project. Returns an error if not, which the caller maps to 403.
func (h *ScheduledChangeHandler) requireProjectAccess(r *http.Request, projectID uuid.UUID) error {
	actorID, ok := parseActorID(r)
	if !ok {
		return repository.ErrNotFound
	}
	roles, err := h.store.RoleRepo().GetUserRoles(r.Context(), actorID, &projectID)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return repository.ErrNotFound
	}
	return nil
}
