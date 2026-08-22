package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/services"
)

type userRepository struct {
	db     *sql.DB
	crypto services.CryptoService
}

func NewUserRepository(db *sql.DB, crypto services.CryptoService) UserRepository {
	return &userRepository{db: db, crypto: crypto}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, email_hash, password_hash, auth_provider, external_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	if user.AuthProvider == "" {
		user.AuthProvider = "local"
	}
	
	emailHash := r.crypto.HashToken(user.Email)
	encryptedEmail, err := r.crypto.EncryptAES(user.Email)
	if err != nil {
		return err
	}
	user.EmailHash = emailHash

	err = r.db.QueryRowContext(ctx, query, user.ID, encryptedEmail, emailHash, user.PasswordHash, user.AuthProvider, user.ExternalID).Scan(&user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	emailHash := r.crypto.HashToken(email)
	query := `SELECT id, email, email_hash, password_hash, auth_provider, external_id, created_at, updated_at FROM users WHERE email_hash = $1`
	row := r.db.QueryRowContext(ctx, query, emailHash)

	var user models.User
	var encEmail string
	err := row.Scan(&user.ID, &encEmail, &user.EmailHash, &user.PasswordHash, &user.AuthProvider, &user.ExternalID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	
	decEmail, _ := r.crypto.DecryptAES(encEmail)
	user.Email = decEmail
	
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, email_hash, password_hash, auth_provider, external_id, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var user models.User
	var encEmail string
	err := row.Scan(&user.ID, &encEmail, &user.EmailHash, &user.PasswordHash, &user.AuthProvider, &user.ExternalID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	
	decEmail, _ := r.crypto.DecryptAES(encEmail)
	user.Email = decEmail
	
	return &user, nil
}

func (r *userRepository) GetByExternalID(ctx context.Context, provider string, externalID string) (*models.User, error) {
	query := `SELECT id, email, email_hash, password_hash, auth_provider, external_id, created_at, updated_at FROM users WHERE auth_provider = $1 AND external_id = $2`
	row := r.db.QueryRowContext(ctx, query, provider, externalID)

	var user models.User
	var encEmail string
	err := row.Scan(&user.ID, &encEmail, &user.EmailHash, &user.PasswordHash, &user.AuthProvider, &user.ExternalID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	
	decEmail, _ := r.crypto.DecryptAES(encEmail)
	user.Email = decEmail
	
	return &user, nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*models.User, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, email, email_hash, password_hash, auth_provider, external_id, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		var encEmail string
		err := rows.Scan(&user.ID, &encEmail, &user.EmailHash, &user.PasswordHash, &user.AuthProvider, &user.ExternalID, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		
		decEmail, _ := r.crypto.DecryptAES(encEmail)
		user.Email = decEmail
		
		users = append(users, &user)
	}

	return users, total, nil
}
