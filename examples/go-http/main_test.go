package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/malah-code/flagmanagment/sdk/go"
	"github.com/open-feature/go-sdk/openfeature"
)

func TestEndpoints(t *testing.T) {
	provider := sdk.NewProvider("test-key", sdk.Options{
		StreamURL: "http://localhost:18080/unused",
	})
	openfeature.SetProvider(provider)
	ofClient := openfeature.NewClient("test-client")
	server := &Server{client: ofClient}

	t.Run("Health check", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		server.handleHealth(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var body map[string]string
		json.NewDecoder(w.Body).Decode(&body)
		if body["status"] != "ok" {
			t.Errorf("expected status ok, got %s", body["status"])
		}
	})

	t.Run("Checkout flag fallback evaluation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/checkout?userId=u1&tier=gold", nil)
		w := httptest.NewRecorder()
		server.handleCheckout(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var body map[string]interface{}
		json.NewDecoder(w.Body).Decode(&body)
		if body["checkoutFlow"] != "v1-standard" {
			t.Errorf("expected v1-standard fallback, got %v", body["checkoutFlow"])
		}
	})
}
