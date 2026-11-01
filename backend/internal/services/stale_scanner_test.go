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

type MockStaleStore struct {
	mock.Mock
	flagStateRepo *MockStaleFlagStateRepo
	envRepo       *MockStaleEnvRepo
	flagRepo      *MockStaleFlagRepo
	policyRepo    *MockStalePolicyRepo
	crRepo        *MockStaleCRRepo
}

func (m *MockStaleStore) ProjectRepo() repository.ProjectRepository             { return nil }
func (m *MockStaleStore) EnvironmentRepo() repository.EnvironmentRepository     { return m.envRepo }
func (m *MockStaleStore) FlagRepo() repository.FlagRepository                   { return m.flagRepo }
func (m *MockStaleStore) FlagStateRepo() repository.FlagStateRepository         { return m.flagStateRepo }
func (m *MockStaleStore) AuditRepo() repository.AuditRepository                 { return nil }
func (m *MockStaleStore) ChangeRequestRepo() repository.ChangeRequestRepository { return m.crRepo }
func (m *MockStaleStore) KillSwitchRepo() repository.KillSwitchRepository       { return nil }
func (m *MockStaleStore) SlackConfigRepo() repository.SlackConfigRepository     { return nil }
func (m *MockStaleStore) ScheduledChangeRepo() repository.ScheduledChangeRepository { return nil }
func (m *MockStaleStore) StalePolicyRepo() repository.StalePolicyRepository     { return m.policyRepo }
func (m *MockStaleStore) RoleRepo() repository.RoleRepository                   { return nil }
func (m *MockStaleStore) UserRepo() repository.UserRepository                   { return nil }
func (m *MockStaleStore) ServiceAccountRepo() repository.ServiceAccountRepository { return nil }
func (m *MockStaleStore) WebhookIntegrationRepo() repository.WebhookIntegrationRepository { return nil }
func (m *MockStaleStore) InvitationRepo() repository.InvitationRepository       { return nil }
func (m *MockStaleStore) SystemConfigRepo() repository.SystemConfigRepository   { return nil }
func (m *MockStaleStore) WithTx(ctx context.Context, fn func(repository.Store) error) error {
	return fn(m)
}
func (m *MockStaleStore) MigrateUp() error   { return nil }
func (m *MockStaleStore) MigrateDown() error { return nil }

type MockStaleFlagStateRepo struct {
	mock.Mock
}

func (m *MockStaleFlagStateRepo) FindActiveFlagsForStalenessScan(ctx context.Context, limit int) ([]*models.EnvironmentFlagState, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*models.EnvironmentFlagState), args.Error(1)
}

func (m *MockStaleFlagStateRepo) UpdateLifecycleState(ctx context.Context, id uuid.UUID, state models.FlagLifecycleState) error {
	args := m.Called(ctx, id, state)
	return args.Error(0)
}

func (m *MockStaleFlagStateRepo) Create(ctx context.Context, state *models.EnvironmentFlagState) error { return nil }
func (m *MockStaleFlagStateRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.EnvironmentFlagState, error) { return nil, nil }
func (m *MockStaleFlagStateRepo) GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*models.EnvironmentFlagState, error) { return nil, nil }
func (m *MockStaleFlagStateRepo) ListByEnvironment(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentFlagState, error) { return nil, nil }
func (m *MockStaleFlagStateRepo) ListByEnvironmentAndLifecycle(ctx context.Context, envID uuid.UUID, lifecycle models.FlagLifecycleState) ([]*models.EnvironmentFlagState, error) { return nil, nil }
func (m *MockStaleFlagStateRepo) Update(ctx context.Context, state *models.EnvironmentFlagState) error { return nil }
func (m *MockStaleFlagStateRepo) UpdateLastEvaluatedAtBatch(ctx context.Context, updates map[uuid.UUID]time.Time) error { return nil }
func (m *MockStaleFlagStateRepo) CloneEnvironmentState(ctx context.Context, sourceEnvID, targetEnvID uuid.UUID) error { return nil }
func (m *MockStaleFlagStateRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }

type MockStaleEnvRepo struct {
	mock.Mock
}

func (m *MockStaleEnvRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Environment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Environment), args.Error(1)
}
func (m *MockStaleEnvRepo) Create(ctx context.Context, env *models.Environment) error { return nil }
func (m *MockStaleEnvRepo) GetByAPIKeyHash(ctx context.Context, apiKeyHash string) (*models.Environment, error) { return nil, nil }
func (m *MockStaleEnvRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.Environment, error) { return nil, nil }
func (m *MockStaleEnvRepo) Update(ctx context.Context, env *models.Environment) error { return nil }
func (m *MockStaleEnvRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockStaleEnvRepo) CreateServerKey(ctx context.Context, key *models.EnvironmentServerKey) error { return nil }
func (m *MockStaleEnvRepo) GetServerKeyByHash(ctx context.Context, keyHash string) (*models.EnvironmentServerKey, error) { return nil, nil }
func (m *MockStaleEnvRepo) ListServerKeys(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentServerKey, error) { return nil, nil }
func (m *MockStaleEnvRepo) DeleteServerKey(ctx context.Context, id uuid.UUID) error { return nil }

