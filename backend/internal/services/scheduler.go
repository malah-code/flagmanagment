package services

import (
	"context"
	"sync"
	"time"

	"github.com/flagmanagment/backend/internal/cache"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Scheduler struct {
	store       repository.Store
	audit       *AuditService
	crService   *ChangeRequestService
	cacheClient *cache.Client
	logger      zerolog.Logger
	interval    time.Duration
	workerCount int
	stopCh      chan struct{}
	stopOnce    sync.Once
}

func NewScheduler(store repository.Store, audit *AuditService, crService *ChangeRequestService, cacheClient *cache.Client, logger zerolog.Logger) *Scheduler {
	return &Scheduler{
		store:       store,
		audit:       audit,
		crService:   crService,
		cacheClient: cacheClient,
		logger:      logger,
		interval:    30 * time.Second,
		workerCount: 20,
		stopCh:      make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.logger.Info().Msg("scheduler started")
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Initial startup catchup check
	s.processDueSchedules(ctx)

	for {
		select {
		case <-ticker.C:
			s.processDueSchedules(ctx)
		case <-s.stopCh:
			s.logger.Info().Msg("scheduler stopped via stop channel")
			return
		case <-ctx.Done():
			s.logger.Info().Msg("scheduler stopped via context cancellation")
			return
		}
	}
}

// Stop safely stops the scheduler. It is safe to call multiple times.
func (s *Scheduler) Stop() {
	s.stopOnce.Do(func() { close(s.stopCh) })
}

func (s *Scheduler) processDueSchedules(ctx context.Context) {
	now := time.Now().UTC()
	// Fetch up to 500 due schedules per tick to handle burst catch-up after downtime.
	schedules, err := s.store.ScheduledChangeRepo().GetDueSchedules(ctx, now, 500)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to fetch due scheduled changes")
		return
	}

	if len(schedules) == 0 {
		return
	}

	s.logger.Info().Int("count", len(schedules)).Msg("processing due scheduled changes")

	jobs := make(chan *models.ScheduledChange, len(schedules))
	var wg sync.WaitGroup

	numWorkers := s.workerCount
	if len(schedules) < numWorkers {
		numWorkers = len(schedules)
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for sc := range jobs {
				s.executeOne(ctx, sc)
			}
		}()
	}

	for _, sc := range schedules {
		jobs <- sc
	}
	close(jobs)

	wg.Wait()
}

func (s *Scheduler) executeOne(ctx context.Context, sc *models.ScheduledChange) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error().Interface("panic", r).Str("id", sc.ID.String()).Msg("recovered panic during scheduled change execution")
		}
	}()

	s.logger.Info().Str("id", sc.ID.String()).Str("target_type", string(sc.TargetType)).Msg("executing scheduled change")

	var execErr error

	switch sc.TargetType {
	case models.TargetTypeFlag:
		// Wrap flag state update and MarkExecuted in a single transaction so
		// a partial failure can't leave the record in PENDING while the flag
		// state has already been changed (which would cause double execution).
		execErr = s.store.WithTx(ctx, func(txStore repository.Store) error {
			state, err := txStore.FlagStateRepo().GetByEnvAndFlag(ctx, sc.EnvironmentID, sc.TargetID)
			if err != nil {
				return err
			}
			state.Enabled = (sc.Action == models.ActionEnable)
			state.UpdatedAt = time.Now().UTC()
			if err := txStore.FlagStateRepo().Update(ctx, state); err != nil {
				return err
			}
			return txStore.ScheduledChangeRepo().MarkExecuted(ctx, sc.ID, state.UpdatedAt)
		})

	case models.TargetTypeChangeRequest:
		if s.crService == nil {
			s.logger.Warn().Str("id", sc.ID.String()).Msg("CHANGE_REQUEST scheduling skipped: crService not provided")
			return
		}
		execErr = s.crService.ApplyScheduled(ctx, sc.TargetID)
		if execErr == nil {
			// ApplyScheduled handles its own transaction; mark as executed separately.
			execErr = s.store.ScheduledChangeRepo().MarkExecuted(ctx, sc.ID, time.Now().UTC())
		}

	default:
		s.logger.Warn().Str("id", sc.ID.String()).Str("target_type", string(sc.TargetType)).Msg("unknown target type")
		return
	}

	if execErr != nil {
		s.logger.Error().Err(execErr).Str("id", sc.ID.String()).Str("target_type", string(sc.TargetType)).Msg("scheduled change execution failed")
		return
	}

	now := time.Now().UTC()

	if s.audit != nil {
		s.audit.LogAction(ctx, &models.AuditLog{
			ID:            uuid.New(), // fresh ID — sc.ID is the scheduled change, not the log entry
			ProjectID:     &sc.ProjectID,
			EnvironmentID: &sc.EnvironmentID,
			ActorID:       models.SystemActor(),
			Action:        "SCHEDULED_EXECUTION",
			TargetType:    string(sc.TargetType),
			TargetID:      sc.TargetID,
			CreatedAt:     now,
		})
	}

	if s.cacheClient != nil {
		_ = s.cacheClient.PublishRulesetUpdate(ctx, sc.EnvironmentID.String(), now.String())
	}

	s.logger.Info().Str("id", sc.ID.String()).Msg("scheduled change executed successfully")
}
