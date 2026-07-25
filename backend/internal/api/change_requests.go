package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/flagmanagment/backend/internal/cache"
	"github.com/flagmanagment/backend/internal/dto"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ChangeRequestHandler struct {
	store       repository.Store
	crService   *services.ChangeRequestService
	rbac        *RBACMiddleware
	cacheClient *cache.Client
}

func NewChangeRequestHandler(store repository.Store, crService *services.ChangeRequestService, rbac *RBACMiddleware, cacheClient *cache.Client) *ChangeRequestHandler {
	return &ChangeRequestHandler{
		store:       store,
		crService:   crService,
		rbac:        rbac,
		cacheClient: cacheClient,
	}
}

func (h *ChangeRequestHandler) RegisterRoutes(r chi.Router) {
	r.With(h.rbac.RequireRole("VIEWER")).Get("/environments/{envId}/change-requests", h.ListByEnvironment)
	r.With(h.rbac.RequireRole("VIEWER")).Get("/change-requests/{id}", h.GetByID)
	r.With(h.rbac.RequireRole("RELEASE_MANAGER")).Post("/change-requests/{id}/approve", h.Approve)
	r.With(h.rbac.RequireRole("RELEASE_MANAGER")).Post("/change-requests/{id}/reject", h.Reject)
}

func (h *ChangeRequestHandler) ListByEnvironment(w http.ResponseWriter, r *http.Request) {
	envIDStr := chi.URLParam(r, "envId")
	envID, err := uuid.Parse(envIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}

	status := r.URL.Query().Get("status")
	pagination := GetPagination(r)
	offset := TokenToOffset(pagination.PageToken)

	requests, total, err := h.crService.ListByEnvironment(r.Context(), envID, status, pagination.PageSize, offset)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to list change requests")
		return
	}

	resp := dto.PaginatedResponse{Data: requests}
	if offset+len(requests) < total {
		resp.NextPageToken = OffsetToToken(offset + len(requests))
	}

	RespondWithJSON(w, http.StatusOK, resp)
}

func (h *ChangeRequestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid change request ID")
		return
	}

	cr, err := h.crService.GetByID(r.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Change request not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve change request")
		return
	}

	RespondWithJSON(w, http.StatusOK, cr)
}

func (h *ChangeRequestHandler) Approve(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid change request ID")
		return
	}

	actorIDStr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	actorID, _ := uuid.Parse(actorIDStr)

	var req dto.ApproveChangeRequestRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	err = h.crService.Approve(r.Context(), id, actorID, req.Comment)
	if err != nil {
		if err == services.ErrSelfApprovalNotAllowed {
			RespondWithError(w, http.StatusForbidden, err.Error())
			return
		}
		if err == services.ErrChangeRequestNotPending {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Change request not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to approve change request")
		return
	}

	cr, _ := h.crService.GetByID(r.Context(), id)
	
	if h.cacheClient != nil {
		_ = h.cacheClient.PublishRulesetUpdate(r.Context(), cr.EnvironmentID.String(), time.Now().UTC().String())
	}

	RespondWithJSON(w, http.StatusOK, cr)
}

func (h *ChangeRequestHandler) Reject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid change request ID")
		return
	}

	actorIDStr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	actorID, _ := uuid.Parse(actorIDStr)

	var req dto.RejectChangeRequestRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	err = h.crService.Reject(r.Context(), id, actorID, req.Reason)
	if err != nil {
		if err == services.ErrChangeRequestNotPending {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Change request not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to reject change request")
		return
	}

	cr, _ := h.crService.GetByID(r.Context(), id)
	RespondWithJSON(w, http.StatusOK, cr)
}
