package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrScheduleNotPending = errors.New("scheduled change is not in PENDING state")
)

type ScheduledChangeService struct {
	store repository.Store
	audit *AuditService
}

func NewScheduledChangeService(store repository.Store, audit *AuditService) *ScheduledChangeService {
	return &ScheduledChangeService{
		store: store,
		audit: audit,
	}
}

func (s *ScheduledChangeService) Create(ctx context.Context, sc *models.ScheduledChange) error {
	now := time.Now().UTC()
	if !sc.ScheduledFor.After(now) {
		return errors.New("scheduled_for must be in the future")
	}

	if sc.TargetType == models.TargetTypeFlag {
		if sc.Action != models.ActionEnable && sc.Action != models.ActionDisable {
			return fmt.Errorf("action %s invalid for FLAG target type (must be ENABLE or DISABLE)", sc.Action)
		}
	} else if sc.TargetType == models.TargetTypeChangeRequest {
		if sc.Action != models.ActionApply {
			return fmt.Errorf("action %s invalid for CHANGE_REQUEST target type (must be APPLY)", sc.Action)
		}
	} else {
		return fmt.Errorf("invalid target_type: %s", sc.TargetType)
	}

	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	sc.Status = models.ScheduleStatusPending
	sc.CreatedAt = now
	sc.UpdatedAt = now

	// Check for existing pending schedule (belt-and-suspenders guard; the
	// partial unique index in the DB is the authoritative enforcement).
	existing, err := s.store.ScheduledChangeRepo().GetPendingByTargetID(ctx, sc.TargetID)
	if err == nil && existing != nil {
		return repository.ErrPendingScheduleExists
	} else if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	if err := s.store.ScheduledChangeRepo().Create(ctx, sc); err != nil {
		// The repo maps 23505 unique violations to repository.ErrPendingScheduleExists.
		return err
	}

	if s.audit != nil {
		s.audit.LogAction(ctx, &models.AuditLog{
			ID:            uuid.New(),
			ProjectID:     &sc.ProjectID,
			EnvironmentID: &sc.EnvironmentID,
			ActorID:       sc.CreatedBy,
			Action:        "CREATE",
			TargetType:    "SCHEDULED_CHANGE",
			TargetID:      sc.ID,
			CreatedAt:     now,
		})
	}

	return nil
}

func (s *ScheduledChangeService) GetByID(ctx context.Context, id uuid.UUID) (*models.ScheduledChange, error) {
	return s.store.ScheduledChangeRepo().GetByID(ctx, id)
}

func (s *ScheduledChangeService) ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*models.ScheduledChange, int, error) {
	return s.store.ScheduledChangeRepo().ListByEnvironment(ctx, envID, status, limit, offset)
}

func (s *ScheduledChangeService) Cancel(ctx context.Context, id uuid.UUID, actorID uuid.UUID) (*models.ScheduledChange, error) {
	sc, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if sc.Status != models.ScheduleStatusPending {
		return nil, ErrScheduleNotPending
	}

	now := time.Now().UTC()
	if err := s.store.ScheduledChangeRepo().MarkCancelled(ctx, id, now); err != nil {
		return nil, err
	}

	if s.audit != nil {
		s.audit.LogAction(ctx, &models.AuditLog{
			ID:            uuid.New(),
			ProjectID:     &sc.ProjectID,
			EnvironmentID: &sc.EnvironmentID,
			ActorID:       actorID,
			Action:        "CANCEL",
			TargetType:    "SCHEDULED_CHANGE",
			TargetID:      sc.ID,
			CreatedAt:     now,
		})
	}

	return s.GetByID(ctx, id)
}

func (s *ScheduledChangeService) UpdateScheduledFor(ctx context.Context, id uuid.UUID, newTime time.Time, actorID uuid.UUID) (*models.ScheduledChange, error) {
	now := time.Now().UTC()
	if !newTime.After(now) {
		return nil, errors.New("scheduled_for must be in the future")
	}

	sc, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if sc.Status != models.ScheduleStatusPending {
		return nil, ErrScheduleNotPending
	}

	if err := s.store.ScheduledChangeRepo().UpdateScheduledFor(ctx, id, newTime); err != nil {
		return nil, err
	}

	if s.audit != nil {
		s.audit.LogAction(ctx, &models.AuditLog{
			ID:            uuid.New(),
			ProjectID:     &sc.ProjectID,
			EnvironmentID: &sc.EnvironmentID,
			ActorID:       actorID,
			Action:        "UPDATE",
			TargetType:    "SCHEDULED_CHANGE",
			TargetID:      sc.ID,
			CreatedAt:     now,
		})
	}

	return s.GetByID(ctx, id)
}
