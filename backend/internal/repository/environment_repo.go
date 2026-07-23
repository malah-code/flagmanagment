package repository

import (
	"context"
	"database/sql"
	"github.com/flagmanagment/backend/internal/models"

	"github.com/google/uuid"
)

type environmentRepository struct {
	db *sql.DB
}

func NewEnvironmentRepository(db *sql.DB) EnvironmentRepository {
	return &environmentRepository{db: db}
}

func (r *environmentRepository) Create(ctx context.Context, env *models.Environment) error {
	query := `INSERT INTO environments (id, project_id, name, key, api_key_hash, is_protected, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, env.ID, env.ProjectID, env.Name, env.Key, env.APIKeyHash, env.IsProtected, env.CreatedAt, env.UpdatedAt)
	return err
}

func (r *environmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Environment, error) {
	query := `SELECT id, project_id, name, key, api_key_hash, is_protected, created_at, updated_at FROM environments WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var e models.Environment
	if err := row.Scan(&e.ID, &e.ProjectID, &e.Name, &e.Key, &e.APIKeyHash, &e.IsProtected, &e.CreatedAt, &e.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *environmentRepository) GetByAPIKeyHash(ctx context.Context, apiKeyHash string) (*models.Environment, error) {
	query := `SELECT id, project_id, name, key, api_key_hash, is_protected, created_at, updated_at FROM environments WHERE api_key_hash = $1`
	row := r.db.QueryRowContext(ctx, query, apiKeyHash)
	var e models.Environment
	if err := row.Scan(&e.ID, &e.ProjectID, &e.Name, &e.Key, &e.APIKeyHash, &e.IsProtected, &e.CreatedAt, &e.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *environmentRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.Environment, error) {
	query := `SELECT id, project_id, name, key, api_key_hash, is_protected, created_at, updated_at FROM environments WHERE project_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envs []*models.Environment
	for rows.Next() {
		var e models.Environment
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.Name, &e.Key, &e.APIKeyHash, &e.IsProtected, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		envs = append(envs, &e)
	}
	return envs, nil
}

func (r *environmentRepository) Update(ctx context.Context, env *models.Environment) error {
	query := `UPDATE environments SET name = $1, key = $2, is_protected = $3, updated_at = $4 WHERE id = $5`
	res, err := r.db.ExecContext(ctx, query, env.Name, env.Key, env.IsProtected, env.UpdatedAt, env.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *environmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM environments WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
