package repository

import (
	"context"
	"database/sql"
	"github.com/flagmanagment/backend/internal/models"

	"github.com/google/uuid"
)

type changeRequestRepository struct {
	db *sql.DB
}

func NewChangeRequestRepository(db *sql.DB) ChangeRequestRepository {
	return &changeRequestRepository{db: db}
}

func (r *changeRequestRepository) Create(ctx context.Context, cr *models.ChangeRequest) error {
	query := `INSERT INTO change_requests (id, project_id, environment_id, title, description, status, proposed_changes, created_by, applied_by, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query, cr.ID, cr.ProjectID, cr.EnvironmentID, cr.Title, cr.Description, cr.Status, cr.ProposedChanges, cr.CreatedBy, cr.AppliedBy, cr.CreatedAt, cr.UpdatedAt)
	return err
}

func (r *changeRequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ChangeRequest, error) {
	query := `SELECT id, project_id, environment_id, title, description, status, proposed_changes, created_by, applied_by, created_at, updated_at FROM change_requests WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var cr models.ChangeRequest
	if err := row.Scan(&cr.ID, &cr.ProjectID, &cr.EnvironmentID, &cr.Title, &cr.Description, &cr.Status, &cr.ProposedChanges, &cr.CreatedBy, &cr.AppliedBy, &cr.CreatedAt, &cr.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &cr, nil
}

func (r *changeRequestRepository) ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*models.ChangeRequest, int, error) {
	return nil, 0, nil // Stub
}

func (r *changeRequestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, appliedBy *uuid.UUID) error {
	query := `UPDATE change_requests SET status = $1, applied_by = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, appliedBy, id)
	return err
}

func (r *changeRequestRepository) AddApproval(ctx context.Context, approval *models.ChangeRequestApproval) error {
	query := `INSERT INTO change_request_approvals (id, change_request_id, approver_id, decision, comment, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, approval.ID, approval.ChangeRequestID, approval.ApproverID, approval.Decision, approval.Comment, approval.CreatedAt)
	return err
}

func (r *changeRequestRepository) ListApprovals(ctx context.Context, crID uuid.UUID) ([]*models.ChangeRequestApproval, error) {
	query := `SELECT id, change_request_id, approver_id, decision, comment, created_at FROM change_request_approvals WHERE change_request_id = $1`
	rows, err := r.db.QueryContext(ctx, query, crID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var approvals []*models.ChangeRequestApproval
	for rows.Next() {
		var a models.ChangeRequestApproval
		if err := rows.Scan(&a.ID, &a.ChangeRequestID, &a.ApproverID, &a.Decision, &a.Comment, &a.CreatedAt); err != nil {
			return nil, err
		}
		approvals = append(approvals, &a)
	}
	return approvals, nil
}
