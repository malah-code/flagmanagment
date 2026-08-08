package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
)

type WebhookIntegrationRepository interface {
	Create(ctx context.Context, wh *models.WebhookIntegration) error
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.WebhookIntegration, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type webhookIntegrationRepository struct {
	db *sql.DB
}

func NewWebhookIntegrationRepository(db *sql.DB) WebhookIntegrationRepository {
	return &webhookIntegrationRepository{db: db}
}

func (r *webhookIntegrationRepository) Create(ctx context.Context, wh *models.WebhookIntegration) error {
	if wh.ID == uuid.Nil {
		wh.ID = uuid.New()
	}
	now := time.Now()
	wh.CreatedAt = now
	wh.UpdatedAt = now

	query := `INSERT INTO webhook_integrations (id, project_id, url, secret_key, events, is_active, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, wh.ID, wh.ProjectID, wh.URL, wh.SecretKey, wh.Events, wh.IsActive, wh.CreatedAt, wh.UpdatedAt)
	return err
}

func (r *webhookIntegrationRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.WebhookIntegration, error) {
	query := `SELECT id, project_id, url, secret_key, events, is_active, created_at, updated_at
	          FROM webhook_integrations WHERE project_id = $1 AND is_active = true`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.WebhookIntegration
	for rows.Next() {
		var wh models.WebhookIntegration
		if err := rows.Scan(&wh.ID, &wh.ProjectID, &wh.URL, &wh.SecretKey, &wh.Events, &wh.IsActive, &wh.CreatedAt, &wh.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &wh)
	}
	return list, nil
}

func (r *webhookIntegrationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM webhook_integrations WHERE id = $1`, id)
	return err
}
