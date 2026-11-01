package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	sdk "github.com/malah-code/flagmanagment/sdk/go"
	"github.com/open-feature/go-sdk/openfeature"
)

type Server struct {
	client *openfeature.Client
}

func main() {
	apiKey := os.Getenv("FLAGMANAGMENT_SERVER_KEY")
	if apiKey == "" {
		apiKey = "sample-go-server-key"
	}

	apiURL := os.Getenv("FLAGMANAGMENT_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	// Initialize the FlagManagment OpenFeature Go Provider
	provider := sdk.NewProvider(apiKey, sdk.Options{
		StreamURL:   fmt.Sprintf("%s/api/v1/sdk/stream", apiURL),
		EndpointURL: apiURL,
	})
	defer provider.Shutdown()

	openfeature.SetProvider(provider)
	ofClient := openfeature.NewClient("go-http-service")

	s := &Server{client: ofClient}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/evaluate", s.handleEvaluate)
	mux.HandleFunc("/api/checkout", s.handleCheckout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("🚀 FlagManagment Go HTTP Sample starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server exited: %v", err)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "go-http-sample",
	})
}

func (s *Server) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		userId = "user-go-default"
	}
	tier := r.URL.Query().Get("tier")
	if tier == "" {
		tier = "standard"
	}

	evalCtx := openfeature.NewEvaluationContext(
		userId,
		map[string]interface{}{
			"tier": tier,
		},
	)

	// Sub-microsecond local in-memory evaluation
	isNewCheckout, err := s.client.BooleanValue(r.Context(), "new-checkout-v2", false, evalCtx)
	if err != nil {
		log.Printf("Error evaluating flag: %v", err)
	}

	flow := "v1-standard"
	if isNewCheckout {
		flow = "v2-one-click"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"userId":         userId,
		"tier":           tier,
		"isNewCheckout":  isNewCheckout,
		"checkoutFlow":   flow,
	})
}

func (s *Server) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	flagKey := r.URL.Query().Get("flag")
	if flagKey == "" {
		flagKey = "new-checkout-v2"
	}
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		userId = "anon-user"
	}

	evalCtx := openfeature.NewEvaluationContext(userId, nil)

	details, err := s.client.BooleanValueDetails(r.Context(), flagKey, false, evalCtx)
	if err != nil {
		log.Printf("Evaluation details error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"flagKey": flagKey,
		"value":   details.Value,
		"variant": details.Variant,
		"reason":  string(details.Reason),
	})
}
