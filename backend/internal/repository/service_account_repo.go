package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
)

type serviceAccountRepository struct {
	db *sql.DB
}

func NewServiceAccountRepository(db *sql.DB) ServiceAccountRepository {
	return &serviceAccountRepository{db: db}
}

func (r *serviceAccountRepository) Create(ctx context.Context, sa *models.ServiceAccount) error {
	query := `
		INSERT INTO service_accounts (id, name, description, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	if sa.ID == uuid.Nil {
		sa.ID = uuid.New()
	}
	err := r.db.QueryRowContext(ctx, query, sa.ID, sa.Name, sa.Description, sa.CreatedBy).Scan(&sa.CreatedAt, &sa.UpdatedAt)
	return err
}

func (r *serviceAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ServiceAccount, error) {
	query := `SELECT id, name, description, created_by, created_at, updated_at FROM service_accounts WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var sa models.ServiceAccount
	err := row.Scan(&sa.ID, &sa.Name, &sa.Description, &sa.CreatedBy, &sa.CreatedAt, &sa.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &sa, nil
}

func (r *serviceAccountRepository) CreateKey(ctx context.Context, key *models.ServiceAccountKey) error {
	query := `
		INSERT INTO service_account_keys (id, service_account_id, key_hash, name, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	if key.ID == uuid.Nil {
		key.ID = uuid.New()
	}
	err := r.db.QueryRowContext(ctx, query, key.ID, key.ServiceAccountID, key.KeyHash, key.Name, key.ExpiresAt).Scan(&key.CreatedAt, &key.UpdatedAt)
	return err
}

func (r *serviceAccountRepository) GetKeyByHash(ctx context.Context, keyHash string) (*models.ServiceAccountKey, error) {
	query := `SELECT id, service_account_id, key_hash, name, expires_at, last_used_at, created_at, updated_at FROM service_account_keys WHERE key_hash = $1`
	row := r.db.QueryRowContext(ctx, query, keyHash)

	var key models.ServiceAccountKey
	err := row.Scan(&key.ID, &key.ServiceAccountID, &key.KeyHash, &key.Name, &key.ExpiresAt, &key.LastUsedAt, &key.CreatedAt, &key.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (r *serviceAccountRepository) ListKeys(ctx context.Context, saID uuid.UUID) ([]*models.ServiceAccountKey, error) {
	query := `SELECT id, service_account_id, key_hash, name, expires_at, last_used_at, created_at, updated_at FROM service_account_keys WHERE service_account_id = $1`
	rows, err := r.db.QueryContext(ctx, query, saID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*models.ServiceAccountKey
	for rows.Next() {
		var key models.ServiceAccountKey
		if err := rows.Scan(&key.ID, &key.ServiceAccountID, &key.KeyHash, &key.Name, &key.ExpiresAt, &key.LastUsedAt, &key.CreatedAt, &key.UpdatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, &key)
	}
	return keys, rows.Err()
}
