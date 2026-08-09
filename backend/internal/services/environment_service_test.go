package services

import (
	"context"
	"testing"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type EnvMockStore struct {
	mock.Mock
	envRepo       *EnvMockRepo
	flagStateRepo *EnvMockFlagStateRepo
	auditRepo     *EnvMockAuditRepo
}

func (m *EnvMockStore) ProjectRepo() repository.ProjectRepository             { return nil }
func (m *EnvMockStore) EnvironmentRepo() repository.EnvironmentRepository     { return m.envRepo }
func (m *EnvMockStore) FlagRepo() repository.FlagRepository                     { return nil }
func (m *EnvMockStore) FlagStateRepo() repository.FlagStateRepository           { return m.flagStateRepo }
func (m *EnvMockStore) AuditRepo() repository.AuditRepository                   { return m.auditRepo }
func (m *EnvMockStore) ChangeRequestRepo() repository.ChangeRequestRepository { return nil }
func (m *EnvMockStore) KillSwitchRepo() repository.KillSwitchRepository       { return nil }
func (m *EnvMockStore) SlackConfigRepo() repository.SlackConfigRepository       { return nil }
func (m *EnvMockStore) ScheduledChangeRepo() repository.ScheduledChangeRepository { return nil }
func (m *EnvMockStore) StalePolicyRepo() repository.StalePolicyRepository       { return nil }
func (m *EnvMockStore) RoleRepo() repository.RoleRepository                     { return nil }
func (m *EnvMockStore) UserRepo() repository.UserRepository                     { return nil }
func (m *EnvMockStore) ServiceAccountRepo() repository.ServiceAccountRepository { return nil }
func (m *EnvMockStore) WebhookIntegrationRepo() repository.WebhookIntegrationRepository { return nil }
func (m *EnvMockStore) WithTx(ctx context.Context, fn func(repository.Store) error) error {
	return fn(m)
}
func (m *EnvMockStore) MigrateUp() error   { return nil }
func (m *EnvMockStore) MigrateDown() error { return nil }

type EnvMockRepo struct {
	mock.Mock
}

func (m *EnvMockRepo) Create(ctx context.Context, env *models.Environment) error {
	args := m.Called(ctx, env)
	return args.Error(0)
}
func (m *EnvMockRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Environment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Environment), args.Error(1)
}
func (m *EnvMockRepo) GetByAPIKeyHash(ctx context.Context, hash string) (*models.Environment, error) {
	return nil, nil
}
func (m *EnvMockRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.Environment, error) {
	return nil, nil
}
func (m *EnvMockRepo) Update(ctx context.Context, env *models.Environment) error {
	return nil
}
func (m *EnvMockRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type EnvMockFlagStateRepo struct {
	mock.Mock
}

func (m *EnvMockFlagStateRepo) Create(ctx context.Context, state *models.EnvironmentFlagState) error {
	return nil
}
func (m *EnvMockFlagStateRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.EnvironmentFlagState, error) {
	return nil, nil
}
func (m *EnvMockFlagStateRepo) GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*models.EnvironmentFlagState, error) {
	return nil, nil
}
func (m *EnvMockFlagStateRepo) ListByEnvironment(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentFlagState, error) {
	return nil, nil
}
func (m *EnvMockFlagStateRepo) ListByEnvironmentAndLifecycle(ctx context.Context, envID uuid.UUID, l models.FlagLifecycleState) ([]*models.EnvironmentFlagState, error) {
	return nil, nil
}
func (m *EnvMockFlagStateRepo) Update(ctx context.Context, state *models.EnvironmentFlagState) error {
	return nil
}
func (m *EnvMockFlagStateRepo) UpdateLifecycleState(ctx context.Context, id uuid.UUID, state models.FlagLifecycleState) error {
	return nil
}
func (m *EnvMockFlagStateRepo) UpdateLastEvaluatedAtBatch(ctx context.Context, updates map[uuid.UUID]time.Time) error {
	return nil
}
func (m *EnvMockFlagStateRepo) FindActiveFlagsForStalenessScan(ctx context.Context, limit int) ([]*models.EnvironmentFlagState, error) {
	return nil, nil
}
func (m *EnvMockFlagStateRepo) CloneEnvironmentState(ctx context.Context, sourceEnvID, targetEnvID uuid.UUID) error {
	args := m.Called(ctx, sourceEnvID, targetEnvID)
	return args.Error(0)
}
func (m *EnvMockFlagStateRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type EnvMockAuditRepo struct {
	mock.Mock
}

func (m *EnvMockAuditRepo) Create(ctx context.Context, log *models.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}
func (m *EnvMockAuditRepo) ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error) {
	return nil, 0, nil
}
func (m *EnvMockAuditRepo) ListByEnvironment(ctx context.Context, envID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error) {
	return nil, 0, nil
}
func (m *EnvMockAuditRepo) DeleteOlderThan(ctx context.Context, cutoff time.Time) error {
	return nil
}

