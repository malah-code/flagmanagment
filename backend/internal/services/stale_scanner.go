package services

import (
	"context"
	"log"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
)

type StaleScannerService interface {
	ScanStaleFlags(ctx context.Context) (int, error)
	Start(ctx context.Context, interval time.Duration)
}

type staleScannerService struct {
	store repository.Store
}

func NewStaleScannerService(store repository.Store) StaleScannerService {
	return &staleScannerService{store: store}
}

func (s *staleScannerService) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				count, err := s.ScanStaleFlags(ctx)
				if err != nil {
					log.Printf("[StaleScanner] error during scan: %v", err)
				} else if count > 0 {
					log.Printf("[StaleScanner] marked %d flags as STALE", count)
				}
			}
		}
	}()
}

func (s *staleScannerService) ScanStaleFlags(ctx context.Context) (int, error) {
	activeFlags, err := s.store.FlagStateRepo().FindActiveFlagsForStalenessScan(ctx, 1000)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	markedCount := 0

	for _, flagState := range activeFlags {
		// Fetch environment to get projectID
		env, err := s.store.EnvironmentRepo().GetByID(ctx, flagState.EnvironmentID)
		if err != nil {
			continue
		}

		// Fetch policy for environment/project
		staleDays := 30 // default
		if policy, err := s.store.StalePolicyRepo().GetByEnvironment(ctx, env.ProjectID, env.ID); err == nil && policy != nil {
			staleDays = policy.StaleAfterDays
		}

		thresholdDuration := time.Duration(staleDays) * 24 * time.Hour
		isStale := false

		// Check 1: Inactivity (no evaluation in staleDays)
		if flagState.LastEvaluatedAt != nil {
			if now.Sub(*flagState.LastEvaluatedAt) > thresholdDuration {
				isStale = true
			}
		} else {
			// Never evaluated, check creation / last state change
			if now.Sub(flagState.LastStateChangeAt) > thresholdDuration {
				isStale = true
			}
		}

		// Check 2: State unchanged for staleDays
		if now.Sub(flagState.LastStateChangeAt) > thresholdDuration {
			isStale = true
		}

		if isStale {
			if err := s.store.FlagStateRepo().UpdateLifecycleState(ctx, flagState.ID, models.LifecycleStale); err == nil {
				markedCount++
			}
		}
	}

	return markedCount, nil
}
