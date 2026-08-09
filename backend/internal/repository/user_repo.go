package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, auth_provider, external_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	if user.AuthProvider == "" {
		user.AuthProvider = "local"
	}
	err := r.db.QueryRowContext(ctx, query, user.ID, user.Email, user.PasswordHash, user.AuthProvider, user.ExternalID).Scan(&user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, auth_provider, external_id, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.AuthProvider, &user.ExternalID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, password_hash, auth_provider, external_id, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.AuthProvider, &user.ExternalID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByExternalID(ctx context.Context, provider string, externalID string) (*models.User, error) {
	query := `SELECT id, email, password_hash, auth_provider, external_id, created_at, updated_at FROM users WHERE auth_provider = $1 AND external_id = $2`
	row := r.db.QueryRowContext(ctx, query, provider, externalID)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.AuthProvider, &user.ExternalID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
