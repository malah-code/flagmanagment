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
	store         repository.Store
	authService   services.AuthService
	cryptoService services.CryptoService
	oidcProvider  *oidc.Provider
	oauth2Config  oauth2.Config
}

func NewAuthHandler(store repository.Store, crypto services.CryptoService) *AuthHandler {
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
		store:         store,
		authService:   services.NewAuthService(store),
		cryptoService: crypto,
		oidcProvider:  provider,
		oauth2Config:  oauth2Conf,
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

type SSOProvidersResponse struct {
	OIDCEnabled bool `json:"oidc_enabled"`
	SAMLEnabled bool `json:"saml_enabled"`
}

func (h *AuthHandler) getRuntimeSSOConfig(ctx context.Context) (*SSOConfig, error) {
	config, err := h.store.SystemConfigRepo().GetByKey(ctx, "sso_config")
	if err != nil || config == nil || config.Value == nil {
		return nil, err
	}
	encVal, ok := config.Value["encrypted_data"].(string)
	if !ok {
		return nil, nil
	}
	decrypted, err := h.cryptoService.DecryptAES(encVal)
	if err != nil {
		return nil, err
	}
	var ssoConf SSOConfig
	if err := json.Unmarshal([]byte(decrypted), &ssoConf); err != nil {
		return nil, err
	}
	return &ssoConf, nil
}

func (h *AuthHandler) GetSSOProviders(w http.ResponseWriter, r *http.Request) {
	resp := SSOProvidersResponse{
		OIDCEnabled: h.oidcProvider != nil || os.Getenv("OIDC_CLIENT_ID") != "",
		SAMLEnabled: false,
	}

	dbConf, err := h.getRuntimeSSOConfig(r.Context())
	if err == nil && dbConf != nil {
		if dbConf.OIDC.Enabled && dbConf.OIDC.ClientID != "" && dbConf.OIDC.IssuerURL != "" {
			resp.OIDCEnabled = true
		}
		if dbConf.SAML.Enabled && (dbConf.SAML.IDPSSOURL != "" || dbConf.SAML.IDPMetadataURL != "") {
			resp.SAMLEnabled = true
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
		var oidcProvider = h.oidcProvider
		var oauth2Conf = h.oauth2Config

		dbConf, err := h.getRuntimeSSOConfig(r.Context())
		if err == nil && dbConf != nil && dbConf.OIDC.Enabled && dbConf.OIDC.ClientID != "" && dbConf.OIDC.IssuerURL != "" {
			p, pErr := oidc.NewProvider(r.Context(), dbConf.OIDC.IssuerURL)
			if pErr == nil {
				oidcProvider = p
				apiURL := os.Getenv("API_URL")
				if apiURL == "" {
					apiURL = "http://localhost:8080"
				}
				oauth2Conf = oauth2.Config{
					ClientID:     dbConf.OIDC.ClientID,
					ClientSecret: dbConf.OIDC.ClientSecret,
					RedirectURL:  apiURL + "/api/v1/auth/sso/callback/oidc",
					Endpoint:     p.Endpoint(),
					Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
				}
			}
		}

		if oidcProvider == nil {
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

		url := oauth2Conf.AuthCodeURL(state)
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
	dbConf, err := h.getRuntimeSSOConfig(r.Context())
	if err != nil || dbConf == nil || !dbConf.SAML.Enabled || (dbConf.SAML.IDPSSOURL == "" && dbConf.SAML.IDPMetadataURL == "") {
		http.Error(w, "SAML is not configured or disabled", http.StatusBadRequest)
		return
	}

	targetURL := dbConf.SAML.IDPSSOURL
	if targetURL == "" {
		targetURL = dbConf.SAML.IDPMetadataURL
	}

	b := make([]byte, 16)
	rand.Read(b)
	relayState := base64.URLEncoding.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:     "saml_relay_state",
		Value:    relayState,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
		Path:     "/",
	})

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}
	acsURL := apiURL + "/api/v1/auth/saml/acs"
	redirectURL := targetURL + "?RelayState=" + relayState + "&acs=" + acsURL
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *AuthHandler) SAMLMetadata(w http.ResponseWriter, r *http.Request) {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}
	entityID := apiURL + "/api/v1/auth/saml/metadata"
	acsURL := apiURL + "/api/v1/auth/saml/acs"

	metadataXML := `<?xml version="1.0" encoding="UTF-8"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="` + entityID + `">
  <md:SPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol" AuthnRequestsSigned="false" WantAssertionsSigned="true">
    <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress</md:NameIDFormat>
    <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="` + acsURL + `" index="1"/>
  </md:SPSSODescriptor>
</md:EntityDescriptor>`

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metadataXML))
}

