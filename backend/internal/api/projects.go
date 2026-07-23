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
)

type ProjectHandler struct {
	store    repository.Store
	validate *validator.Validate
	rbac     *RBACMiddleware
	audit    *AuditHandler
}

func NewProjectHandler(store repository.Store, rbac *RBACMiddleware, audit *AuditHandler) *ProjectHandler {
	return &ProjectHandler{
		store:    store,
		validate: validator.New(),
		rbac:     rbac,
		audit:    audit,
	}
}

func (h *ProjectHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.Get)
		r.With(h.rbac.RequireRole("EDITOR")).Put("/", h.Update)
		r.With(h.rbac.RequireRole("ADMIN")).Delete("/", h.Delete)
		r.With(h.rbac.RequireRole("VIEWER")).Get("/audit-logs", h.audit.ListByProject)
	})
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	key := req.Key
	if key == "" {
		key = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	}

	now := time.Now().UTC()
	project := &models.Project{
		ID:          uuid.New(),
		Name:        req.Name,
		Key:         key,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.store.ProjectRepo().Create(r.Context(), project); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	RespondWithJSON(w, http.StatusCreated, dto.ProjectResponse{
		ID:          project.ID.String(),
		Name:        project.Name,
		Key:         project.Key,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	})
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	pagination := GetPagination(r)
	offset := TokenToOffset(pagination.PageToken)

	projects, total, err := h.store.ProjectRepo().List(r.Context(), pagination.PageSize, offset)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to list projects")
		return
	}

	var dtos []dto.ProjectResponse
	for _, p := range projects {
		dtos = append(dtos, dto.ProjectResponse{
			ID:          p.ID.String(),
			Name:        p.Name,
			Key:         p.Key,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	if dtos == nil {
		dtos = []dto.ProjectResponse{}
	}

	resp := dto.PaginatedResponse{Data: dtos}
	if offset+len(projects) < total {
		resp.NextPageToken = OffsetToToken(offset + len(projects))
	}

	RespondWithJSON(w, http.StatusOK, resp)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.store.ProjectRepo().GetByID(r.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve project")
		return
	}

	RespondWithJSON(w, http.StatusOK, dto.ProjectResponse{
		ID:          project.ID.String(),
		Name:        project.Name,
		Key:         project.Key,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	})
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var req dto.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	project, err := h.store.ProjectRepo().GetByID(r.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve project")
		return
	}

	project.Name = req.Name
	if req.Description != "" {
		project.Description = req.Description
	}
	project.UpdatedAt = time.Now().UTC()
	if err := h.store.ProjectRepo().Update(r.Context(), project); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to update project")
		return
	}

	RespondWithJSON(w, http.StatusOK, dto.ProjectResponse{
		ID:          project.ID.String(),
		Name:        project.Name,
		Key:         project.Key,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	})
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	if err := h.store.ProjectRepo().Delete(r.Context(), id); err != nil {
		if err == repository.ErrNotFound {
			RespondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
