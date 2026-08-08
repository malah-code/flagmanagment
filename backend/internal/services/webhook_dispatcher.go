package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/google/uuid"
)

type WebhookDispatcher struct {
	repo       repository.WebhookIntegrationRepository
	auditChan  <-chan *models.AuditLog
	httpClient *http.Client
}

func NewWebhookDispatcher(repo repository.WebhookIntegrationRepository, auditChan <-chan *models.AuditLog) *WebhookDispatcher {
	return &WebhookDispatcher{
		repo:       repo,
		auditChan:  auditChan,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (d *WebhookDispatcher) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case log, ok := <-d.auditChan:
			if !ok {
				return
			}
			if log != nil && log.ProjectID != nil {
				go d.dispatchLog(ctx, *log.ProjectID, log)
			}
		}
	}
}

func (d *WebhookDispatcher) dispatchLog(ctx context.Context, projectID uuid.UUID, log *models.AuditLog) {
	hooks, err := d.repo.ListByProject(ctx, projectID)
	if err != nil || len(hooks) == 0 {
		return
	}

	payload := models.WebhookPayload{
		EventID:   uuid.New().String(),
		EventType: "audit." + log.Action,
		Timestamp: time.Now(),
		Data:      log,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	for _, wh := range hooks {
		go d.sendWithRetry(wh.URL, body, 3, 1*time.Second)
	}
}

func (d *WebhookDispatcher) sendWithRetry(url string, body []byte, retries int, backoff time.Duration) {
	for i := 0; i < retries; i++ {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			resp, err := d.httpClient.Do(req)
			if err == nil {
				_ = resp.Body.Close()
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					return // Success
				}
			}
		}
		time.Sleep(backoff)
		backoff *= 2
	}
}
