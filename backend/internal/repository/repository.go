package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrConflict      = errors.New("resource conflict")
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
	RoleRepo() RoleRepository
	UserRepo() UserRepository

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
}

type FlagRepository interface {
	Create(ctx context.Context, flag *models.FeatureFlag) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.FeatureFlag, error)
	GetByKey(ctx context.Context, projectID uuid.UUID, key string) (*models.FeatureFlag, error)
	ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*models.FeatureFlag, int, error)
	Update(ctx context.Context, flag *models.FeatureFlag) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastEvaluatedAt(ctx context.Context, ids []uuid.UUID) error
}

type FlagStateRepository interface {
	Create(ctx context.Context, state *models.EnvironmentFlagState) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.EnvironmentFlagState, error)
	GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*models.EnvironmentFlagState, error)
	ListByEnvironment(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentFlagState, error)
	Update(ctx context.Context, state *models.EnvironmentFlagState) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AuditRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error)
	ListByEnvironment(ctx context.Context, envID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error)
}

type ChangeRequestRepository interface {
	Create(ctx context.Context, cr *models.ChangeRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ChangeRequest, error)
	ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*models.ChangeRequest, int, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, appliedBy *uuid.UUID) error
	AddApproval(ctx context.Context, approval *models.ChangeRequestApproval) error
	ListApprovals(ctx context.Context, crID uuid.UUID) ([]*models.ChangeRequestApproval, error)
}

type RoleRepository interface {
	Create(ctx context.Context, role *models.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
	List(ctx context.Context) ([]*models.Role, error)
	AssignUserRole(ctx context.Context, ur *models.UserRole) error
	RemoveUserRole(ctx context.Context, id uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID, projectID *uuid.UUID) ([]*models.UserRole, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}
