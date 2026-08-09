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

// MockStore & MockRepos for testing ChangeRequestService
type MockStore struct {
	mock.Mock
	crRepo *MockChangeRequestRepo
}

func (m *MockStore) ProjectRepo() repository.ProjectRepository         { return nil }
func (m *MockStore) EnvironmentRepo() repository.EnvironmentRepository { return nil }
func (m *MockStore) FlagRepo() repository.FlagRepository               { return nil }
func (m *MockStore) FlagStateRepo() repository.FlagStateRepository      { return nil }
func (m *MockStore) AuditRepo() repository.AuditRepository             { return nil }
func (m *MockStore) ChangeRequestRepo() repository.ChangeRequestRepository {
	return m.crRepo
}
func (m *MockStore) RoleRepo() repository.RoleRepository { return nil }
func (m *MockStore) UserRepo() repository.UserRepository { return nil }
func (m *MockStore) ServiceAccountRepo() repository.ServiceAccountRepository { return nil }
func (m *MockStore) WebhookIntegrationRepo() repository.WebhookIntegrationRepository { return nil }
func (m *MockStore) KillSwitchRepo() repository.KillSwitchRepository { return nil }
func (m *MockStore) SlackConfigRepo() repository.SlackConfigRepository { return nil }
func (m *MockStore) ScheduledChangeRepo() repository.ScheduledChangeRepository { return nil }
func (m *MockStore) StalePolicyRepo() repository.StalePolicyRepository { return nil }

func (m *MockStore) WithTx(ctx context.Context, fn func(repository.Store) error) error {
	return fn(m)
}

func (m *MockStore) MigrateUp() error   { return nil }
func (m *MockStore) MigrateDown() error { return nil }

type MockChangeRequestRepo struct {
	mock.Mock
}

func (m *MockChangeRequestRepo) Create(ctx context.Context, cr *models.ChangeRequest) error {
	args := m.Called(ctx, cr)
	return args.Error(0)
}

func (m *MockChangeRequestRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.ChangeRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChangeRequest), args.Error(1)
}

func (m *MockChangeRequestRepo) ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*models.ChangeRequest, int, error) {
	args := m.Called(ctx, envID, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*models.ChangeRequest), args.Int(1), args.Error(2)
}

func (m *MockChangeRequestRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, appliedBy *uuid.UUID) error {
	args := m.Called(ctx, id, status, appliedBy)
	return args.Error(0)
}

func (m *MockChangeRequestRepo) AddApproval(ctx context.Context, approval *models.ChangeRequestApproval) error {
	args := m.Called(ctx, approval)
	return args.Error(0)
}

func (m *MockChangeRequestRepo) ListApprovals(ctx context.Context, crID uuid.UUID) ([]*models.ChangeRequestApproval, error) {
	args := m.Called(ctx, crID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ChangeRequestApproval), args.Error(1)
}

func TestChangeRequestService_SelfApprovalPrevention(t *testing.T) {
	mockCRRepo := new(MockChangeRequestRepo)
	mockStore := &MockStore{crRepo: mockCRRepo}
	svc := services.NewChangeRequestService(mockStore, nil)

	authorID := uuid.New()
	crID := uuid.New()

	cr := &models.ChangeRequest{
		ID:            crID,
		ProjectID:     uuid.New(),
		EnvironmentID: uuid.New(),
		Status:        models.StatusPending,
		CreatedBy:     authorID,
	}

	mockCRRepo.On("GetByID", mock.Anything, crID).Return(cr, nil)

	err := svc.Approve(context.Background(), crID, authorID, "Approving my own request")

	assert.ErrorIs(t, err, services.ErrSelfApprovalNotAllowed)
	mockCRRepo.AssertExpectations(t)
}

func TestChangeRequestService_Reject(t *testing.T) {
	mockCRRepo := new(MockChangeRequestRepo)
	mockStore := &MockStore{crRepo: mockCRRepo}
	svc := services.NewChangeRequestService(mockStore, nil)

	authorID := uuid.New()
	reviewerID := uuid.New()
	crID := uuid.New()

	cr := &models.ChangeRequest{
		ID:            crID,
		ProjectID:     uuid.New(),
		EnvironmentID: uuid.New(),
		Status:        models.StatusPending,
		CreatedBy:     authorID,
	}

	mockCRRepo.On("GetByID", mock.Anything, crID).Return(cr, nil)
	mockCRRepo.On("AddApproval", mock.Anything, mock.MatchedBy(func(a *models.ChangeRequestApproval) bool {
		return a.ChangeRequestID == crID && a.Decision == models.DecisionReject
	})).Return(nil)
	mockCRRepo.On("UpdateStatus", mock.Anything, crID, string(models.StatusRejected), &reviewerID).Return(nil)

	err := svc.Reject(context.Background(), crID, reviewerID, "Too risky for prod")

	assert.NoError(t, err)
	mockCRRepo.AssertExpectations(t)
}
