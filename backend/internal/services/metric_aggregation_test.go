package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MetricMockFlagStateRepo struct {
	mock.Mock
}

func (m *MetricMockFlagStateRepo) Create(ctx context.Context, state *models.EnvironmentFlagState) error { return nil }
func (m *MetricMockFlagStateRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.EnvironmentFlagState, error) { return nil, nil }
func (m *MetricMockFlagStateRepo) GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*models.EnvironmentFlagState, error) { return nil, nil }
func (m *MetricMockFlagStateRepo) ListByEnvironment(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentFlagState, error) { return nil, nil }
func (m *MetricMockFlagStateRepo) ListByEnvironmentAndLifecycle(ctx context.Context, envID uuid.UUID, lifecycle models.FlagLifecycleState) ([]*models.EnvironmentFlagState, error) { return nil, nil }
func (m *MetricMockFlagStateRepo) Update(ctx context.Context, state *models.EnvironmentFlagState) error { return nil }
func (m *MetricMockFlagStateRepo) UpdateLifecycleState(ctx context.Context, id uuid.UUID, state models.FlagLifecycleState) error { return nil }
func (m *MetricMockFlagStateRepo) FindActiveFlagsForStalenessScan(ctx context.Context, limit int) ([]*models.EnvironmentFlagState, error) { return nil, nil }
func (m *MetricMockFlagStateRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }

func (m *MetricMockFlagStateRepo) UpdateLastEvaluatedAtBatch(ctx context.Context, updates map[uuid.UUID]time.Time) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func TestMetricAggregationService_RecordAndFlush(t *testing.T) {
	mockRepo := new(MetricMockFlagStateRepo)
	svc := services.NewMetricAggregationService(mockRepo)

	id1 := uuid.New()
	now := time.Now()

	svc.RecordEvaluation(id1, now)

	mockRepo.On("UpdateLastEvaluatedAtBatch", mock.Anything, mock.MatchedBy(func(updates map[uuid.UUID]time.Time) bool {
		val, ok := updates[id1]
		return ok && val.Equal(now)
	})).Return(nil)

	err := svc.Flush(context.Background())
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
