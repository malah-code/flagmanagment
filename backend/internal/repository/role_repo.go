package repository

import (
	"context"
	"database/sql"
	"github.com/flagmanagment/backend/internal/models"

	"github.com/google/uuid"
)

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *models.Role) error {
	query := `INSERT INTO roles (id, name, description, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING created_at, updated_at`
	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}
	return r.db.QueryRowContext(ctx, query, role.ID, role.Name, role.Description).Scan(&role.CreatedAt, &role.UpdatedAt)
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1`
	var role models.Role
	err := r.db.QueryRowContext(ctx, query, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &role, err
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE name = $1`
	var role models.Role
	err := r.db.QueryRowContext(ctx, query, name).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &role, err
}

func (r *roleRepository) List(ctx context.Context) ([]*models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}
	return roles, nil
}

func (r *roleRepository) AssignUserRole(ctx context.Context, ur *models.UserRole) error {
	query := `INSERT INTO user_roles (id, user_id, role_id, project_id, environment_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING created_at, updated_at`
	if ur.ID == uuid.Nil {
		ur.ID = uuid.New()
	}
	return r.db.QueryRowContext(ctx, query, ur.ID, ur.UserID, ur.RoleID, ur.ProjectID, ur.EnvironmentID).Scan(&ur.CreatedAt, &ur.UpdatedAt)
}

func (r *roleRepository) RemoveUserRole(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *roleRepository) RemoveAllUserRoles(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID, projectID *uuid.UUID) ([]*models.UserRole, error) {
	query := `
		SELECT ur.id, ur.user_id, ur.role_id, ur.project_id, ur.environment_id, ur.created_at, ur.updated_at, r.name as role_name
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = $1 AND (ur.project_id = $2 OR ur.project_id IS NULL)
	`
	rows, err := r.db.QueryContext(ctx, query, userID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userRoles []*models.UserRole
	for rows.Next() {
		var ur models.UserRole
		var roleName string
		if err := rows.Scan(&ur.ID, &ur.UserID, &ur.RoleID, &ur.ProjectID, &ur.EnvironmentID, &ur.CreatedAt, &ur.UpdatedAt, &roleName); err != nil {
			return nil, err
		}
		ur.Role = &models.Role{Name: roleName}
		userRoles = append(userRoles, &ur)
	}
	return userRoles, nil
}

func (r *roleRepository) GetServiceAccountRoles(ctx context.Context, saID uuid.UUID, projectID *uuid.UUID) ([]*models.ServiceAccountRole, error) {
	query := `
		SELECT sar.id, sar.service_account_id, sar.role_id, sar.project_id, sar.environment_id, sar.created_at, sar.updated_at, ro.name as role_name
		FROM service_account_roles sar
		JOIN roles ro ON sar.role_id = ro.id
		WHERE sar.service_account_id = $1 AND (sar.project_id = $2 OR sar.project_id IS NULL)
	`
	rows, err := r.db.QueryContext(ctx, query, saID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var saRoles []*models.ServiceAccountRole
	for rows.Next() {
		var sar models.ServiceAccountRole
		var roleName string
		if err := rows.Scan(&sar.ID, &sar.ServiceAccountID, &sar.RoleID, &sar.ProjectID, &sar.EnvironmentID, &sar.CreatedAt, &sar.UpdatedAt, &roleName); err != nil {
			return nil, err
		}
		sar.Role = &models.Role{Name: roleName}
		saRoles = append(saRoles, &sar)
	}
	return saRoles, nil
}
