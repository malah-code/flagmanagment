package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/google/uuid"
)

type StaleScanResult struct {
	FlagsScanned          int `json:"flags_scanned"`
	FlagsMarkedStale      int `json:"flags_marked_stale"`
	ChangeRequestsCreated int `json:"change_requests_created"`
	NotificationsSent     int `json:"notifications_sent"`
}

type StaleScannerService interface {
	ScanStaleFlags(ctx context.Context) (int, error)
	ScanAndGovern(ctx context.Context) (*StaleScanResult, error)
	Start(ctx context.Context, interval time.Duration)
}

type staleScannerService struct {
	store    repository.Store
	notifier NotificationService
}

func NewStaleScannerService(store repository.Store, notifier ...NotificationService) StaleScannerService {
	var n NotificationService
	if len(notifier) > 0 {
		n = notifier[0]
	}
	return &staleScannerService{store: store, notifier: n}
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
				res, err := s.ScanAndGovern(ctx)
				if err != nil {
					log.Printf("[StaleScanner] error during scan: %v", err)
				} else if res.FlagsMarkedStale > 0 || res.ChangeRequestsCreated > 0 {
					log.Printf("[StaleScanner] marked %d stale, drafted %d CRs, sent %d notifications",
						res.FlagsMarkedStale, res.ChangeRequestsCreated, res.NotificationsSent)
				}
			}
		}
	}()
}

func (s *staleScannerService) ScanStaleFlags(ctx context.Context) (int, error) {
	res, err := s.ScanAndGovern(ctx)
	if err != nil {
		return 0, err
	}
	return res.FlagsMarkedStale, nil
}

func (s *staleScannerService) ScanAndGovern(ctx context.Context) (*StaleScanResult, error) {
	activeFlags, err := s.store.FlagStateRepo().FindActiveFlagsForStalenessScan(ctx, 1000)
	if err != nil {
		return nil, err
	}

	result := &StaleScanResult{
		FlagsScanned: len(activeFlags),
	}

	now := time.Now().UTC()
	botUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001") // System Governance Bot

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

		if !isStale {
			continue
		}

		// Mark state as STALE
		if err := s.store.FlagStateRepo().UpdateLifecycleState(ctx, flagState.ID, models.LifecycleStale); err == nil {
			result.FlagsMarkedStale++
		}

		// Retrieve flag details
		flag, err := s.store.FlagRepo().GetByID(ctx, flagState.FeatureFlagID)
		if err != nil {
			continue
		}

		// Check if an automated or pending deprecation Change Request already exists for this flag
		existingCRs, _, err := s.store.ChangeRequestRepo().ListByEnvironment(ctx, env.ID, string(models.StatusPending), 100, 0)
		hasPendingCR := false
		if err == nil {
			for _, cr := range existingCRs {
				if strings.Contains(cr.Title, flag.Key) {
					hasPendingCR = true
					break
				}
			}
		}

		var crID *uuid.UUID
		if !hasPendingCR {
			newCRID := uuid.New()
			cr := &models.ChangeRequest{
				ID:            newCRID,
				ProjectID:     env.ProjectID,
				EnvironmentID: env.ID,
				Title:         fmt.Sprintf("[Automated Governance] Deprecate Stale Flag: %s", flag.Key),
				Description:   fmt.Sprintf("Flag '%s' has been inactive for %d days exceeding the staleness policy threshold. Automated governance bot drafted this proposal to deprecate and safely disable the flag.", flag.Key, staleDays),
				Status:        models.StatusPending,
				ProposedChanges: models.JSONB{
					"flag_id":         flagState.FeatureFlagID.String(),
					"flag_key":        flag.Key,
					"is_enabled":      false,
					"lifecycle_state": string(models.LifecycleDeprecated),
					"governance_bot":  true,
				},
				CurrentState: models.JSONB{
					"flag_id":         flagState.FeatureFlagID.String(),
					"flag_key":        flag.Key,
					"is_enabled":      flagState.Enabled,
					"lifecycle_state": string(models.LifecycleStale),
				},
				CreatedBy: botUserID,
				CreatedAt: now,
				UpdatedAt: now,
			}

			if err := s.store.ChangeRequestRepo().Create(ctx, cr); err == nil {
				result.ChangeRequestsCreated++
				crID = &newCRID
			}
		}

		// Send Slack Notification if configured
		if s.notifier != nil {
			s.notifier.SendStaleFlagAlert(ctx, env.ID, flag.Key, staleDays, crID)
			result.NotificationsSent++
		}
	}

	return result, nil
}
