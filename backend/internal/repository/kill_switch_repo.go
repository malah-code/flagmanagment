package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
)

type KillSwitchRepository interface {
	Create(ctx context.Context, ks *models.KillSwitchRule) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.KillSwitchRule, error)
	ListByEnvironmentAndFlag(ctx context.Context, envID uuid.UUID, flagID uuid.UUID) ([]*models.KillSwitchRule, error)
	ListByEnvironmentAndAlert(ctx context.Context, envID uuid.UUID, alertIdentifier string) ([]*models.KillSwitchRule, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type killSwitchRepository struct {
	db *sql.DB
}

func NewKillSwitchRepository(db *sql.DB) KillSwitchRepository {
	return &killSwitchRepository{db: db}
}

func (r *killSwitchRepository) Create(ctx context.Context, ks *models.KillSwitchRule) error {
	query := `INSERT INTO kill_switches (id, flag_id, environment_id, alert_identifier, action, created_by, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, ks.ID, ks.FlagID, ks.EnvironmentID, ks.AlertIdentifier, ks.Action, ks.CreatedBy, ks.CreatedAt, ks.UpdatedAt)
	return err
}

func (r *killSwitchRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.KillSwitchRule, error) {
	query := `SELECT id, flag_id, environment_id, alert_identifier, action, created_by, created_at, updated_at FROM kill_switches WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	
	var ks models.KillSwitchRule
	err := row.Scan(&ks.ID, &ks.FlagID, &ks.EnvironmentID, &ks.AlertIdentifier, &ks.Action, &ks.CreatedBy, &ks.CreatedAt, &ks.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &ks, nil
}

func (r *killSwitchRepository) ListByEnvironmentAndFlag(ctx context.Context, envID uuid.UUID, flagID uuid.UUID) ([]*models.KillSwitchRule, error) {
	query := `SELECT id, flag_id, environment_id, alert_identifier, action, created_by, created_at, updated_at FROM kill_switches WHERE environment_id = $1 AND flag_id = $2`
	rows, err := r.db.QueryContext(ctx, query, envID, flagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*models.KillSwitchRule
	for rows.Next() {
		var ks models.KillSwitchRule
		if err := rows.Scan(&ks.ID, &ks.FlagID, &ks.EnvironmentID, &ks.AlertIdentifier, &ks.Action, &ks.CreatedBy, &ks.CreatedAt, &ks.UpdatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, &ks)
	}
	return rules, nil
}

func (r *killSwitchRepository) ListByEnvironmentAndAlert(ctx context.Context, envID uuid.UUID, alertIdentifier string) ([]*models.KillSwitchRule, error) {
	query := `SELECT id, flag_id, environment_id, alert_identifier, action, created_by, created_at, updated_at FROM kill_switches WHERE environment_id = $1 AND alert_identifier = $2`
	rows, err := r.db.QueryContext(ctx, query, envID, alertIdentifier)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*models.KillSwitchRule
	for rows.Next() {
		var ks models.KillSwitchRule
		if err := rows.Scan(&ks.ID, &ks.FlagID, &ks.EnvironmentID, &ks.AlertIdentifier, &ks.Action, &ks.CreatedBy, &ks.CreatedAt, &ks.UpdatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, &ks)
	}
	return rules, nil
}

func (r *killSwitchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM kill_switches WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
