package testutil

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
)

// MockWebhookServer captures received HTTP webhook payloads for testing Slack and SIEM integrations.
type MockWebhookServer struct {
	server   *httptest.Server
	mu       sync.Mutex
	payloads [][]byte
}

// NewMockWebhookServer creates and starts an httptest.Server capturing payloads.
func NewMockWebhookServer() *MockWebhookServer {
	m := &MockWebhookServer{
		payloads: make([][]byte, 0),
	}
	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, r.ContentLength)
		if r.ContentLength > 0 {
			_, _ = r.Body.Read(buf)
		}
		m.mu.Lock()
		m.payloads = append(m.payloads, buf)
		m.mu.Unlock()

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	return m
}

func (m *MockWebhookServer) URL() string {
	return m.server.URL
}

func (m *MockWebhookServer) Close() {
	m.server.Close()
}

func (m *MockWebhookServer) PayloadCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.payloads)
}

func (m *MockWebhookServer) Payloads() [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make([][]byte, len(m.payloads))
	copy(cp, m.payloads)
	return cp
}

// MockOIDCServer provides an in-memory OIDC discovery and JWKS endpoint for testing SSO.
type MockOIDCServer struct {
	server *httptest.Server
	key    *rsa.PrivateKey
}

// NewMockOIDCServer creates an in-memory OIDC provider.
func NewMockOIDCServer() *MockOIDCServer {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	m := &MockOIDCServer{key: key}

	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"issuer":                 m.server.URL,
			"authorization_endpoint": m.server.URL + "/auth",
			"token_endpoint":         m.server.URL + "/token",
			"jwks_uri":               m.server.URL + "/jwks",
			"userinfo_endpoint":      m.server.URL + "/userinfo",
		})
	})

	mux.HandleFunc("/jwks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"keys": []map[string]string{
				{
					"kty": "RSA",
					"alg": "RS256",
					"use": "sig",
					"kid": "test-key-id",
					"n":   "test-modulus",
					"e":   "AQAB",
				},
			},
		})
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "mock-access-token",
			"id_token":     "mock-id-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	})

	m.server = httptest.NewServer(mux)
	return m
}

func (m *MockOIDCServer) URL() string {
	return m.server.URL
}

func (m *MockOIDCServer) Close() {
	m.server.Close()
}

