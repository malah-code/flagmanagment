package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
)

type SlackConfigRepository interface {
	GetByEnvironmentID(ctx context.Context, envID uuid.UUID) (*models.SlackWebhookConfig, error)
	Upsert(ctx context.Context, config *models.SlackWebhookConfig) error
	Delete(ctx context.Context, envID uuid.UUID) error
}

type slackConfigRepository struct {
	db *sql.DB
}

func NewSlackConfigRepository(db *sql.DB) SlackConfigRepository {
	return &slackConfigRepository{db: db}
}

func (r *slackConfigRepository) GetByEnvironmentID(ctx context.Context, envID uuid.UUID) (*models.SlackWebhookConfig, error) {
	query := `
		SELECT id, environment_id, webhook_url, enabled, created_at, updated_at
		FROM slack_webhook_configs
		WHERE environment_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, envID)

	var c models.SlackWebhookConfig
	err := row.Scan(&c.ID, &c.EnvironmentID, &c.WebhookURL, &c.Enabled, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *slackConfigRepository) Upsert(ctx context.Context, config *models.SlackWebhookConfig) error {
	query := `
		INSERT INTO slack_webhook_configs (id, environment_id, webhook_url, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (environment_id) DO UPDATE SET
			webhook_url = EXCLUDED.webhook_url,
			enabled = EXCLUDED.enabled,
			updated_at = EXCLUDED.updated_at
	`
	now := time.Now().UTC()
	if config.ID == uuid.Nil {
		config.ID = uuid.New()
	}
	if config.CreatedAt.IsZero() {
		config.CreatedAt = now
	}
	config.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query, config.ID, config.EnvironmentID, config.WebhookURL, config.Enabled, config.CreatedAt, config.UpdatedAt)
	return err
}

func (r *slackConfigRepository) Delete(ctx context.Context, envID uuid.UUID) error {
	query := `DELETE FROM slack_webhook_configs WHERE environment_id = $1`
	_, err := r.db.ExecContext(ctx, query, envID)
	return err
}
