package repository

import (
	"context"
	"database/sql"
	"github.com/flagmanagment/backend/internal/models"
)

type systemConfigRepository struct {
	db *sql.DB
}

func NewSystemConfigRepository(db *sql.DB) SystemConfigRepository {
	return &systemConfigRepository{db: db}
}

func (r *systemConfigRepository) GetByKey(ctx context.Context, key string) (*models.SystemConfig, error) {
	query := `SELECT key, value, updated_at FROM system_configs WHERE key = $1`
	var config models.SystemConfig
	err := r.db.QueryRowContext(ctx, query, key).Scan(&config.Key, &config.Value, &config.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &config, err
}

func (r *systemConfigRepository) Upsert(ctx context.Context, config *models.SystemConfig) error {
	query := `
		INSERT INTO system_configs (key, value, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = NOW()
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query, config.Key, config.Value).Scan(&config.UpdatedAt)
}
