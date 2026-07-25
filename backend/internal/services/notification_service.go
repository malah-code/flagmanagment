package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/repository"
)

type NotificationService interface {
	SendFlagStateChanged(ctx context.Context, envID uuid.UUID, flagKey string, oldState, newState bool, actor string)
}

type notificationService struct {
	store      repository.Store
	httpClient *http.Client
}

func NewNotificationService(store repository.Store) NotificationService {
	return &notificationService{
		store: store,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type SlackBlock struct {
	Type string     `json:"type"`
	Text *SlackText `json:"text,omitempty"`
}

type SlackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (s *notificationService) SendFlagStateChanged(ctx context.Context, envID uuid.UUID, flagKey string, oldState, newState bool, actor string) {
	// Dispatch asynchronously to avoid blocking API caller
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		config, err := s.store.SlackConfigRepo().GetByEnvironmentID(bgCtx, envID)
		if err != nil || config == nil || !config.Enabled || config.WebhookURL == "" {
			return
		}

		statusStr := "DISABLED 🔴"
		if newState {
			statusStr = "ENABLED 🟢"
		}

		text := fmt.Sprintf("🚩 *Feature Flag Update*\n*Flag:* `%s`\n*Status:* %s (was %t)\n*Actor:* %s\n*Time:* %s",
			flagKey, statusStr, oldState, actor, time.Now().UTC().Format(time.RFC3339))

		payload := map[string]interface{}{
			"text": text,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return
		}

		req, err := http.NewRequestWithContext(bgCtx, "POST", config.WebhookURL, bytes.NewBuffer(body))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.httpClient.Do(req)
		if err == nil {
			_ = resp.Body.Close()
		}
	}()
}
