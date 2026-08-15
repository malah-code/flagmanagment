package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
)

type invitationRepository struct {
	db *sql.DB
}

func NewInvitationRepository(db *sql.DB) InvitationRepository {
	return &invitationRepository{db: db}
}

func (r *invitationRepository) Create(ctx context.Context, inv *models.Invitation) error {
	query := `
		INSERT INTO invitations (id, email, token_hash, role, project_ids, expires_at, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING created_at
	`
	if inv.ID == uuid.Nil {
		inv.ID = uuid.New()
	}
	
	err := r.db.QueryRowContext(ctx, query, inv.ID, inv.Email, inv.TokenHash, inv.Role, inv.ProjectIDs, inv.ExpiresAt, inv.CreatedBy).Scan(&inv.CreatedAt)
	return err
}

func (r *invitationRepository) GetByEmail(ctx context.Context, email string) (*models.Invitation, error) {
	query := `SELECT id, email, token_hash, role, project_ids, expires_at, created_by, created_at FROM invitations WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)
	
	var inv models.Invitation
	err := row.Scan(&inv.ID, &inv.Email, &inv.TokenHash, &inv.Role, &inv.ProjectIDs, &inv.ExpiresAt, &inv.CreatedBy, &inv.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &inv, err
}

func (r *invitationRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.Invitation, error) {
	query := `SELECT id, email, token_hash, role, project_ids, expires_at, created_by, created_at FROM invitations WHERE token_hash = $1`
	row := r.db.QueryRowContext(ctx, query, tokenHash)
	
	var inv models.Invitation
	err := row.Scan(&inv.ID, &inv.Email, &inv.TokenHash, &inv.Role, &inv.ProjectIDs, &inv.ExpiresAt, &inv.CreatedBy, &inv.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &inv, err
}

func (r *invitationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM invitations WHERE id = $1", id)
	return err
}
