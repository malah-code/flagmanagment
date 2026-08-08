package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/cache"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
)

type WebhookService interface {
	ProcessAPMAlert(ctx context.Context, envID uuid.UUID, alertIdentifier string, payload interface{}) (int, error)
	CreateWebhook(ctx context.Context, wh *models.WebhookIntegration) error
	ListWebhooks(ctx context.Context, projectID uuid.UUID) ([]*models.WebhookIntegration, error)
	DeleteWebhook(ctx context.Context, id uuid.UUID) error
}

type webhookService struct {
	store           repository.Store
	audit           *AuditService
	cacheClient     *cache.Client
	notificationSvc NotificationService
}

func NewWebhookService(store repository.Store, audit *AuditService, cacheClient *cache.Client, notificationSvc NotificationService) WebhookService {
	return &webhookService{
		store:           store,
		audit:           audit,
		cacheClient:     cacheClient,
		notificationSvc: notificationSvc,
	}
}

func (s *webhookService) ProcessAPMAlert(ctx context.Context, envID uuid.UUID, alertIdentifier string, payload interface{}) (int, error) {
	// T010: Query kill_switches
	rules, err := s.store.KillSwitchRepo().ListByEnvironmentAndAlert(ctx, envID, alertIdentifier)
	if err != nil {
		return 0, err
	}

	flagsKilled := 0

	for _, rule := range rules {
		if rule.Action != "DISABLE" {
			continue // Currently we only support DISABLE
		}

		err := s.store.WithTx(ctx, func(txStore repository.Store) error {
			// T011: Disable flag
			state, err := txStore.FlagStateRepo().GetByEnvAndFlag(ctx, envID, rule.FlagID)
			if err != nil {
				if err == repository.ErrNotFound {
					// Flag state doesn't exist yet, create it as disabled
					now := time.Now().UTC()
					state = &models.EnvironmentFlagState{
						ID:             uuid.New(),
						EnvironmentID:  envID,
						FeatureFlagID:  rule.FlagID,
						Enabled:        false,
						CreatedAt:      now,
						UpdatedAt:      now,
					}
					if err := txStore.FlagStateRepo().Create(ctx, state); err != nil {
						return err
					}
					s.logAudit(ctx, txStore, envID, state, nil, state, payload)
					return nil
				}
				return err
			}
			
			if !state.Enabled {
				// Already disabled
				return nil
			}
			
			// Copy old state for audit
			bPrev, _ := json.Marshal(state)
			var prevState models.JSONB
			json.Unmarshal(bPrev, &prevState)

			state.Enabled = false
			state.UpdatedAt = time.Now().UTC()
			if err := txStore.FlagStateRepo().Update(ctx, state); err != nil {
				return err
			}
			s.logAudit(ctx, txStore, envID, state, prevState, state, payload)
			return nil
		})

		if err == nil {
			flagsKilled++
			if s.notificationSvc != nil {
				s.notificationSvc.SendFlagStateChanged(ctx, envID, rule.FlagID.String(), true, false, "Automated APM Kill-Switch")
			}
		}
	}

	if s.cacheClient != nil && flagsKilled > 0 {
		_ = s.cacheClient.PublishRulesetUpdate(ctx, envID.String(), time.Now().UTC().String())
	}

	return flagsKilled, nil
}

func (s *webhookService) logAudit(ctx context.Context, txStore repository.Store, envID uuid.UUID, state *models.EnvironmentFlagState, prevState interface{}, newStateStruct *models.EnvironmentFlagState, apmPayload interface{}) {
	bStateNew, _ := json.Marshal(map[string]interface{}{
		"flag_state":  newStateStruct,
		"apm_payload": apmPayload,
	})
	var newState models.JSONB
	json.Unmarshal(bStateNew, &newState)

	// Since it's an automated action via webhook, there's no user actor.
	// We use the nil actor ID (or a system UUID) for automation.
	var prev models.JSONB
	if p, ok := prevState.(models.JSONB); ok {
		prev = p
	}

	_ = txStore.AuditRepo().Create(ctx, &models.AuditLog{
		ID:            uuid.New(),
		EnvironmentID: &envID,
		ActorID:       uuid.Nil, // System
		Action:        "FLAG_KILLED_AUTOMATICALLY",
		TargetType:    "FLAG_STATE",
		TargetID:      state.ID,
		PreviousState: prev,
		NewState:      newState,
		ActorIP:       "webhook",
		CreatedAt:     time.Now().UTC(),
	})
}

func (s *webhookService) CreateWebhook(ctx context.Context, wh *models.WebhookIntegration) error {
	return s.store.WebhookIntegrationRepo().Create(ctx, wh)
}

func (s *webhookService) ListWebhooks(ctx context.Context, projectID uuid.UUID) ([]*models.WebhookIntegration, error) {
	return s.store.WebhookIntegrationRepo().ListByProject(ctx, projectID)
}

func (s *webhookService) DeleteWebhook(ctx context.Context, id uuid.UUID) error {
	return s.store.WebhookIntegrationRepo().Delete(ctx, id)
}
