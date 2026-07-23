package repository

import (
	"context"
	"database/sql"
	"github.com/flagmanagment/backend/internal/models"

	"github.com/google/uuid"
)

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *models.Project) error {
	query := `INSERT INTO projects (id, name, key, description, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, project.ID, project.Name, project.Key, project.Description, project.CreatedAt, project.UpdatedAt)
	return err
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	query := `SELECT id, name, key, description, created_at, updated_at FROM projects WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var p models.Project
	if err := row.Scan(&p.ID, &p.Name, &p.Key, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *projectRepository) GetByKey(ctx context.Context, key string) (*models.Project, error) {
	query := `SELECT id, name, key, description, created_at, updated_at FROM projects WHERE key = $1`
	row := r.db.QueryRowContext(ctx, query, key)
	var p models.Project
	if err := row.Scan(&p.ID, &p.Name, &p.Key, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *projectRepository) List(ctx context.Context, limit, offset int) ([]*models.Project, int, error) {
	query := `SELECT id, name, key, description, created_at, updated_at FROM projects ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Key, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		projects = append(projects, &p)
	}

	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM projects`).Scan(&count); err != nil {
		return nil, 0, err
	}

	return projects, count, nil
}

func (r *projectRepository) Update(ctx context.Context, project *models.Project) error {
	query := `UPDATE projects SET name = $1, key = $2, description = $3, updated_at = $4 WHERE id = $5`
	res, err := r.db.ExecContext(ctx, query, project.Name, project.Key, project.Description, project.UpdatedAt, project.ID)
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

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
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
