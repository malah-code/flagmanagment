package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
)

type ServiceAccountHandler struct {
	store                 repository.Store
	serviceAccountService services.ServiceAccountService
}

func NewServiceAccountHandler(store repository.Store) *ServiceAccountHandler {
	return &ServiceAccountHandler{
		store:                 store,
		serviceAccountService: services.NewServiceAccountService(store),
	}
}

type CreateServiceAccountRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (h *ServiceAccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateServiceAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Assuming a user is making this request (we should get this from context, mock for now)
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		// Fallback to nil if no user ID provided in header (for MVP simplicity if auth isn't fully wired)
		userID = uuid.Nil
	}

	sa := &models.ServiceAccount{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
	}

	if err := h.store.ServiceAccountRepo().Create(r.Context(), sa); err != nil {
		http.Error(w, "Failed to create service account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sa)
}

type CreateServiceAccountKeyRequest struct {
	Name          string `json:"name"`
	ExpiresInDays *int   `json:"expires_in_days"`
}

type CreateServiceAccountKeyResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Token string    `json:"token"`
}

func (h *ServiceAccountHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	saIDStr := chi.URLParam(r, "id")
	saID, err := uuid.Parse(saIDStr)
	if err != nil {
		http.Error(w, "Invalid service account ID", http.StatusBadRequest)
		return
	}

	var req CreateServiceAccountKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var expiresAt *time.Time
	if req.ExpiresInDays != nil {
		t := time.Now().Add(time.Duration(*req.ExpiresInDays) * 24 * time.Hour)
		expiresAt = &t
	}

	key, plaintextKey, err := h.serviceAccountService.CreateKey(r.Context(), saID, req.Name, expiresAt)
	if err != nil {
		http.Error(w, "Failed to create key", http.StatusInternalServerError)
		return
	}

	resp := CreateServiceAccountKeyResponse{
		ID:    key.ID,
		Name:  key.Name,
		Token: plaintextKey,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *ServiceAccountHandler) ListKeys(w http.ResponseWriter, r *http.Request) {
	saIDStr := chi.URLParam(r, "id")
	saID, err := uuid.Parse(saIDStr)
	if err != nil {
		http.Error(w, "Invalid service account ID", http.StatusBadRequest)
		return
	}

	keys, err := h.store.ServiceAccountRepo().ListKeys(r.Context(), saID)
	if err != nil {
		http.Error(w, "Failed to list keys", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

func (h *ServiceAccountHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.Create)
	r.Post("/{id}/keys", h.CreateKey)
	r.Get("/{id}/keys", h.ListKeys)
}