type MockStaleFlagRepo struct {
	mock.Mock
}

func (m *MockStaleFlagRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.FeatureFlag, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.FeatureFlag), args.Error(1)
}
func (m *MockStaleFlagRepo) Create(ctx context.Context, flag *models.FeatureFlag) error { return nil }
func (m *MockStaleFlagRepo) GetByKey(ctx context.Context, projectID uuid.UUID, key string) (*models.FeatureFlag, error) { return nil, nil }
func (m *MockStaleFlagRepo) ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*models.FeatureFlag, int, error) { return nil, 0, nil }
func (m *MockStaleFlagRepo) Update(ctx context.Context, flag *models.FeatureFlag) error { return nil }
func (m *MockStaleFlagRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockStaleFlagRepo) UpdateLastEvaluatedAt(ctx context.Context, ids []uuid.UUID) error { return nil }
func (m *MockStaleFlagRepo) ListDependencyMap(ctx context.Context, projectID uuid.UUID) (map[uuid.UUID]*uuid.UUID, error) { return nil, nil }

type MockStalePolicyRepo struct {
	mock.Mock
}

func (m *MockStalePolicyRepo) GetByEnvironment(ctx context.Context, projectID, envID uuid.UUID) (*models.StaleFlagPolicy, error) {
	args := m.Called(ctx, projectID, envID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StaleFlagPolicy), args.Error(1)
}
func (m *MockStalePolicyRepo) Upsert(ctx context.Context, policy *models.StaleFlagPolicy) error { return nil }
func (m *MockStalePolicyRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }

type MockStaleCRRepo struct {
	mock.Mock
}

func (m *MockStaleCRRepo) ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*models.ChangeRequest, int, error) {
	args := m.Called(ctx, envID, status, limit, offset)
	return args.Get(0).([]*models.ChangeRequest), args.Int(1), args.Error(2)
}

func (m *MockStaleCRRepo) Create(ctx context.Context, cr *models.ChangeRequest) error {
	args := m.Called(ctx, cr)
	return args.Error(0)
}
func (m *MockStaleCRRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.ChangeRequest, error) { return nil, nil }
func (m *MockStaleCRRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, appliedBy *uuid.UUID) error { return nil }
func (m *MockStaleCRRepo) AddApproval(ctx context.Context, approval *models.ChangeRequestApproval) error { return nil }
func (m *MockStaleCRRepo) ListApprovals(ctx context.Context, crID uuid.UUID) ([]*models.ChangeRequestApproval, error) { return nil, nil }

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) SendFlagStateChanged(ctx context.Context, envID uuid.UUID, flagKey string, oldState, newState bool, actor string) {
}

func (m *MockNotifier) SendStaleFlagAlert(ctx context.Context, envID uuid.UUID, flagKey string, staleDays int, changeRequestID *uuid.UUID) {
	m.Called(ctx, envID, flagKey, staleDays, changeRequestID)
}

