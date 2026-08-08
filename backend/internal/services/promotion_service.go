package services

import (
	"context"
	"errors"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/google/uuid"
)

type PromotionService struct {
	store     repository.Store
	audit     *AuditService
	crService *ChangeRequestService
}

func NewPromotionService(store repository.Store, audit *AuditService, crService *ChangeRequestService) *PromotionService {
	return &PromotionService{
		store:     store,
		audit:     audit,
		crService: crService,
	}
}

func (s *PromotionService) PromoteFlag(ctx context.Context, flagID, sourceEnvID, targetEnvID, actorID uuid.UUID) (interface{}, error) {
	sourceState, err := s.store.FlagStateRepo().GetByEnvAndFlag(ctx, sourceEnvID, flagID)
	if err != nil {
		return nil, err
	}

	targetEnv, err := s.store.EnvironmentRepo().GetByID(ctx, targetEnvID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	if targetEnv.IsProtected && s.crService != nil {
		proposed := models.JSONB{
			"flag_id":         flagID.String(),
			"enabled":         sourceState.Enabled,
			"targeting_rules": sourceState.TargetingRules,
			"remote_config":   sourceState.RemoteConfig,
		}

		var currentState models.JSONB
		state, err := s.store.FlagStateRepo().GetByEnvAndFlag(ctx, targetEnvID, flagID)
		if err == nil && state != nil {
			currentState = models.JSONB{
				"flag_id":         state.FeatureFlagID.String(),
				"enabled":         state.Enabled,
				"targeting_rules": state.TargetingRules,
				"remote_config":   state.RemoteConfig,
			}
		} else {
			currentState = models.JSONB{
				"flag_id":         flagID.String(),
				"enabled":         false,
				"targeting_rules": nil,
				"remote_config":   nil,
			}
		}

		cr := &models.ChangeRequest{
			ID:              uuid.New(),
			ProjectID:       targetEnv.ProjectID,
			EnvironmentID:   targetEnvID,
			Title:           "Promote Flag State",
			Description:     "Promotion of flag state from source environment",
			Status:          models.StatusPending,
			ProposedChanges: proposed,
			CurrentState:    currentState,
			CreatedBy:       actorID,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		if err := s.crService.Create(ctx, cr); err != nil {
			return nil, err
		}

		return cr, nil
	}

	// Unprotected target environment: Update directly
	var targetState *models.EnvironmentFlagState
	isUpdate := true
	targetState, err = s.store.FlagStateRepo().GetByEnvAndFlag(ctx, targetEnvID, flagID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			isUpdate = false
			targetState = &models.EnvironmentFlagState{
				ID:             uuid.New(),
				EnvironmentID:  targetEnvID,
				FeatureFlagID:  flagID,
				CreatedAt:      now,
			}
		} else {
			return nil, err
		}
	}

	targetState.Enabled = sourceState.Enabled
	targetState.TargetingRules = sourceState.TargetingRules
	targetState.RemoteConfig = sourceState.RemoteConfig
	targetState.UpdatedAt = now

	if isUpdate {
		if err := s.store.FlagStateRepo().Update(ctx, targetState); err != nil {
			return nil, err
		}
	} else {
		if err := s.store.FlagStateRepo().Create(ctx, targetState); err != nil {
			return nil, err
		}
	}

	if s.audit != nil {
		action := "UPDATE"
		if !isUpdate {
			action = "CREATE"
		}
		
		s.audit.LogAction(ctx, &models.AuditLog{
			ID:            uuid.New(),
			EnvironmentID: &targetEnvID,
			ActorID:       actorID,
			Action:        action,
			TargetType:    "FLAG_STATE_PROMOTION",
			TargetID:      targetState.ID,
			CreatedAt:     now,
		})
	}

	return targetState, nil
}
