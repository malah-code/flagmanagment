package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type WebhookMockKillSwitchRepo struct {
	mock.Mock
}

func (m *WebhookMockKillSwitchRepo) Create(ctx context.Context, rule *models.KillSwitchRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *WebhookMockKillSwitchRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.KillSwitchRule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KillSwitchRule), args.Error(1)
}

func (m *WebhookMockKillSwitchRepo) ListByEnvironmentAndFlag(ctx context.Context, envID, flagID uuid.UUID) ([]*models.KillSwitchRule, error) {
	args := m.Called(ctx, envID, flagID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.KillSwitchRule), args.Error(1)
}

func (m *WebhookMockKillSwitchRepo) ListByEnvironmentAndAlert(ctx context.Context, envID uuid.UUID, alertIdentifier string) ([]*models.KillSwitchRule, error) {
	args := m.Called(ctx, envID, alertIdentifier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.KillSwitchRule), args.Error(1)
}

func (m *WebhookMockKillSwitchRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type WebhookMockFlagStateRepo struct {
	mock.Mock
}

func (m *WebhookMockFlagStateRepo) GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, envID, flagID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EnvironmentFlagState), args.Error(1)
}

func (m *WebhookMockFlagStateRepo) Create(ctx context.Context, state *models.EnvironmentFlagState) error {
	args := m.Called(ctx, state)
	return args.Error(0)
}

func (m *WebhookMockFlagStateRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EnvironmentFlagState), args.Error(1)
}

func (m *WebhookMockFlagStateRepo) Update(ctx context.Context, state *models.EnvironmentFlagState) error {
	args := m.Called(ctx, state)
	return args.Error(0)
}

func (m *WebhookMockFlagStateRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *WebhookMockFlagStateRepo) ListByEnvironment(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, envID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.EnvironmentFlagState), args.Error(1)
}

func (m *WebhookMockFlagStateRepo) ListByEnvironmentAndLifecycle(ctx context.Context, envID uuid.UUID, lifecycle models.FlagLifecycleState) ([]*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, envID, lifecycle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.EnvironmentFlagState), args.Error(1)
}

func (m *WebhookMockFlagStateRepo) UpdateLifecycleState(ctx context.Context, id uuid.UUID, state models.FlagLifecycleState) error {
	args := m.Called(ctx, id, state)
	return args.Error(0)
}

func (m *WebhookMockFlagStateRepo) UpdateLastEvaluatedAtBatch(ctx context.Context, updates map[uuid.UUID]time.Time) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func (m *WebhookMockFlagStateRepo) FindActiveFlagsForStalenessScan(ctx context.Context, limit int) ([]*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.EnvironmentFlagState), args.Error(1)
}

type WebhookMockAuditRepo struct {
	mock.Mock
}

func (m *WebhookMockAuditRepo) Create(ctx context.Context, log *models.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *WebhookMockAuditRepo) ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error) {
	args := m.Called(ctx, projectID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*models.AuditLog), args.Int(1), args.Error(2)
}

func (m *WebhookMockAuditRepo) ListByEnvironment(ctx context.Context, envID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error) {
	args := m.Called(ctx, envID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*models.AuditLog), args.Int(1), args.Error(2)
}

type WebhookMockStore struct {
	ksRepo    *WebhookMockKillSwitchRepo
	fsRepo    *WebhookMockFlagStateRepo
	auditRepo *WebhookMockAuditRepo
}

func (m *WebhookMockStore) KillSwitchRepo() repository.KillSwitchRepository { return m.ksRepo }
func (m *WebhookMockStore) FlagStateRepo() repository.FlagStateRepository { return m.fsRepo }
func (m *WebhookMockStore) AuditRepo() repository.AuditRepository         { return m.auditRepo }
func (m *WebhookMockStore) ProjectRepo() repository.ProjectRepository     { return nil }
func (m *WebhookMockStore) EnvironmentRepo() repository.EnvironmentRepository { return nil }
func (m *WebhookMockStore) FlagRepo() repository.FlagRepository           { return nil }
func (m *WebhookMockStore) ChangeRequestRepo() repository.ChangeRequestRepository { return nil }
func (m *WebhookMockStore) RoleRepo() repository.RoleRepository           { return nil }
func (m *WebhookMockStore) UserRepo() repository.UserRepository           { return nil }
func (m *WebhookMockStore) SlackConfigRepo() repository.SlackConfigRepository { return nil }
func (m *WebhookMockStore) ScheduledChangeRepo() repository.ScheduledChangeRepository { return nil }
func (m *WebhookMockStore) StalePolicyRepo() repository.StalePolicyRepository { return nil }
func (m *WebhookMockStore) WithTx(ctx context.Context, fn func(repository.Store) error) error {
	return fn(m)
}
func (m *WebhookMockStore) MigrateUp() error   { return nil }
func (m *WebhookMockStore) MigrateDown() error { return nil }

func TestWebhookService_ProcessAPMAlert_DisablesFlag(t *testing.T) {
	ksRepo := new(WebhookMockKillSwitchRepo)
	fsRepo := new(WebhookMockFlagStateRepo)
	auditRepo := new(WebhookMockAuditRepo)

	store := &WebhookMockStore{
		ksRepo:    ksRepo,
		fsRepo:    fsRepo,
		auditRepo: auditRepo,
	}

	auditSvc := services.NewAuditService(store)
	svc := services.NewWebhookService(store, auditSvc, nil, nil)

	envID := uuid.New()
	flagID := uuid.New()
	alertID := "high_error_rate"

	ksRule := &models.KillSwitchRule{
		ID:              uuid.New(),
		FlagID:          flagID,
		EnvironmentID:   envID,
		AlertIdentifier: alertID,
		Action:          "DISABLE",
	}

	flagState := &models.EnvironmentFlagState{
		ID:            uuid.New(),
		EnvironmentID: envID,
		FeatureFlagID: flagID,
		Enabled:       true,
	}

	ksRepo.On("ListByEnvironmentAndAlert", mock.Anything, envID, alertID).Return([]*models.KillSwitchRule{ksRule}, nil)
	fsRepo.On("GetByEnvAndFlag", mock.Anything, envID, flagID).Return(flagState, nil)
	fsRepo.On("Update", mock.Anything, mock.MatchedBy(func(s *models.EnvironmentFlagState) bool {
		return s.Enabled == false
	})).Return(nil)
	auditRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	_, err := svc.ProcessAPMAlert(context.Background(), envID, alertID, map[string]string{"alert_identifier": alertID})

	assert.NoError(t, err)
	ksRepo.AssertExpectations(t)
	fsRepo.AssertExpectations(t)
	auditRepo.AssertExpectations(t)
}
