package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/google/uuid"
)

type stalePolicyRepository struct {
	db *sql.DB
}

func NewStalePolicyRepository(db *sql.DB) StalePolicyRepository {
	return &stalePolicyRepository{db: db}
}

func (r *stalePolicyRepository) GetByEnvironment(ctx context.Context, projectID, envID uuid.UUID) (*models.StaleFlagPolicy, error) {
	query := `SELECT id, project_id, environment_id, stale_after_days, created_at, updated_at 
	          FROM stale_flag_policies 
	          WHERE project_id = $1 AND (environment_id = $2 OR environment_id IS NULL)
	          ORDER BY environment_id DESC NULLS LAST LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, projectID, envID)
	var p models.StaleFlagPolicy
	if err := row.Scan(&p.ID, &p.ProjectID, &p.EnvironmentID, &p.StaleAfterDays, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *stalePolicyRepository) Upsert(ctx context.Context, policy *models.StaleFlagPolicy) error {
	now := time.Now()
	if policy.ID == uuid.Nil {
		policy.ID = uuid.New()
	}
	policy.UpdatedAt = now
	if policy.CreatedAt.IsZero() {
		policy.CreatedAt = now
	}
	query := `INSERT INTO stale_flag_policies (id, project_id, environment_id, stale_after_days, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          ON CONFLICT (project_id, COALESCE(environment_id, '00000000-0000-0000-0000-000000000000'))
	          DO UPDATE SET stale_after_days = EXCLUDED.stale_after_days, updated_at = EXCLUDED.updated_at`
	_, err := r.db.ExecContext(ctx, query, policy.ID, policy.ProjectID, policy.EnvironmentID, policy.StaleAfterDays, policy.CreatedAt, policy.UpdatedAt)
	return err
}

func (r *stalePolicyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM stale_flag_policies WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
