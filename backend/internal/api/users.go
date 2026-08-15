package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/services"
)

type UsersHandler struct {
	userService services.UserService
}

func NewUsersHandler(userService services.UserService) *UsersHandler {
	return &UsersHandler{
		userService: userService,
	}
}

type GetUsersResponse struct {
	Users []interface{} `json:"users"`
	Total int           `json:"total"`
}

func (h *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	users, total, err := h.userService.GetUsers(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// We wrap it in interface to match the JSON structure expected by the frontend
	var usersInterface []interface{}
	for _, u := range users {
		usersInterface = append(usersInterface, u)
	}

	if usersInterface == nil {
		usersInterface = []interface{}{}
	}

	resp := GetUsersResponse{
		Users: usersInterface,
		Total: total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type InviteUserRequest struct {
	Email      string   `json:"email"`
	Role       string   `json:"role"`
	ProjectIDs []string `json:"project_ids"`
}

func (h *UsersHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
	var req InviteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Role == "" {
		http.Error(w, "Email and role are required", http.StatusBadRequest)
		return
	}

	userIDVal := r.Context().Value(UserIDKey)
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	currentUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	inv, err := h.userService.CreateInvitation(r.Context(), req.Email, req.Role, req.ProjectIDs, currentUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inv)
}

type UpdateAccessRequest struct {
	Role       string   `json:"role"`
	ProjectIDs []string `json:"project_ids"`
}

func (h *UsersHandler) UpdateUserAccess(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req UpdateAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Role == "" {
		http.Error(w, "Role is required", http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdateUserAccess(r.Context(), userID, req.Role, req.ProjectIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *UsersHandler) RegisterRoutes(r chi.Router) {
	// The caller will prefix with /api/v1/users
	r.Get("/", h.ListUsers)
	r.Post("/invite", h.InviteUser)
	r.Put("/{id}/access", h.UpdateUserAccess)
}
