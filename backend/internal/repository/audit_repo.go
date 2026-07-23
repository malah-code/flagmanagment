package repository

import (
	"context"
	"database/sql"
	"github.com/flagmanagment/backend/internal/models"

	"github.com/google/uuid"
)

type auditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Create(ctx context.Context, log *models.AuditLog) error {
	query := `INSERT INTO audit_logs (id, project_id, environment_id, actor_id, action, target_type, target_id, previous_state, new_state, actor_ip, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query, log.ID, log.ProjectID, log.EnvironmentID, log.ActorID, log.Action, log.TargetType, log.TargetID, log.PreviousState, log.NewState, log.ActorIP, log.CreatedAt)
	return err
}

func (r *auditRepository) ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error) {
	query := `
		SELECT id, project_id, environment_id, actor_id, action, target_type, target_id, previous_state, new_state, actor_ip, created_at
		FROM audit_logs
		WHERE project_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		if err := rows.Scan(&log.ID, &log.ProjectID, &log.EnvironmentID, &log.ActorID, &log.Action, &log.TargetType, &log.TargetID, &log.PreviousState, &log.NewState, &log.ActorIP, &log.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, &log)
	}

	var total int
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_logs WHERE project_id = $1`, projectID).Scan(&total)
	return logs, total, err
}

func (r *auditRepository) ListByEnvironment(ctx context.Context, envID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error) {
	query := `
		SELECT id, project_id, environment_id, actor_id, action, target_type, target_id, previous_state, new_state, actor_ip, created_at
		FROM audit_logs
		WHERE environment_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, envID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		if err := rows.Scan(&log.ID, &log.ProjectID, &log.EnvironmentID, &log.ActorID, &log.Action, &log.TargetType, &log.TargetID, &log.PreviousState, &log.NewState, &log.ActorIP, &log.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, &log)
	}

	var total int
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_logs WHERE environment_id = $1`, envID).Scan(&total)
	return logs, total, err
}
