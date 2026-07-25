package services

import (
	"context"
	"errors"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrSelfApprovalNotAllowed = errors.New("author cannot approve their own change request")
	ErrChangeRequestNotPending = errors.New("change request is not in PENDING state")
)

type ChangeRequestService struct {
	store repository.Store
	audit *AuditService
}

func NewChangeRequestService(store repository.Store, audit *AuditService) *ChangeRequestService {
	return &ChangeRequestService{
		store: store,
		audit: audit,
	}
}

func (s *ChangeRequestService) Create(ctx context.Context, cr *models.ChangeRequest) error {
	if cr.ID == uuid.Nil {
		cr.ID = uuid.New()
	}
	now := time.Now().UTC()
	cr.CreatedAt = now
	cr.UpdatedAt = now
	cr.Status = models.StatusPending

	if err := s.store.ChangeRequestRepo().Create(ctx, cr); err != nil {
		return err
	}

	if s.audit != nil {
		s.audit.LogAction(ctx, &models.AuditLog{
			ID:            uuid.New(),
			ProjectID:     &cr.ProjectID,
			EnvironmentID: &cr.EnvironmentID,
			ActorID:       cr.CreatedBy,
			Action:        "CREATE",
			TargetType:    "CHANGE_REQUEST",
			TargetID:      cr.ID,
			NewState:      cr.ProposedChanges,
			CreatedAt:     now,
		})
	}

	return nil
}

func (s *ChangeRequestService) GetByID(ctx context.Context, id uuid.UUID) (*models.ChangeRequest, error) {
	return s.store.ChangeRequestRepo().GetByID(ctx, id)
}

func (s *ChangeRequestService) ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*models.ChangeRequest, int, error) {
	return s.store.ChangeRequestRepo().ListByEnvironment(ctx, envID, status, limit, offset)
}

func (s *ChangeRequestService) Approve(ctx context.Context, crID uuid.UUID, reviewerID uuid.UUID, comment string) error {
	cr, err := s.store.ChangeRequestRepo().GetByID(ctx, crID)
	if err != nil {
		return err
	}

	if cr.Status != models.StatusPending {
		return ErrChangeRequestNotPending
	}

	if cr.CreatedBy == reviewerID {
		return ErrSelfApprovalNotAllowed
	}

	now := time.Now().UTC()

	// Perform state application and status update in a transaction
	err = s.store.WithTx(ctx, func(txStore repository.Store) error {
		// 1. Add approval record
		approval := &models.ChangeRequestApproval{
			ID:              uuid.New(),
			ChangeRequestID: crID,
			ApproverID:      reviewerID,
			Decision:        models.DecisionApprove,
			Comment:         comment,
			CreatedAt:       now,
		}
		if err := txStore.ChangeRequestRepo().AddApproval(ctx, approval); err != nil {
			return err
		}

		// 2. Mark Change Request as APPLIED
		if err := txStore.ChangeRequestRepo().UpdateStatus(ctx, crID, string(models.StatusApplied), &reviewerID); err != nil {
			return err
		}

		// 3. Apply proposed changes to target environment flag state if available in proposed_changes
		flagIDRaw, ok1 := cr.ProposedChanges["flag_id"].(string)
		enabled, ok2 := cr.ProposedChanges["enabled"].(bool)
		targetingRules, _ := cr.ProposedChanges["targeting_rules"].(map[string]interface{})
		remoteConfig, _ := cr.ProposedChanges["remote_config"].(map[string]interface{})

		if ok1 && ok2 {
			flagID, err := uuid.Parse(flagIDRaw)
			if err == nil {
				existingState, err := txStore.FlagStateRepo().GetByEnvAndFlag(ctx, cr.EnvironmentID, flagID)
				if err == nil {
					existingState.Enabled = enabled
					if targetingRules != nil {
						existingState.TargetingRules = targetingRules
					}
					if remoteConfig != nil {
						existingState.RemoteConfig = remoteConfig
					}
					existingState.UpdatedAt = now
					if err := txStore.FlagStateRepo().Update(ctx, existingState); err != nil {
						return err
					}
				} else if errors.Is(err, repository.ErrNotFound) {
					newState := &models.EnvironmentFlagState{
						ID:             uuid.New(),
						EnvironmentID:  cr.EnvironmentID,
						FeatureFlagID:  flagID,
						Enabled:        enabled,
						TargetingRules: targetingRules,
						RemoteConfig:   remoteConfig,
						CreatedAt:      now,
						UpdatedAt:      now,
					}
					if err := txStore.FlagStateRepo().Create(ctx, newState); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if s.audit != nil {
		s.audit.LogAction(ctx, &models.AuditLog{
			ID:            uuid.New(),
			ProjectID:     &cr.ProjectID,
			EnvironmentID: &cr.EnvironmentID,
			ActorID:       reviewerID,
			Action:        "APPROVE",
			TargetType:    "CHANGE_REQUEST",
			TargetID:      cr.ID,
			NewState:      cr.ProposedChanges,
			CreatedAt:     now,
		})
	}

	return nil
}

func (s *ChangeRequestService) Reject(ctx context.Context, crID uuid.UUID, reviewerID uuid.UUID, reason string) error {
	cr, err := s.store.ChangeRequestRepo().GetByID(ctx, crID)
	if err != nil {
		return err
	}

	if cr.Status != models.StatusPending {
		return ErrChangeRequestNotPending
	}

	now := time.Now().UTC()

	err = s.store.WithTx(ctx, func(txStore repository.Store) error {
		approval := &models.ChangeRequestApproval{
			ID:              uuid.New(),
			ChangeRequestID: crID,
			ApproverID:      reviewerID,
			Decision:        models.DecisionReject,
			Comment:         reason,
			CreatedAt:       now,
		}
		if err := txStore.ChangeRequestRepo().AddApproval(ctx, approval); err != nil {
			return err
		}

		return txStore.ChangeRequestRepo().UpdateStatus(ctx, crID, string(models.StatusRejected), &reviewerID)
	})

	if err != nil {
		return err
	}

	if s.audit != nil {
		s.audit.LogAction(ctx, &models.AuditLog{
			ID:            uuid.New(),
			ProjectID:     &cr.ProjectID,
			EnvironmentID: &cr.EnvironmentID,
			ActorID:       reviewerID,
			Action:        "REJECT",
			TargetType:    "CHANGE_REQUEST",
			TargetID:      cr.ID,
			CreatedAt:     now,
		})
	}

	return nil
}