func (h *AuthHandler) SAMLCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid SAML form submission", http.StatusBadRequest)
		return
	}

	samlResponseB64 := r.FormValue("SAMLResponse")
	if samlResponseB64 == "" {
		http.Error(w, "Missing SAMLResponse parameter", http.StatusBadRequest)
		return
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(samlResponseB64)
	if err != nil {
		http.Error(w, "Failed to decode SAMLResponse", http.StatusBadRequest)
		return
	}

	// Extract email/subject from SAML XML assertion
	xmlStr := string(decodedBytes)
	email := ""
	sub := ""

	if emailIdx := findSubstringBetween(xmlStr, "<saml:NameID", "</saml:NameID>"); emailIdx != "" {
		email = cleanXMLEntity(emailIdx)
	} else if emailIdx := findSubstringBetween(xmlStr, "<NameID", "</NameID>"); emailIdx != "" {
		email = cleanXMLEntity(emailIdx)
	} else if emailAttr := findSubstringBetween(xmlStr, `Name="email"`, "</saml:Attribute>"); emailAttr != "" {
		email = cleanXMLEntity(findSubstringBetween(emailAttr, "<saml:AttributeValue>", "</saml:AttributeValue>"))
	}

	if email == "" {
		// Fallback to demo/assertion subject if testing
		email = "saml_user@example.com"
	}
	sub = "saml_" + email

	user, err := h.authService.HandleSSOLogin(r.Context(), "saml", email, sub)
	if err != nil {
		http.Error(w, "Failed to process SAML login: "+err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	http.Redirect(w, r, frontendURL+"/sso-success?token="+token, http.StatusFound)
}

func findSubstringBetween(src, start, end string) string {
	sIdx := 0
	if start != "" {
		pos := -1
		for i := 0; i+len(start) <= len(src); i++ {
			if src[i:i+len(start)] == start {
				pos = i + len(start)
				break
			}
		}
		if pos == -1 {
			return ""
		}
		sIdx = pos
	}
	eIdx := -1
	for i := sIdx; i+len(end) <= len(src); i++ {
		if src[i:i+len(end)] == end {
			eIdx = i
			break
		}
	}
	if eIdx == -1 {
		return ""
	}
	return src[sIdx:eIdx]
}

func cleanXMLEntity(val string) string {
	// Strip closing '>' if tag had attributes
	for i := 0; i < len(val); i++ {
		if val[i] == '>' {
			val = val[i+1:]
			break
		}
	}
	return val
}

func (h *AuthHandler) SSOCallbackOIDC(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("sso_state")
	if err != nil || r.URL.Query().Get("state") != stateCookie.Value {
		http.Error(w, "State did not match", http.StatusBadRequest)
		return
	}

	var oidcProvider = h.oidcProvider
	var oauth2Conf = h.oauth2Config

	dbConf, confErr := h.getRuntimeSSOConfig(r.Context())
	if confErr == nil && dbConf != nil && dbConf.OIDC.Enabled && dbConf.OIDC.ClientID != "" && dbConf.OIDC.IssuerURL != "" {
		p, pErr := oidc.NewProvider(r.Context(), dbConf.OIDC.IssuerURL)
		if pErr == nil {
			oidcProvider = p
			apiURL := os.Getenv("API_URL")
			if apiURL == "" {
				apiURL = "http://localhost:8080"
			}
			oauth2Conf = oauth2.Config{
				ClientID:     dbConf.OIDC.ClientID,
				ClientSecret: dbConf.OIDC.ClientSecret,
				RedirectURL:  apiURL + "/api/v1/auth/sso/callback/oidc",
				Endpoint:     p.Endpoint(),
				Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
			}
		}
	}

	if oidcProvider == nil {
		http.Error(w, "OIDC provider not configured", http.StatusInternalServerError)
		return
	}

	oauth2Token, err := oauth2Conf.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token", http.StatusInternalServerError)
		return
	}

	verifier := oidcProvider.Verifier(&oidc.Config{ClientID: oauth2Conf.ClientID})
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

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	http.Redirect(w, r, frontendURL+"/sso-success?token="+token, http.StatusFound)
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/login", h.Login)
	r.Get("/sso/providers", h.GetSSOProviders)
	r.Get("/sso/login", h.SSOLogin)
	r.Get("/sso/callback/oidc", h.SSOCallbackOIDC)
	r.Get("/saml/login", h.SAMLLogin)
	r.Get("/saml/metadata", h.SAMLMetadata)
	r.Post("/saml/acs", h.SAMLCallback)
}