func TestStaleScanner_AutomatedGovernance(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	envID := uuid.New()
	flagID := uuid.New()
	flagStateID := uuid.New()

	staleTime := time.Now().UTC().Add(-45 * 24 * time.Hour) // 45 days ago

	flagState := &models.EnvironmentFlagState{
		ID:                flagStateID,
		FeatureFlagID:     flagID,
		EnvironmentID:     envID,
		Enabled:           true,
		LifecycleState:    models.LifecycleActive,
		LastStateChangeAt: staleTime,
		LastEvaluatedAt:   &staleTime,
	}

	env := &models.Environment{
		ID:        envID,
		ProjectID: projectID,
		Name:      "Production",
	}

	flag := &models.FeatureFlag{
		ID:        flagID,
		ProjectID: projectID,
		Key:       "legacy-payment-flow",
		Name:      "Legacy Payment Flow",
	}

	policy := &models.StaleFlagPolicy{
		ProjectID:      projectID,
		StaleAfterDays: 30,
	}

	mockFlagStateRepo := new(MockStaleFlagStateRepo)
	mockEnvRepo := new(MockStaleEnvRepo)
	mockFlagRepo := new(MockStaleFlagRepo)
	mockPolicyRepo := new(MockStalePolicyRepo)
	mockCRRepo := new(MockStaleCRRepo)
	mockNotifier := new(MockNotifier)

	store := &MockStaleStore{
		flagStateRepo: mockFlagStateRepo,
		envRepo:       mockEnvRepo,
		flagRepo:      mockFlagRepo,
		policyRepo:    mockPolicyRepo,
		crRepo:        mockCRRepo,
	}

	mockFlagStateRepo.On("FindActiveFlagsForStalenessScan", ctx, 1000).Return([]*models.EnvironmentFlagState{flagState}, nil)
	mockEnvRepo.On("GetByID", ctx, envID).Return(env, nil)
	mockPolicyRepo.On("GetByEnvironment", ctx, projectID, envID).Return(policy, nil)
	mockFlagStateRepo.On("UpdateLifecycleState", ctx, flagStateID, models.LifecycleStale).Return(nil)
	mockFlagRepo.On("GetByID", ctx, flagID).Return(flag, nil)

	// No existing CRs
	mockCRRepo.On("ListByEnvironment", ctx, envID, string(models.StatusPending), 100, 0).Return([]*models.ChangeRequest{}, 0, nil)

	// Expect creation of ChangeRequest
	mockCRRepo.On("Create", ctx, mock.MatchedBy(func(cr *models.ChangeRequest) bool {
		return cr.ProjectID == projectID &&
			cr.EnvironmentID == envID &&
			cr.Status == models.StatusPending &&
			cr.ProposedChanges["is_enabled"] == false &&
			cr.ProposedChanges["flag_key"] == "legacy-payment-flow"
	})).Return(nil)

	// Expect Slack notification
	mockNotifier.On("SendStaleFlagAlert", ctx, envID, "legacy-payment-flow", 30, mock.AnythingOfType("*uuid.UUID")).Return()

	scanner := NewStaleScannerService(store, mockNotifier)
	result, err := scanner.ScanAndGovern(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.FlagsScanned)
	assert.Equal(t, 1, result.FlagsMarkedStale)
	assert.Equal(t, 1, result.ChangeRequestsCreated)
	assert.Equal(t, 1, result.NotificationsSent)

	mockFlagStateRepo.AssertExpectations(t)
	mockCRRepo.AssertExpectations(t)
	mockNotifier.AssertExpectations(t)
}

func TestStaleScanner_IdempotentDuplicatePrevention(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	envID := uuid.New()
	flagID := uuid.New()
	flagStateID := uuid.New()

	staleTime := time.Now().UTC().Add(-40 * 24 * time.Hour)

	flagState := &models.EnvironmentFlagState{
		ID:                flagStateID,
		FeatureFlagID:     flagID,
		EnvironmentID:     envID,
		Enabled:           true,
		LifecycleState:    models.LifecycleActive,
		LastStateChangeAt: staleTime,
		LastEvaluatedAt:   &staleTime,
	}

	env := &models.Environment{
		ID:        envID,
		ProjectID: projectID,
		Name:      "Staging",
	}

	flag := &models.FeatureFlag{
		ID:        flagID,
		ProjectID: projectID,
		Key:       "already-deprecated-flag",
	}

	existingCR := &models.ChangeRequest{
		ID:            uuid.New(),
		Title:         "[Automated Governance] Deprecate Stale Flag: already-deprecated-flag",
		Status:        models.StatusPending,
		EnvironmentID: envID,
	}

	mockFlagStateRepo := new(MockStaleFlagStateRepo)
	mockEnvRepo := new(MockStaleEnvRepo)
	mockFlagRepo := new(MockStaleFlagRepo)
	mockPolicyRepo := new(MockStalePolicyRepo)
	mockCRRepo := new(MockStaleCRRepo)

	store := &MockStaleStore{
		flagStateRepo: mockFlagStateRepo,
		envRepo:       mockEnvRepo,
		flagRepo:      mockFlagRepo,
		policyRepo:    mockPolicyRepo,
		crRepo:        mockCRRepo,
	}

	mockFlagStateRepo.On("FindActiveFlagsForStalenessScan", ctx, 1000).Return([]*models.EnvironmentFlagState{flagState}, nil)
	mockEnvRepo.On("GetByID", ctx, envID).Return(env, nil)
	mockPolicyRepo.On("GetByEnvironment", ctx, projectID, envID).Return((*models.StaleFlagPolicy)(nil), repository.ErrNotFound)
	mockFlagStateRepo.On("UpdateLifecycleState", ctx, flagStateID, models.LifecycleStale).Return(nil)
	mockFlagRepo.On("GetByID", ctx, flagID).Return(flag, nil)

	// An existing pending CR is found
	mockCRRepo.On("ListByEnvironment", ctx, envID, string(models.StatusPending), 100, 0).Return([]*models.ChangeRequest{existingCR}, 1, nil)

	// mockCRRepo.Create should NOT be called
	scanner := NewStaleScannerService(store)
	result, err := scanner.ScanAndGovern(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.FlagsMarkedStale)
	assert.Equal(t, 0, result.ChangeRequestsCreated, "Should not create duplicate Change Request")

	mockCRRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}
