package services_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
	"github.com/google/uuid"
)

type mockStore struct {
	repository.Store
	auditRepo *mockAuditRepo
}

func (m *mockStore) AuditRepo() repository.AuditRepository {
	return m.auditRepo
}

type mockAuditRepo struct {
	repository.AuditRepository
	mu   sync.Mutex
	logs []*models.AuditLog
}

func (m *mockAuditRepo) Create(ctx context.Context, log *models.AuditLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, log)
	return nil
}

func TestAuditService_LogAction_SanitizesSensitiveData(t *testing.T) {
	mockRepo := &mockAuditRepo{}
	store := &mockStore{auditRepo: mockRepo}
	auditService := services.NewAuditService(store)

	rawNewState := models.JSONB{
		"name":    "Production Env",
		"api_key": "secret-12345",
		"token":   "bearer-xyz",
		"secret":  "my-super-secret",
	}

	rawPrevState := models.JSONB{
		"password": "old-password",
		"normal":   "value",
	}

	log := &models.AuditLog{
		ID:            uuid.New(),
		ActorID:       uuid.New(),
		Action:        "CREATE",
		TargetType:    "ENVIRONMENT",
		TargetID:      uuid.New(),
		PreviousState: rawPrevState,
		NewState:      rawNewState,
		CreatedAt:     time.Now().UTC(),
	}

	err := auditService.LogAction(context.Background(), log)
	if err != nil {
		t.Fatalf("unexpected error logging action: %v", err)
	}

	// Give async goroutine a brief moment to execute
	time.Sleep(50 * time.Millisecond)

	mockRepo.mu.Lock()
	defer mockRepo.mu.Unlock()

	if len(mockRepo.logs) != 1 {
		t.Fatalf("expected 1 audit log, got %d", len(mockRepo.logs))
	}

	createdLog := mockRepo.logs[0]

	// Verify sensitive keys were redacted
	if createdLog.NewState["api_key"] != "***" {
		t.Errorf("expected api_key to be redacted to '***', got %v", createdLog.NewState["api_key"])
	}
	if createdLog.NewState["token"] != "***" {
		t.Errorf("expected token to be redacted to '***', got %v", createdLog.NewState["token"])
	}
	if createdLog.NewState["secret"] != "***" {
		t.Errorf("expected secret to be redacted to '***', got %v", createdLog.NewState["secret"])
	}
	if createdLog.NewState["name"] != "Production Env" {
		t.Errorf("expected name to be 'Production Env', got %v", createdLog.NewState["name"])
	}

	if createdLog.PreviousState["password"] != "***" {
		t.Errorf("expected password to be redacted to '***', got %v", createdLog.PreviousState["password"])
	}
	if createdLog.PreviousState["normal"] != "value" {
		t.Errorf("expected normal to be 'value', got %v", createdLog.PreviousState["normal"])
	}
}
