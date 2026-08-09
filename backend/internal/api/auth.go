package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"github.com/go-chi/chi/v5"
	"github.com/flagmanagment/backend/internal/auth"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
)

type AuthHandler struct {
	store        repository.Store
	authService  services.AuthService
	oidcProvider *oidc.Provider
	oauth2Config oauth2.Config
}

func NewAuthHandler(store repository.Store) *AuthHandler {
	var provider *oidc.Provider
	var oauth2Conf oauth2.Config

	clientID := os.Getenv("OIDC_CLIENT_ID")
	issuerURL := os.Getenv("OIDC_ISSUER_URL")
	if clientID != "" && issuerURL != "" {
		p, err := oidc.NewProvider(context.Background(), issuerURL)
		if err == nil {
			provider = p
			apiURL := os.Getenv("API_URL")
			if apiURL == "" {
				apiURL = "http://localhost:8080"
			}
			oauth2Conf = oauth2.Config{
				ClientID:     clientID,
				ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
				RedirectURL:  apiURL + "/api/v1/auth/sso/callback/oidc",
				Endpoint:     provider.Endpoint(),
				Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
			}
		}
	}

	return &AuthHandler{
		store:        store,
		authService:  services.NewAuthService(store),
		oidcProvider: provider,
		oauth2Config: oauth2Conf,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.store.UserRepo().GetByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if user.PasswordHash == nil || !auth.CheckPasswordHash(req.Password, *user.PasswordHash) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var resp LoginResponse
	resp.Token = token
	resp.User.ID = user.ID.String()
	resp.User.Email = user.Email

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) SSOLogin(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	if provider == "oidc" {
		if h.oidcProvider == nil {
			http.Error(w, "OIDC not configured", http.StatusInternalServerError)
			return
		}
		// Generate random state
		b := make([]byte, 16)
		rand.Read(b)
		state := base64.URLEncoding.EncodeToString(b)
		
		// Set state cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "sso_state",
			Value:    state,
			MaxAge:   int(time.Hour.Seconds()),
			Secure:   r.TLS != nil,
			HttpOnly: true,
			Path:     "/",
		})

		url := h.oauth2Config.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	if provider == "saml" {
		http.Redirect(w, r, "/api/v1/auth/saml/login", http.StatusFound)
		return
	}
	http.Error(w, "Unsupported provider", http.StatusBadRequest)
}

func (h *AuthHandler) SAMLLogin(w http.ResponseWriter, r *http.Request) {
	// In a complete implementation, this would use samlsp.Middleware.
	// For this MVP, we return a 501 Not Implemented to satisfy the route structure.
	http.Error(w, "SAML Login Not Implemented yet. Please configure OIDC.", http.StatusNotImplemented)
}

func (h *AuthHandler) SAMLCallback(w http.ResponseWriter, r *http.Request) {
	// SAML ACS Callback
	http.Error(w, "SAML Callback Not Implemented", http.StatusNotImplemented)
}

func (h *AuthHandler) SSOCallbackOIDC(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("sso_state")
	if err != nil || r.URL.Query().Get("state") != stateCookie.Value {
		http.Error(w, "State did not match", http.StatusBadRequest)
		return
	}

	oauth2Token, err := h.oauth2Config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token", http.StatusInternalServerError)
		return
	}

	verifier := h.oidcProvider.Verifier(&oidc.Config{ClientID: h.oauth2Config.ClientID})
	idToken, err := verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID token", http.StatusInternalServerError)
		return
	}

	var claims struct {
		Email string `json:"email"`
		Sub   string `json:"sub"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
		return
	}

	user, err := h.authService.HandleSSOLogin(r.Context(), "oidc", claims.Email, claims.Sub)
	if err != nil {
		http.Error(w, "Failed to process SSO login: "+err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Redirect to frontend with token (stateless approach for MVP)
	// In production, setting an HttpOnly cookie is better, but this satisfies the basic UI callback.
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	http.Redirect(w, r, frontendURL+"/sso-success?token="+token, http.StatusFound)
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/login", h.Login)
	r.Get("/sso/login", h.SSOLogin)
	r.Get("/sso/callback/oidc", h.SSOCallbackOIDC)
	r.Get("/saml/login", h.SAMLLogin)
	r.Post("/saml/acs", h.SAMLCallback)
}
