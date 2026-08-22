package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/open-feature/go-sdk/openfeature"
)

// PostHogHook implements the openfeature.Hook interface to send evaluation events to PostHog.
type PostHogHook struct {
	posthogAPIKey string
	endpoint      string
	client        *http.Client
}

// NewPostHogHook creates a new PostHogHook with the provided API key.
// By default it sends to https://us.i.posthog.com/capture/
func NewPostHogHook(posthogAPIKey string, endpoint string) *PostHogHook {
	if endpoint == "" {
		endpoint = "https://us.i.posthog.com/capture/"
	}
	return &PostHogHook{
		posthogAPIKey: posthogAPIKey,
		endpoint:      endpoint,
		client:        &http.Client{Timeout: 5 * time.Second},
	}
}

func (h *PostHogHook) Before(ctx context.Context, hookContext openfeature.HookContext, hookHints openfeature.HookHints) (*openfeature.EvaluationContext, error) {
	return nil, nil
}

func (h *PostHogHook) After(ctx context.Context, hookContext openfeature.HookContext, flagEvaluationDetails openfeature.InterfaceResolutionDetail, hookHints openfeature.HookHints) error {
	distinctID := hookContext.EvaluationContext().TargetingKey()
	if distinctID == "" {
		distinctID = "anonymous-server-user"
	}

	payload := map[string]interface{}{
		"api_key": h.posthogAPIKey,
		"event":   "$feature_flag_called",
		"properties": map[string]interface{}{
			"distinct_id":            distinctID,
			"$feature_flag":          hookContext.FlagKey(),
			"$feature_flag_response": flagEvaluationDetails.Value,
			"timestamp":              time.Now().UTC().Format(time.RFC3339),
		},
	}

	go h.sendEvent(payload)
	return nil
}

func (h *PostHogHook) Error(ctx context.Context, hookContext openfeature.HookContext, err error, hookHints openfeature.HookHints) {
	distinctID := hookContext.EvaluationContext().TargetingKey()
	if distinctID == "" {
		distinctID = "anonymous-server-user"
	}

	payload := map[string]interface{}{
		"api_key": h.posthogAPIKey,
		"event":   "Feature Flag Error",
		"properties": map[string]interface{}{
			"distinct_id":   distinctID,
			"$feature_flag": hookContext.FlagKey(),
			"error":         err.Error(),
			"timestamp":     time.Now().UTC().Format(time.RFC3339),
		},
	}

	go h.sendEvent(payload)
}

func (h *PostHogHook) Finally(ctx context.Context, hookContext openfeature.HookContext, hookHints openfeature.HookHints) {
	// No-op
}

func (h *PostHogHook) sendEvent(payload map[string]interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[flagmanagment-posthog] failed to marshal event: %v", err)
		return
	}

	req, err := http.NewRequest("POST", h.endpoint, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := h.client.Do(req)
	if err != nil {
		log.Printf("[flagmanagment-posthog] failed to send event: %v", err)
		return
	}
	defer resp.Body.Close()
}
