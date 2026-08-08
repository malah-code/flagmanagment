package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flagmanagment/backend/internal/api"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWebhookService struct {
	mock.Mock
}

func (m *MockWebhookService) ProcessAPMAlert(ctx context.Context, envID uuid.UUID, alertIdentifier string, payload interface{}) (int, error) {
	args := m.Called(ctx, envID, alertIdentifier, payload)
	return args.Int(0), args.Error(1)
}

func (m *MockWebhookService) CreateWebhook(ctx context.Context, wh *models.WebhookIntegration) error {
	args := m.Called(ctx, wh)
	return args.Error(0)
}

func (m *MockWebhookService) ListWebhooks(ctx context.Context, projectID uuid.UUID) ([]*models.WebhookIntegration, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]*models.WebhookIntegration), args.Error(1)
}

func (m *MockWebhookService) DeleteWebhook(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestHandleAPMWebhook_Unauthorized(t *testing.T) {
	mockSvc := new(MockWebhookService)
	handler := api.NewWebhookHandler(mockSvc)

	reqPayload := map[string]string{
		"alert_identifier": "high_cpu",
	}
	body, _ := json.Marshal(reqPayload)

	req, _ := http.NewRequest("POST", "/api/v1/webhooks/apm", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.HandleAPMWebhook(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandleAPMWebhook_Success(t *testing.T) {
	mockSvc := new(MockWebhookService)
	handler := api.NewWebhookHandler(mockSvc)

	envID := uuid.New()
	env := &models.Environment{
		ID:   envID,
		Name: "Production",
	}

	reqPayload := map[string]string{
		"alert_identifier": "high_error_rate",
	}
	body, _ := json.Marshal(reqPayload)

	req, _ := http.NewRequest("POST", "/api/v1/webhooks/apm", bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), api.EnvKey, env)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	mockSvc.On("ProcessAPMAlert", mock.Anything, envID, "high_error_rate", mock.Anything).Return(1, nil)

	handler.HandleAPMWebhook(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code)
	mockSvc.AssertExpectations(t)
}
