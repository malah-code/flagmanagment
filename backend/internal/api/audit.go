package api

import (
	"net/http"

	"github.com/flagmanagment/backend/internal/dto"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuditHandler struct {
	store repository.Store
}

func NewAuditHandler(store repository.Store) *AuditHandler {
	return &AuditHandler{store: store}
}

func (h *AuditHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	pagination := GetPagination(r)
	offset := TokenToOffset(pagination.PageToken)

	logs, total, err := h.store.AuditRepo().ListByProject(r.Context(), projectID, pagination.PageSize, offset)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to list audit logs")
		return
	}

	// Just reuse PaginatedResponse logic with anon struct or generic map. We can use map since no DTO yet.
	resp := dto.PaginatedResponse{Data: logs}
	if offset+len(logs) < total {
		resp.NextPageToken = OffsetToToken(offset + len(logs))
	}

	RespondWithJSON(w, http.StatusOK, resp)
}
