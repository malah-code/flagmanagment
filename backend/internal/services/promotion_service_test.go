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

type MockEnvironmentRepo struct {
	mock.Mock
}

func (m *MockEnvironmentRepo) Create(ctx context.Context, env *models.Environment) error {
	args := m.Called(ctx, env)
	return args.Error(0)
}

func (m *MockEnvironmentRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Environment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Environment), args.Error(1)
}

func (m *MockEnvironmentRepo) GetByAPIKeyHash(ctx context.Context, hash string) (*models.Environment, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Environment), args.Error(1)
}

func (m *MockEnvironmentRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.Environment, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Environment), args.Error(1)
}

func (m *MockEnvironmentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEnvironmentRepo) Update(ctx context.Context, env *models.Environment) error {
	args := m.Called(ctx, env)
	return args.Error(0)
}

type MockFlagStateRepo struct {
	mock.Mock
}

func (m *MockFlagStateRepo) Create(ctx context.Context, state *models.EnvironmentFlagState) error {
	args := m.Called(ctx, state)
	return args.Error(0)
}

func (m *MockFlagStateRepo) GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, envID, flagID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EnvironmentFlagState), args.Error(1)
}

func (m *MockFlagStateRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EnvironmentFlagState), args.Error(1)
}

func (m *MockFlagStateRepo) Update(ctx context.Context, state *models.EnvironmentFlagState) error {
	args := m.Called(ctx, state)
	return args.Error(0)
}

func (m *MockFlagStateRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFlagStateRepo) ListByEnvironment(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, envID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.EnvironmentFlagState), args.Error(1)
}

func (m *MockFlagStateRepo) ListByEnvironmentAndLifecycle(ctx context.Context, envID uuid.UUID, lifecycle models.FlagLifecycleState) ([]*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, envID, lifecycle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.EnvironmentFlagState), args.Error(1)
}

func (m *MockFlagStateRepo) UpdateLifecycleState(ctx context.Context, id uuid.UUID, state models.FlagLifecycleState) error {
	args := m.Called(ctx, id, state)
	return args.Error(0)
}

func (m *MockFlagStateRepo) UpdateLastEvaluatedAtBatch(ctx context.Context, updates map[uuid.UUID]time.Time) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func (m *MockFlagStateRepo) FindActiveFlagsForStalenessScan(ctx context.Context, limit int) ([]*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.EnvironmentFlagState), args.Error(1)
}

type PromoMockStore struct {
	mock.Mock
	crRepo *MockChangeRequestRepo
	envRepo *MockEnvironmentRepo
	fsRepo *MockFlagStateRepo
}

func (m *PromoMockStore) ProjectRepo() repository.ProjectRepository         { return nil }
func (m *PromoMockStore) EnvironmentRepo() repository.EnvironmentRepository { return m.envRepo }
func (m *PromoMockStore) FlagRepo() repository.FlagRepository               { return nil }
func (m *PromoMockStore) FlagStateRepo() repository.FlagStateRepository      { return m.fsRepo }
func (m *PromoMockStore) AuditRepo() repository.AuditRepository             { return nil }
func (m *PromoMockStore) ChangeRequestRepo() repository.ChangeRequestRepository { return m.crRepo }
func (m *PromoMockStore) RoleRepo() repository.RoleRepository { return nil }
func (m *PromoMockStore) UserRepo() repository.UserRepository { return nil }
func (m *PromoMockStore) SlackConfigRepo() repository.SlackConfigRepository { return nil }
func (m *PromoMockStore) ScheduledChangeRepo() repository.ScheduledChangeRepository { return nil }
func (m *PromoMockStore) StalePolicyRepo() repository.StalePolicyRepository { return nil }
func (m *PromoMockStore) KillSwitchRepo() repository.KillSwitchRepository { return nil }

func (m *PromoMockStore) WithTx(ctx context.Context, fn func(repository.Store) error) error {
	return fn(m)
}

func (m *PromoMockStore) MigrateUp() error   { return nil }
func (m *PromoMockStore) MigrateDown() error { return nil }

func TestPromotionService_PromoteToUnprotected(t *testing.T) {
	envRepo := new(MockEnvironmentRepo)
	fsRepo := new(MockFlagStateRepo)
	store := &PromoMockStore{envRepo: envRepo, fsRepo: fsRepo}

	svc := services.NewPromotionService(store, nil, nil)

	sourceEnvID := uuid.New()
	targetEnvID := uuid.New()
	flagID := uuid.New()
	actorID := uuid.New()

	sourceState := &models.EnvironmentFlagState{
		EnvironmentID: sourceEnvID,
		FeatureFlagID: flagID,
		Enabled:       true,
	}

	targetEnv := &models.Environment{
		ID:          targetEnvID,
		IsProtected: false,
	}

	fsRepo.On("GetByEnvAndFlag", mock.Anything, sourceEnvID, flagID).Return(sourceState, nil)
	envRepo.On("GetByID", mock.Anything, targetEnvID).Return(targetEnv, nil)
	fsRepo.On("GetByEnvAndFlag", mock.Anything, targetEnvID, flagID).Return(nil, repository.ErrNotFound)
	fsRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.EnvironmentFlagState")).Return(nil)

	res, err := svc.PromoteFlag(context.Background(), flagID, sourceEnvID, targetEnvID, actorID)
	assert.NoError(t, err)

	state, ok := res.(*models.EnvironmentFlagState)
	assert.True(t, ok)
	assert.True(t, state.Enabled)
	assert.Equal(t, targetEnvID, state.EnvironmentID)
}

func TestPromotionService_PromoteToProtected(t *testing.T) {
	envRepo := new(MockEnvironmentRepo)
	fsRepo := new(MockFlagStateRepo)
	crRepo := new(MockChangeRequestRepo)
	store := &PromoMockStore{envRepo: envRepo, fsRepo: fsRepo, crRepo: crRepo}

	crSvc := services.NewChangeRequestService(store, nil)
	svc := services.NewPromotionService(store, nil, crSvc)

	sourceEnvID := uuid.New()
	targetEnvID := uuid.New()
	flagID := uuid.New()
	actorID := uuid.New()

	sourceState := &models.EnvironmentFlagState{
		EnvironmentID: sourceEnvID,
		FeatureFlagID: flagID,
		Enabled:       true,
	}

	targetEnv := &models.Environment{
		ID:          targetEnvID,
		IsProtected: true,
	}

	fsRepo.On("GetByEnvAndFlag", mock.Anything, sourceEnvID, flagID).Return(sourceState, nil)
	envRepo.On("GetByID", mock.Anything, targetEnvID).Return(targetEnv, nil)
	fsRepo.On("GetByEnvAndFlag", mock.Anything, targetEnvID, flagID).Return(nil, repository.ErrNotFound)
	crRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.ChangeRequest")).Return(nil)

	res, err := svc.PromoteFlag(context.Background(), flagID, sourceEnvID, targetEnvID, actorID)
	assert.NoError(t, err)

	cr, ok := res.(*models.ChangeRequest)
	assert.True(t, ok)
	assert.Equal(t, models.StatusPending, cr.Status)
	assert.Equal(t, targetEnvID, cr.EnvironmentID)
	enabled, _ := cr.ProposedChanges["enabled"].(bool)
	assert.True(t, enabled)
}
