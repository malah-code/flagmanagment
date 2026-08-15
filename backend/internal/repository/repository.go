package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
)

var (
	ErrNotFound             = errors.New("resource not found")
	ErrAlreadyExists        = errors.New("resource already exists")
	ErrConflict             = errors.New("resource conflict")
	ErrPendingScheduleExists = errors.New("a pending schedule already exists for this flag")
)

// Store provides access to all repository operations.
type Store interface {
	ProjectRepo() ProjectRepository
	EnvironmentRepo() EnvironmentRepository
	FlagRepo() FlagRepository
	FlagStateRepo() FlagStateRepository
	AuditRepo() AuditRepository
	ChangeRequestRepo() ChangeRequestRepository
	KillSwitchRepo() KillSwitchRepository
	SlackConfigRepo() SlackConfigRepository
	ScheduledChangeRepo() ScheduledChangeRepository
	StalePolicyRepo() StalePolicyRepository
	RoleRepo() RoleRepository
	UserRepo() UserRepository
	InvitationRepo() InvitationRepository
	SystemConfigRepo() SystemConfigRepository
	ServiceAccountRepo() ServiceAccountRepository
	WebhookIntegrationRepo() WebhookIntegrationRepository

	// Transaction support
	WithTx(ctx context.Context, fn func(Store) error) error

	// Migration support
	MigrateUp() error
	MigrateDown() error
}

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	GetByKey(ctx context.Context, key string) (*models.Project, error)
	List(ctx context.Context, limit, offset int) ([]*models.Project, int, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type EnvironmentRepository interface {
	Create(ctx context.Context, env *models.Environment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Environment, error)
	GetByAPIKeyHash(ctx context.Context, apiKeyHash string) (*models.Environment, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.Environment, error)
	Update(ctx context.Context, env *models.Environment) error
	Delete(ctx context.Context, id uuid.UUID) error

	CreateServerKey(ctx context.Context, key *models.EnvironmentServerKey) error
	GetServerKeyByHash(ctx context.Context, keyHash string) (*models.EnvironmentServerKey, error)
	ListServerKeys(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentServerKey, error)
	DeleteServerKey(ctx context.Context, id uuid.UUID) error
}

type FlagRepository interface {
	Create(ctx context.Context, flag *models.FeatureFlag) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.FeatureFlag, error)
	GetByKey(ctx context.Context, projectID uuid.UUID, key string) (*models.FeatureFlag, error)
	ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*models.FeatureFlag, int, error)
	Update(ctx context.Context, flag *models.FeatureFlag) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastEvaluatedAt(ctx context.Context, ids []uuid.UUID) error
	ListDependencyMap(ctx context.Context, projectID uuid.UUID) (map[uuid.UUID]*uuid.UUID, error)
}

type FlagStateRepository interface {
	Create(ctx context.Context, state *models.EnvironmentFlagState) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.EnvironmentFlagState, error)
	GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*models.EnvironmentFlagState, error)
	ListByEnvironment(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentFlagState, error)
	ListByEnvironmentAndLifecycle(ctx context.Context, envID uuid.UUID, lifecycle models.FlagLifecycleState) ([]*models.EnvironmentFlagState, error)
	Update(ctx context.Context, state *models.EnvironmentFlagState) error
	UpdateLifecycleState(ctx context.Context, id uuid.UUID, state models.FlagLifecycleState) error
	UpdateLastEvaluatedAtBatch(ctx context.Context, updates map[uuid.UUID]time.Time) error
	FindActiveFlagsForStalenessScan(ctx context.Context, limit int) ([]*models.EnvironmentFlagState, error)
	CloneEnvironmentState(ctx context.Context, sourceEnvID, targetEnvID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type StalePolicyRepository interface {
	GetByEnvironment(ctx context.Context, projectID, envID uuid.UUID) (*models.StaleFlagPolicy, error)
	Upsert(ctx context.Context, policy *models.StaleFlagPolicy) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AuditRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error)
	ListByEnvironment(ctx context.Context, envID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error)
	DeleteOlderThan(ctx context.Context, cutoff time.Time) error
}

type ChangeRequestRepository interface {
	Create(ctx context.Context, cr *models.ChangeRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ChangeRequest, error)
	ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*models.ChangeRequest, int, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, appliedBy *uuid.UUID) error
	AddApproval(ctx context.Context, approval *models.ChangeRequestApproval) error
	ListApprovals(ctx context.Context, crID uuid.UUID) ([]*models.ChangeRequestApproval, error)
}

type ScheduledChangeRepository interface {
	Create(ctx context.Context, sc *models.ScheduledChange) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ScheduledChange, error)
	ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*models.ScheduledChange, int, error)
	GetPendingByTargetID(ctx context.Context, targetID uuid.UUID) (*models.ScheduledChange, error)
	GetDueSchedules(ctx context.Context, now time.Time, limit int) ([]*models.ScheduledChange, error)
	MarkExecuted(ctx context.Context, id uuid.UUID, executedAt time.Time) error
	MarkCancelled(ctx context.Context, id uuid.UUID, cancelledAt time.Time) error
	UpdateScheduledFor(ctx context.Context, id uuid.UUID, newTime time.Time) error
}

type RoleRepository interface {
	Create(ctx context.Context, role *models.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
	List(ctx context.Context) ([]*models.Role, error)
	AssignUserRole(ctx context.Context, ur *models.UserRole) error
	RemoveUserRole(ctx context.Context, id uuid.UUID) error
	RemoveAllUserRoles(ctx context.Context, userID uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID, projectID *uuid.UUID) ([]*models.UserRole, error)
	GetServiceAccountRoles(ctx context.Context, saID uuid.UUID, projectID *uuid.UUID) ([]*models.ServiceAccountRole, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByExternalID(ctx context.Context, provider string, externalID string) (*models.User, error)
	List(ctx context.Context, limit, offset int) ([]*models.User, int, error)
}

type InvitationRepository interface {
	Create(ctx context.Context, inv *models.Invitation) error
	GetByEmail(ctx context.Context, email string) (*models.Invitation, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*models.Invitation, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SystemConfigRepository interface {
	GetByKey(ctx context.Context, key string) (*models.SystemConfig, error)
	Upsert(ctx context.Context, config *models.SystemConfig) error
}

type ServiceAccountRepository interface {
	Create(ctx context.Context, sa *models.ServiceAccount) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ServiceAccount, error)
	CreateKey(ctx context.Context, key *models.ServiceAccountKey) error
	GetKeyByHash(ctx context.Context, keyHash string) (*models.ServiceAccountKey, error)
	ListKeys(ctx context.Context, saID uuid.UUID) ([]*models.ServiceAccountKey, error)
}
