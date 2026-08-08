package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/google/uuid"
)

type flagStateRepository struct {
	db *sql.DB
}

func NewFlagStateRepository(db *sql.DB) FlagStateRepository {
	return &flagStateRepository{db: db}
}

const selectFlagStateCols = `id, environment_id, feature_flag_id, enabled, targeting_rules, remote_config, variations, COALESCE(lifecycle_state, 'ACTIVE'), last_evaluated_at, COALESCE(last_state_change_at, created_at), created_at, updated_at`

func (r *flagStateRepository) Create(ctx context.Context, state *models.EnvironmentFlagState) error {
	if state.LifecycleState == "" {
		state.LifecycleState = models.LifecycleActive
	}
	if state.LastStateChangeAt.IsZero() {
		state.LastStateChangeAt = time.Now()
	}
	query := `INSERT INTO environment_flag_states (id, environment_id, feature_flag_id, enabled, targeting_rules, remote_config, variations, lifecycle_state, last_evaluated_at, last_state_change_at, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query, state.ID, state.EnvironmentID, state.FeatureFlagID, state.Enabled, state.TargetingRules, state.RemoteConfig, state.Variations, state.LifecycleState, state.LastEvaluatedAt, state.LastStateChangeAt, state.CreatedAt, state.UpdatedAt)
	return err
}

func (r *flagStateRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.EnvironmentFlagState, error) {
	query := `SELECT ` + selectFlagStateCols + ` FROM environment_flag_states WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var s models.EnvironmentFlagState
	if err := row.Scan(&s.ID, &s.EnvironmentID, &s.FeatureFlagID, &s.Enabled, &s.TargetingRules, &s.RemoteConfig, &s.Variations, &s.LifecycleState, &s.LastEvaluatedAt, &s.LastStateChangeAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *flagStateRepository) GetByEnvAndFlag(ctx context.Context, envID, flagID uuid.UUID) (*models.EnvironmentFlagState, error) {
	query := `SELECT ` + selectFlagStateCols + ` FROM environment_flag_states WHERE environment_id = $1 AND feature_flag_id = $2`
	row := r.db.QueryRowContext(ctx, query, envID, flagID)
	var s models.EnvironmentFlagState
	if err := row.Scan(&s.ID, &s.EnvironmentID, &s.FeatureFlagID, &s.Enabled, &s.TargetingRules, &s.RemoteConfig, &s.Variations, &s.LifecycleState, &s.LastEvaluatedAt, &s.LastStateChangeAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *flagStateRepository) ListByEnvironment(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentFlagState, error) {
	query := `SELECT ` + selectFlagStateCols + ` FROM environment_flag_states WHERE environment_id = $1`
	rows, err := r.db.QueryContext(ctx, query, envID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*models.EnvironmentFlagState
	for rows.Next() {
		var s models.EnvironmentFlagState
		if err := rows.Scan(&s.ID, &s.EnvironmentID, &s.FeatureFlagID, &s.Enabled, &s.TargetingRules, &s.RemoteConfig, &s.Variations, &s.LifecycleState, &s.LastEvaluatedAt, &s.LastStateChangeAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		states = append(states, &s)
	}
	return states, nil
}

func (r *flagStateRepository) ListByEnvironmentAndLifecycle(ctx context.Context, envID uuid.UUID, lifecycle models.FlagLifecycleState) ([]*models.EnvironmentFlagState, error) {
	query := `SELECT ` + selectFlagStateCols + ` FROM environment_flag_states WHERE environment_id = $1 AND lifecycle_state = $2`
	rows, err := r.db.QueryContext(ctx, query, envID, lifecycle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*models.EnvironmentFlagState
	for rows.Next() {
		var s models.EnvironmentFlagState
		if err := rows.Scan(&s.ID, &s.EnvironmentID, &s.FeatureFlagID, &s.Enabled, &s.TargetingRules, &s.RemoteConfig, &s.Variations, &s.LifecycleState, &s.LastEvaluatedAt, &s.LastStateChangeAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		states = append(states, &s)
	}
	return states, nil
}

func (r *flagStateRepository) Update(ctx context.Context, state *models.EnvironmentFlagState) error {
	now := time.Now()
	// If rule or state changed, reset lifecycle to ACTIVE and update last_state_change_at
	query := `UPDATE environment_flag_states SET enabled = $1, targeting_rules = $2, remote_config = $3, variations = $4, lifecycle_state = 'ACTIVE', last_state_change_at = $5, updated_at = $5 WHERE id = $6`
	res, err := r.db.ExecContext(ctx, query, state.Enabled, state.TargetingRules, state.RemoteConfig, state.Variations, now, state.ID)
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

func (r *flagStateRepository) UpdateLifecycleState(ctx context.Context, id uuid.UUID, state models.FlagLifecycleState) error {
	now := time.Now()
	query := `UPDATE environment_flag_states SET lifecycle_state = $1, updated_at = $2 WHERE id = $3`
	res, err := r.db.ExecContext(ctx, query, state, now, id)
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

func (r *flagStateRepository) UpdateLastEvaluatedAtBatch(ctx context.Context, updates map[uuid.UUID]time.Time) error {
	if len(updates) == 0 {
		return nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `UPDATE environment_flag_states SET last_evaluated_at = $1 WHERE id = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id, evaluatedAt := range updates {
		if _, err := stmt.ExecContext(ctx, evaluatedAt, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *flagStateRepository) FindActiveFlagsForStalenessScan(ctx context.Context, limit int) ([]*models.EnvironmentFlagState, error) {
	if limit <= 0 {
		limit = 500
	}
	query := `SELECT ` + selectFlagStateCols + ` FROM environment_flag_states WHERE lifecycle_state = 'ACTIVE' LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*models.EnvironmentFlagState
	for rows.Next() {
		var s models.EnvironmentFlagState
		if err := rows.Scan(&s.ID, &s.EnvironmentID, &s.FeatureFlagID, &s.Enabled, &s.TargetingRules, &s.RemoteConfig, &s.Variations, &s.LifecycleState, &s.LastEvaluatedAt, &s.LastStateChangeAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		states = append(states, &s)
	}
	return states, nil
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
