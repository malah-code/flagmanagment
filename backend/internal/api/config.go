package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
)

type ConfigHandler struct {
	store         repository.Store
	cryptoService services.CryptoService
	emailService  services.EmailService
}

func NewConfigHandler(store repository.Store, crypto services.CryptoService, email services.EmailService) *ConfigHandler {
	return &ConfigHandler{
		store:         store,
		cryptoService: crypto,
		emailService:  email,
	}
}

type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *ConfigHandler) GetSMTP(w http.ResponseWriter, r *http.Request) {
	config, err := h.store.SystemConfigRepo().GetByKey(r.Context(), "smtp_config")
	if err != nil && err != repository.ErrNotFound {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var smtpConfig SMTPConfig
	if config != nil && config.Value != nil {
		if encVal, ok := config.Value["encrypted_data"].(string); ok {
			decrypted, err := h.cryptoService.DecryptAES(encVal)
			if err == nil {
				json.Unmarshal([]byte(decrypted), &smtpConfig)
				smtpConfig.Password = "" // Hide password in GET
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(smtpConfig)
}

func (h *ConfigHandler) UpdateSMTP(w http.ResponseWriter, r *http.Request) {
	var req SMTPConfig
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		existing, err := h.store.SystemConfigRepo().GetByKey(r.Context(), "smtp_config")
		if err == nil && existing != nil && existing.Value != nil {
			if encVal, ok := existing.Value["encrypted_data"].(string); ok {
				var existingConfig SMTPConfig
				decrypted, _ := h.cryptoService.DecryptAES(encVal)
				json.Unmarshal([]byte(decrypted), &existingConfig)
				req.Password = existingConfig.Password
			}
		}
	}

	configBytes, _ := json.Marshal(req)
	encrypted, err := h.cryptoService.EncryptAES(string(configBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config := &models.SystemConfig{
		Key: "smtp_config",
		Value: models.JSONB{
			"encrypted_data": encrypted,
		},
	}

	if err := h.store.SystemConfigRepo().Upsert(r.Context(), config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func (h *ConfigHandler) TestSMTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if err := h.emailService.SendInvitation(req.Email, "https://example.com/test-link"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ConfigHandler) RegisterRoutes(r chi.Router) {
	r.Get("/smtp", h.GetSMTP)
	r.Put("/smtp", h.UpdateSMTP)
	r.Post("/smtp/test", h.TestSMTP)
}
