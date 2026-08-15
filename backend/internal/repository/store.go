package repository

import (
	"context"
	"database/sql"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &store{db: db}
}

func (s *store) ProjectRepo() ProjectRepository {
	return NewProjectRepository(s.db)
}

func (s *store) EnvironmentRepo() EnvironmentRepository {
	return NewEnvironmentRepository(s.db)
}

func (s *store) FlagRepo() FlagRepository {
	return NewFlagRepository(s.db)
}

func (s *store) FlagStateRepo() FlagStateRepository {
	return NewFlagStateRepository(s.db)
}

func (s *store) AuditRepo() AuditRepository {
	return NewAuditRepository(s.db)
}

func (s *store) ChangeRequestRepo() ChangeRequestRepository {
	return NewChangeRequestRepository(s.db)
}

func (s *store) KillSwitchRepo() KillSwitchRepository {
	return NewKillSwitchRepository(s.db)
}

func (s *store) SlackConfigRepo() SlackConfigRepository {
	return NewSlackConfigRepository(s.db)
}

func (s *store) ScheduledChangeRepo() ScheduledChangeRepository {
	return NewScheduledChangeRepository(s.db)
}

func (s *store) StalePolicyRepo() StalePolicyRepository {
	return NewStalePolicyRepository(s.db)
}

func (s *store) RoleRepo() RoleRepository {
	return NewRoleRepository(s.db)
}

func (s *store) UserRepo() UserRepository {
	return NewUserRepository(s.db)
}

func (s *store) InvitationRepo() InvitationRepository {
	return NewInvitationRepository(s.db)
}

func (s *store) SystemConfigRepo() SystemConfigRepository {
	return NewSystemConfigRepository(s.db)
}

func (s *store) WebhookIntegrationRepo() WebhookIntegrationRepository {
	return NewWebhookIntegrationRepository(s.db)
}

func (s *store) ServiceAccountRepo() ServiceAccountRepository {
	return NewServiceAccountRepository(s.db)
}

func (s *store) WithTx(ctx context.Context, fn func(Store) error) error {
	// Simple non-transactional stub. Proper tx management requires extending Store to wrap *sql.Tx
	return fn(s)
}

func (s *store) MigrateUp() error {
	return nil
}

func (s *store) MigrateDown() error {
	return nil
}
