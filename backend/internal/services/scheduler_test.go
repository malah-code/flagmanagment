package services_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type mockFlagStateRepo struct {
	repository.FlagStateRepository
	states map[string]*models.EnvironmentFlagState
}

func (m *mockFlagStateRepo) GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*models.EnvironmentFlagState, error) {
	key := envID.String() + ":" + flagID.String()
	st, ok := m.states[key]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return st, nil
}

func (m *mockFlagStateRepo) Update(ctx context.Context, state *models.EnvironmentFlagState) error {
	key := state.EnvironmentID.String() + ":" + state.FeatureFlagID.String()
	m.states[key] = state
	return nil
}

type schedulerMockStore struct {
	repository.Store
	scRepo        *mockSCRepo
	flagStateRepo *mockFlagStateRepo
	auditRepo     *mockAuditRepo
}

func (m *schedulerMockStore) ScheduledChangeRepo() repository.ScheduledChangeRepository {
	return m.scRepo
}

func (m *schedulerMockStore) FlagStateRepo() repository.FlagStateRepository {
	return m.flagStateRepo
}

func (m *schedulerMockStore) AuditRepo() repository.AuditRepository {
	return m.auditRepo
}

func (m *schedulerMockStore) WithTx(ctx context.Context, fn func(repository.Store) error) error {
	return fn(m)
}

func (m *mockSCRepo) GetDueSchedules(ctx context.Context, now time.Time, limit int) ([]*models.ScheduledChange, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var due []*models.ScheduledChange
	for _, sc := range m.schedules {
		if sc.Status == models.ScheduleStatusPending && !sc.ScheduledFor.After(now) {
			due = append(due, sc)
		}
	}
	return due, nil
}

func (m *mockSCRepo) MarkExecuted(ctx context.Context, id uuid.UUID, executedAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	sc, ok := m.schedules[id]
	if !ok {
		return repository.ErrNotFound
	}
	sc.Status = models.ScheduleStatusExecuted
	sc.ExecutedAt = &executedAt
	return nil
}

func TestScheduler_ExecutesDueSchedules(t *testing.T) {
	scRepo := newMockSCRepo()
	flagStateRepo := &mockFlagStateRepo{states: make(map[string]*models.EnvironmentFlagState)}
	auditRepo := &mockAuditRepo{}
	store := &schedulerMockStore{scRepo: scRepo, flagStateRepo: flagStateRepo, auditRepo: auditRepo}

	envID := uuid.New()
	flagID := uuid.New()
	key := envID.String() + ":" + flagID.String()

	// Initial flag state is OFF
	flagStateRepo.states[key] = &models.EnvironmentFlagState{
		ID:            uuid.New(),
		EnvironmentID: envID,
		FeatureFlagID: flagID,
		Enabled:       false,
	}

	scID := uuid.New()
	scRepo.schedules[scID] = &models.ScheduledChange{
		ID:            scID,
		ProjectID:     uuid.New(),
		EnvironmentID: envID,
		TargetType:    models.TargetTypeFlag,
		TargetID:      flagID,
		Action:        models.ActionEnable,
		ScheduledFor:  time.Now().UTC().Add(-1 * time.Minute), // Due in the past
		Status:        models.ScheduleStatusPending,
	}

	logger := zerolog.New(io.Discard)
	scheduler := services.NewScheduler(store, nil, nil, nil, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start scheduler in goroutine and give it time to process
	go scheduler.Start(ctx)

	time.Sleep(100 * time.Millisecond)
	cancel()

	// Verify flag state turned ON
	if !flagStateRepo.states[key].Enabled {
		t.Fatalf("expected flag state to be enabled (ON), got false")
	}

	// Verify schedule marked EXECUTED
	scRepo.mu.Lock()
	executedSC := scRepo.schedules[scID]
	scRepo.mu.Unlock()

	if executedSC.Status != models.ScheduleStatusExecuted {
		t.Fatalf("expected schedule status EXECUTED, got %s", executedSC.Status)
	}
}
