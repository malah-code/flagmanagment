package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type scheduledChangeRepository struct {
	db *sql.DB
}

func NewScheduledChangeRepository(db *sql.DB) ScheduledChangeRepository {
	return &scheduledChangeRepository{db: db}
}

func (r *scheduledChangeRepository) Create(ctx context.Context, sc *models.ScheduledChange) error {
	query := `INSERT INTO scheduled_changes 
		(id, project_id, environment_id, target_type, target_id, action, scheduled_for, status, created_by, executed_at, cancelled_at, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.db.ExecContext(ctx, query,
		sc.ID, sc.ProjectID, sc.EnvironmentID, sc.TargetType, sc.TargetID, sc.Action,
		sc.ScheduledFor, sc.Status, sc.CreatedBy, sc.ExecutedAt, sc.CancelledAt,
		sc.CreatedAt, sc.UpdatedAt,
	)
	if err != nil {
		// Map Postgres unique_violation (23505) on the pending-flag partial index
		// to a clean sentinel so service callers don't have to inspect raw driver errors.
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrPendingScheduleExists
		}
		return err
	}
	return nil
}

func (r *scheduledChangeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ScheduledChange, error) {
	query := `SELECT id, project_id, environment_id, target_type, target_id, action, scheduled_for, status, created_by, executed_at, cancelled_at, created_at, updated_at 
		FROM scheduled_changes WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var sc models.ScheduledChange
	err := row.Scan(
		&sc.ID, &sc.ProjectID, &sc.EnvironmentID, &sc.TargetType, &sc.TargetID, &sc.Action,
		&sc.ScheduledFor, &sc.Status, &sc.CreatedBy, &sc.ExecutedAt, &sc.CancelledAt,
		&sc.CreatedAt, &sc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &sc, nil
}

func (r *scheduledChangeRepository) ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*models.ScheduledChange, int, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `SELECT id, project_id, environment_id, target_type, target_id, action, scheduled_for, status, created_by, executed_at, cancelled_at, created_at, updated_at, COUNT(*) OVER() 
			FROM scheduled_changes WHERE environment_id = $1 AND status = $2 ORDER BY scheduled_for DESC LIMIT $3 OFFSET $4`
		args = []interface{}{envID, status, limit, offset}
	} else {
		query = `SELECT id, project_id, environment_id, target_type, target_id, action, scheduled_for, status, created_by, executed_at, cancelled_at, created_at, updated_at, COUNT(*) OVER() 
			FROM scheduled_changes WHERE environment_id = $1 ORDER BY scheduled_for DESC LIMIT $2 OFFSET $3`
		args = []interface{}{envID, limit, offset}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var schedules []*models.ScheduledChange
	total := 0

	for rows.Next() {
		var sc models.ScheduledChange
		if err := rows.Scan(
			&sc.ID, &sc.ProjectID, &sc.EnvironmentID, &sc.TargetType, &sc.TargetID, &sc.Action,
			&sc.ScheduledFor, &sc.Status, &sc.CreatedBy, &sc.ExecutedAt, &sc.CancelledAt,
			&sc.CreatedAt, &sc.UpdatedAt, &total,
		); err != nil {
			return nil, 0, err
		}
		schedules = append(schedules, &sc)
	}

	return schedules, total, nil
}

func (r *scheduledChangeRepository) GetPendingByTargetID(ctx context.Context, targetID uuid.UUID) (*models.ScheduledChange, error) {
	query := `SELECT id, project_id, environment_id, target_type, target_id, action, scheduled_for, status, created_by, executed_at, cancelled_at, created_at, updated_at 
		FROM scheduled_changes WHERE target_id = $1 AND status = 'PENDING' LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, targetID)

	var sc models.ScheduledChange
	err := row.Scan(
		&sc.ID, &sc.ProjectID, &sc.EnvironmentID, &sc.TargetType, &sc.TargetID, &sc.Action,
		&sc.ScheduledFor, &sc.Status, &sc.CreatedBy, &sc.ExecutedAt, &sc.CancelledAt,
		&sc.CreatedAt, &sc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &sc, nil
}

func (r *scheduledChangeRepository) GetDueSchedules(ctx context.Context, now time.Time, limit int) ([]*models.ScheduledChange, error) {
	query := `SELECT id, project_id, environment_id, target_type, target_id, action, scheduled_for, status, created_by, executed_at, cancelled_at, created_at, updated_at 
		FROM scheduled_changes WHERE status = 'PENDING' AND scheduled_for <= $1 ORDER BY scheduled_for ASC LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*models.ScheduledChange
	for rows.Next() {
		var sc models.ScheduledChange
		if err := rows.Scan(
			&sc.ID, &sc.ProjectID, &sc.EnvironmentID, &sc.TargetType, &sc.TargetID, &sc.Action,
			&sc.ScheduledFor, &sc.Status, &sc.CreatedBy, &sc.ExecutedAt, &sc.CancelledAt,
			&sc.CreatedAt, &sc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		schedules = append(schedules, &sc)
	}
	return schedules, nil
}

func (r *scheduledChangeRepository) MarkExecuted(ctx context.Context, id uuid.UUID, executedAt time.Time) error {
	query := `UPDATE scheduled_changes SET status = 'EXECUTED', executed_at = $2, updated_at = $2 WHERE id = $1 AND status = 'PENDING'`
	res, err := r.db.ExecContext(ctx, query, id, executedAt)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("scheduled change %s not found or not pending", id)
	}
	return nil
}

func (r *scheduledChangeRepository) MarkCancelled(ctx context.Context, id uuid.UUID, cancelledAt time.Time) error {
	query := `UPDATE scheduled_changes SET status = 'CANCELLED', cancelled_at = $2, updated_at = $2 WHERE id = $1 AND status = 'PENDING'`
	res, err := r.db.ExecContext(ctx, query, id, cancelledAt)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("scheduled change %s not found or not pending", id)
	}
	return nil
}

func (r *scheduledChangeRepository) UpdateScheduledFor(ctx context.Context, id uuid.UUID, newTime time.Time) error {
	// Use $3 (newTime) for updated_at to be consistent with MarkExecuted/MarkCancelled
	// which also use the application-provided timestamp rather than DB NOW().
	query := `UPDATE scheduled_changes SET scheduled_for = $2, updated_at = $2 WHERE id = $1 AND status = 'PENDING'`
	res, err := r.db.ExecContext(ctx, query, id, newTime)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("scheduled change %s not found or not pending", id)
	}
	return nil
}
