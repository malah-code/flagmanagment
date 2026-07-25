package services_test

import (
	"context"
	"testing"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKillSwitchRepo struct {
	mock.Mock
}

func (m *MockKillSwitchRepo) Create(ctx context.Context, rule *models.KillSwitchRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockKillSwitchRepo) ListByEnvironmentAndFlag(ctx context.Context, envID, flagID uuid.UUID) ([]*models.KillSwitchRule, error) {
	args := m.Called(ctx, envID, flagID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.KillSwitchRule), args.Error(1)
}

func (m *MockKillSwitchRepo) ListByEnvironmentAndAlert(ctx context.Context, envID uuid.UUID, alertIdentifier string) ([]*models.KillSwitchRule, error) {
	args := m.Called(ctx, envID, alertIdentifier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.KillSwitchRule), args.Error(1)
}

func (m *MockKillSwitchRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockFlagStateRepo struct {
	mock.Mock
}

func (m *MockFlagStateRepo) GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, envID, flagID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EnvironmentFlagState), args.Error(1)
}

func (m *MockFlagStateRepo) Create(ctx context.Context, state *models.EnvironmentFlagState) error {
	args := m.Called(ctx, state)
	return args.Error(0)
}

func (m *MockFlagStateRepo) Update(ctx context.Context, state *models.EnvironmentFlagState) error {
	args := m.Called(ctx, state)
	return args.Error(0)
}

func (m *MockFlagStateRepo) ListByEnvironment(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, envID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.EnvironmentFlagState), args.Error(1)
}

type MockAuditRepo struct {
	mock.Mock
}

func (m *MockAuditRepo) Create(ctx context.Context, log *models.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockAuditRepo) List(ctx context.Context, filter repository.AuditFilter, limit, offset int) ([]*models.AuditLog, int, error) {
	args := m.Called(ctx, filter, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*models.AuditLog), args.Int(1), args.Error(2)
}

type WebhookMockStore struct {
	mockStore MockStore
	ksRepo    *MockKillSwitchRepo
	fsRepo    *MockFlagStateRepo
	auditRepo *MockAuditRepo
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
func (m *WebhookMockStore) WithTx(ctx context.Context, fn func(repository.Store) error) error {
	return fn(m)
}
func (m *WebhookMockStore) MigrateUp() error   { return nil }
func (m *WebhookMockStore) MigrateDown() error { return nil }

func TestWebhookService_ProcessAPMAlert_DisablesFlag(t *testing.T) {
	ksRepo := new(MockKillSwitchRepo)
	fsRepo := new(MockFlagStateRepo)
	auditRepo := new(MockAuditRepo)

	store := &WebhookMockStore{
		ksRepo:    ksRepo,
		fsRepo:    fsRepo,
		auditRepo: auditRepo,
	}

	auditSvc := services.NewAuditService(store)
	svc := services.NewWebhookService(store, auditSvc, nil)

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

	err := svc.ProcessAPMAlert(context.Background(), envID, alertID, map[string]string{"alert_identifier": alertID})

	assert.NoError(t, err)
	ksRepo.AssertExpectations(t)
	fsRepo.AssertExpectations(t)
	auditRepo.AssertExpectations(t)
}
