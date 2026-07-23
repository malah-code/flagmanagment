package repository

import (
	"context"
	"database/sql"
	"github.com/flagmanagment/backend/internal/models"

	"github.com/google/uuid"
)

type flagStateRepository struct {
	db *sql.DB
}

func NewFlagStateRepository(db *sql.DB) FlagStateRepository {
	return &flagStateRepository{db: db}
}

func (r *flagStateRepository) Create(ctx context.Context, state *models.EnvironmentFlagState) error {
	query := `INSERT INTO environment_flag_states (id, environment_id, feature_flag_id, enabled, targeting_rules, remote_config, variations, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, state.ID, state.EnvironmentID, state.FeatureFlagID, state.Enabled, state.TargetingRules, state.RemoteConfig, state.Variations, state.CreatedAt, state.UpdatedAt)
	return err
}

func (r *flagStateRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.EnvironmentFlagState, error) {
	query := `SELECT id, environment_id, feature_flag_id, enabled, targeting_rules, remote_config, variations, created_at, updated_at FROM environment_flag_states WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var s models.EnvironmentFlagState
	if err := row.Scan(&s.ID, &s.EnvironmentID, &s.FeatureFlagID, &s.Enabled, &s.TargetingRules, &s.RemoteConfig, &s.Variations, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *flagStateRepository) GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*models.EnvironmentFlagState, error) {
	query := `SELECT id, environment_id, feature_flag_id, enabled, targeting_rules, remote_config, variations, created_at, updated_at FROM environment_flag_states WHERE environment_id = $1 AND feature_flag_id = $2`
	row := r.db.QueryRowContext(ctx, query, envID, flagID)
	var s models.EnvironmentFlagState
	if err := row.Scan(&s.ID, &s.EnvironmentID, &s.FeatureFlagID, &s.Enabled, &s.TargetingRules, &s.RemoteConfig, &s.Variations, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *flagStateRepository) ListByEnvironment(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentFlagState, error) {
	query := `SELECT id, environment_id, feature_flag_id, enabled, targeting_rules, remote_config, variations, created_at, updated_at FROM environment_flag_states WHERE environment_id = $1`
	rows, err := r.db.QueryContext(ctx, query, envID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*models.EnvironmentFlagState
	for rows.Next() {
		var s models.EnvironmentFlagState
		if err := rows.Scan(&s.ID, &s.EnvironmentID, &s.FeatureFlagID, &s.Enabled, &s.TargetingRules, &s.RemoteConfig, &s.Variations, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		states = append(states, &s)
	}
	return states, nil
}

func (r *flagStateRepository) Update(ctx context.Context, state *models.EnvironmentFlagState) error {
	query := `UPDATE environment_flag_states SET enabled = $1, targeting_rules = $2, remote_config = $3, variations = $4, updated_at = $5 WHERE id = $6`
	res, err := r.db.ExecContext(ctx, query, state.Enabled, state.TargetingRules, state.RemoteConfig, state.Variations, state.UpdatedAt, state.ID)
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

func (r *flagStateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM environment_flag_states WHERE id = $1`
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