func TestCloneEnvironment_Success(t *testing.T) {
	store := new(EnvMockStore)
	envRepo := new(EnvMockRepo)
	flagStateRepo := new(EnvMockFlagStateRepo)
	auditRepo := new(EnvMockAuditRepo)

	store.envRepo = envRepo
	store.flagStateRepo = flagStateRepo
	store.auditRepo = auditRepo

	auditService := NewAuditService(store)
	service := NewEnvironmentService(store, auditService)

	projectID := uuid.New()
	sourceEnvID := uuid.New()
	actorID := uuid.New()

	sourceEnv := &models.Environment{
		ID:          sourceEnvID,
		ProjectID:   projectID,
		Name:        "Staging",
		IsProtected: false,
	}

	envRepo.On("GetByID", mock.Anything, sourceEnvID).Return(sourceEnv, nil)
	envRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Environment")).Return(nil)
	flagStateRepo.On("CloneEnvironmentState", mock.Anything, sourceEnvID, mock.AnythingOfType("uuid.UUID")).Return(nil)
	auditRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.AuditLog")).Return(nil).Maybe()

	newEnv, apiKey, err := service.CloneEnvironment(context.Background(), projectID, sourceEnvID, "PR-123-Test", actorID, "127.0.0.1")
	time.Sleep(10 * time.Millisecond)

	assert.NoError(t, err)
	assert.NotNil(t, newEnv)
	assert.Equal(t, "PR-123-Test", newEnv.Name)
	assert.Equal(t, projectID, newEnv.ProjectID)
	assert.False(t, newEnv.IsProtected)
	assert.True(t, len(apiKey) > 0)
	assert.Contains(t, apiKey, "env_")

	envRepo.AssertExpectations(t)
	flagStateRepo.AssertExpectations(t)
	auditRepo.AssertExpectations(t)
}

func TestDeleteEnvironment_Protected_ReturnsForbidden(t *testing.T) {
	store := new(EnvMockStore)
	envRepo := new(EnvMockRepo)
	store.envRepo = envRepo

	auditService := NewAuditService(store)
	service := NewEnvironmentService(store, auditService)

	envID := uuid.New()
	actorID := uuid.New()

	protectedEnv := &models.Environment{
		ID:          envID,
		ProjectID:   uuid.New(),
		Name:        "Production",
		IsProtected: true,
	}

	envRepo.On("GetByID", mock.Anything, envID).Return(protectedEnv, nil)

	err := service.DeleteEnvironment(context.Background(), envID, actorID, "127.0.0.1")

	assert.ErrorIs(t, err, ErrProtectedEnvironment)
	envRepo.AssertExpectations(t)
}

func TestDeleteEnvironment_Success(t *testing.T) {
	store := new(EnvMockStore)
	envRepo := new(EnvMockRepo)
	auditRepo := new(EnvMockAuditRepo)
	store.envRepo = envRepo
	store.auditRepo = auditRepo

	auditService := NewAuditService(store)
	service := NewEnvironmentService(store, auditService)

	envID := uuid.New()
	actorID := uuid.New()

	ephemeralEnv := &models.Environment{
		ID:          envID,
		ProjectID:   uuid.New(),
		Name:        "PR-123-Test",
		IsProtected: false,
	}

	envRepo.On("GetByID", mock.Anything, envID).Return(ephemeralEnv, nil)
	envRepo.On("Delete", mock.Anything, envID).Return(nil)
	auditRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.AuditLog")).Return(nil).Maybe()

	err := service.DeleteEnvironment(context.Background(), envID, actorID, "127.0.0.1")
	time.Sleep(10 * time.Millisecond)

	assert.NoError(t, err)
	envRepo.AssertExpectations(t)
	auditRepo.AssertExpectations(t)
}
