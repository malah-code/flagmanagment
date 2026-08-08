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

type mockSCRepo struct {
	repository.ScheduledChangeRepository
	mu        sync.Mutex
	schedules map[uuid.UUID]*models.ScheduledChange
}

func newMockSCRepo() *mockSCRepo {
	return &mockSCRepo{
		schedules: make(map[uuid.UUID]*models.ScheduledChange),
	}
}

func (m *mockSCRepo) Create(ctx context.Context, sc *models.ScheduledChange) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.schedules[sc.ID] = sc
	return nil
}

func (m *mockSCRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.ScheduledChange, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sc, ok := m.schedules[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return sc, nil
}

func (m *mockSCRepo) GetPendingByTargetID(ctx context.Context, targetID uuid.UUID) (*models.ScheduledChange, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sc := range m.schedules {
		if sc.TargetID == targetID && sc.Status == models.ScheduleStatusPending {
			return sc, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockSCRepo) MarkCancelled(ctx context.Context, id uuid.UUID, cancelledAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	sc, ok := m.schedules[id]
	if !ok {
		return repository.ErrNotFound
	}
	if sc.Status != models.ScheduleStatusPending {
		return repository.ErrConflict
	}
	sc.Status = models.ScheduleStatusCancelled
	sc.CancelledAt = &cancelledAt
	return nil
}

func (m *mockSCRepo) UpdateScheduledFor(ctx context.Context, id uuid.UUID, newTime time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	sc, ok := m.schedules[id]
	if !ok {
		return repository.ErrNotFound
	}
	if sc.Status != models.ScheduleStatusPending {
		return repository.ErrConflict
	}
	sc.ScheduledFor = newTime
	return nil
}

type fullMockStore struct {
	repository.Store
	scRepo    *mockSCRepo
	auditRepo *mockAuditRepo
}

func (m *fullMockStore) ScheduledChangeRepo() repository.ScheduledChangeRepository {
	return m.scRepo
}

func (m *fullMockStore) AuditRepo() repository.AuditRepository {
	return m.auditRepo
}

func (m *fullMockStore) RoleRepo() repository.RoleRepository { return nil }
func (m *fullMockStore) UserRepo() repository.UserRepository { return nil }
func (m *fullMockStore) WebhookIntegrationRepo() repository.WebhookIntegrationRepository { return nil }

func TestScheduledChangeService_Create(t *testing.T) {
	scRepo := newMockSCRepo()
	auditRepo := &mockAuditRepo{}
	store := &fullMockStore{scRepo: scRepo, auditRepo: auditRepo}
	svc := services.NewScheduledChangeService(store, nil)

	// Valid future schedule
	sc := &models.ScheduledChange{
		ProjectID:     uuid.New(),
		EnvironmentID: uuid.New(),
		TargetType:    models.TargetTypeFlag,
		TargetID:      uuid.New(),
		Action:        models.ActionEnable,
		ScheduledFor:  time.Now().UTC().Add(1 * time.Hour),
		CreatedBy:     uuid.New(),
	}

	err := svc.Create(context.Background(), sc)
	if err != nil {
		t.Fatalf("expected nil error on valid create, got: %v", err)
	}

	// Schedule in the past should fail
	pastSc := &models.ScheduledChange{
		ProjectID:     sc.ProjectID,
		EnvironmentID: sc.EnvironmentID,
		TargetType:    models.TargetTypeFlag,
		TargetID:      uuid.New(),
		Action:        models.ActionEnable,
		ScheduledFor:  time.Now().UTC().Add(-1 * time.Hour),
		CreatedBy:     sc.CreatedBy,
	}
	err = svc.Create(context.Background(), pastSc)
	if err == nil {
		t.Fatalf("expected error for past scheduled_for, got nil")
	}

	// Duplicate pending schedule for same target_id should fail with ErrPendingScheduleExists
	dupSc := &models.ScheduledChange{
		ProjectID:     sc.ProjectID,
		EnvironmentID: sc.EnvironmentID,
		TargetType:    models.TargetTypeFlag,
		TargetID:      sc.TargetID, // same target_id
		Action:        models.ActionDisable,
		ScheduledFor:  time.Now().UTC().Add(2 * time.Hour),
		CreatedBy:     sc.CreatedBy,
	}
	err = svc.Create(context.Background(), dupSc)
	if err != repository.ErrPendingScheduleExists {
		t.Fatalf("expected ErrPendingScheduleExists, got: %v", err)
	}
}

func TestScheduledChangeService_Cancel(t *testing.T) {
	scRepo := newMockSCRepo()
	store := &fullMockStore{scRepo: scRepo}
	svc := services.NewScheduledChangeService(store, nil)

	id := uuid.New()
	scRepo.schedules[id] = &models.ScheduledChange{
		ID:           id,
		TargetID:     uuid.New(),
		Status:       models.ScheduleStatusPending,
		ScheduledFor: time.Now().UTC().Add(1 * time.Hour),
	}

	cancelled, err := svc.Cancel(context.Background(), id, uuid.New())
	if err != nil {
		t.Fatalf("expected successful cancel, got: %v", err)
	}
	if cancelled.Status != models.ScheduleStatusCancelled {
		t.Fatalf("expected status CANCELLED, got %s", cancelled.Status)
	}

	// Cancelling again should fail with ErrScheduleNotPending
	_, err = svc.Cancel(context.Background(), id, uuid.New())
	if err != services.ErrScheduleNotPending {
		t.Fatalf("expected ErrScheduleNotPending, got: %v", err)
	}
}
