package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/services"
)

type invitationRepository struct {
	db     *sql.DB
	crypto services.CryptoService
}

func NewInvitationRepository(db *sql.DB, crypto services.CryptoService) InvitationRepository {
	return &invitationRepository{db: db, crypto: crypto}
}

func (r *invitationRepository) Create(ctx context.Context, inv *models.Invitation) error {
	query := `
		INSERT INTO invitations (id, email, email_hash, token_hash, role, project_ids, expires_at, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING created_at
	`
	if inv.ID == uuid.Nil {
		inv.ID = uuid.New()
	}
	
	emailHash := r.crypto.HashToken(inv.Email)
	encryptedEmail, err := r.crypto.EncryptAES(inv.Email)
	if err != nil {
		return err
	}
	inv.EmailHash = emailHash
	
	err = r.db.QueryRowContext(ctx, query, inv.ID, encryptedEmail, emailHash, inv.TokenHash, inv.Role, inv.ProjectIDs, inv.ExpiresAt, inv.CreatedBy).Scan(&inv.CreatedAt)
	return err
}

func (r *invitationRepository) GetByEmail(ctx context.Context, email string) (*models.Invitation, error) {
	emailHash := r.crypto.HashToken(email)
	query := `SELECT id, email, email_hash, token_hash, role, project_ids, expires_at, created_by, created_at FROM invitations WHERE email_hash = $1`
	row := r.db.QueryRowContext(ctx, query, emailHash)
	
	var inv models.Invitation
	var encEmail string
	err := row.Scan(&inv.ID, &encEmail, &inv.EmailHash, &inv.TokenHash, &inv.Role, &inv.ProjectIDs, &inv.ExpiresAt, &inv.CreatedBy, &inv.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	
	decEmail, _ := r.crypto.DecryptAES(encEmail)
	inv.Email = decEmail
	
	return &inv, nil
}

func (r *invitationRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.Invitation, error) {
	query := `SELECT id, email, email_hash, token_hash, role, project_ids, expires_at, created_by, created_at FROM invitations WHERE token_hash = $1`
	row := r.db.QueryRowContext(ctx, query, tokenHash)
	
	var inv models.Invitation
	var encEmail string
	err := row.Scan(&inv.ID, &encEmail, &inv.EmailHash, &inv.TokenHash, &inv.Role, &inv.ProjectIDs, &inv.ExpiresAt, &inv.CreatedBy, &inv.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	
	decEmail, _ := r.crypto.DecryptAES(encEmail)
	inv.Email = decEmail
	
	return &inv, nil
}

func (r *invitationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM invitations WHERE id = $1", id)
	return err
}
