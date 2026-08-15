package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
)

type UserService interface {
	GetUsers(ctx context.Context, limit, offset int) ([]*UserResponse, int, error)
	CreateInvitation(ctx context.Context, email, role string, projectIDs []string, createdBy uuid.UUID) (*models.Invitation, error)
	UpdateUserAccess(ctx context.Context, userID uuid.UUID, roleName string, projectIDs []string) error
}

type userService struct {
	store         repository.Store
	cryptoService CryptoService
	emailService  EmailService
}

func NewUserService(store repository.Store, crypto CryptoService, email EmailService) UserService {
	return &userService{
		store:         store,
		cryptoService: crypto,
		emailService:  email,
	}
}

// UserResponse represents the user data returned by the API
type UserResponse struct {
	*models.User
	Roles    []string `json:"roles"`
	Projects []string `json:"projects"`
}

func (s *userService) GetUsers(ctx context.Context, limit, offset int) ([]*UserResponse, int, error) {
	users, total, err := s.store.UserRepo().List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var response []*UserResponse
	for _, user := range users {
		userRoles, err := s.store.RoleRepo().GetUserRoles(ctx, user.ID, nil)
		if err != nil {
			return nil, 0, err
		}

		rolesMap := make(map[string]bool)
		projectsMap := make(map[string]bool)
		
		for _, ur := range userRoles {
			if ur.Role != nil {
				rolesMap[ur.Role.Name] = true
			}
			if ur.ProjectID != nil {
				projectsMap[ur.ProjectID.String()] = true
			}
		}

		roles := []string{}
		for r := range rolesMap {
			roles = append(roles, r)
		}

		projects := []string{}
		for p := range projectsMap {
			projects = append(projects, p)
		}

		response = append(response, &UserResponse{
			User:     user,
			Roles:    roles,
			Projects: projects,
		})
	}

	return response, total, nil
}

func (s *userService) CreateInvitation(ctx context.Context, email, role string, projectIDs []string, createdBy uuid.UUID) (*models.Invitation, error) {
	if _, err := s.store.UserRepo().GetByEmail(ctx, email); err == nil {
		return nil, errors.New("user already exists")
	}

	token, err := s.cryptoService.GenerateToken()
	if err != nil {
		return nil, err
	}

	tokenHash := s.cryptoService.HashToken(token)

	inv := &models.Invitation{
		Email:      email,
		TokenHash:  tokenHash,
		Role:       role,
		ProjectIDs: models.JSONB{"ids": projectIDs},
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour),
		CreatedBy:  &createdBy,
	}

	if err := s.store.InvitationRepo().Create(ctx, inv); err != nil {
		return nil, err
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	inviteLink := fmt.Sprintf("%s/accept-invite?token=%s", frontendURL, token)

	if err := s.emailService.SendInvitation(email, inviteLink); err != nil {
		return nil, fmt.Errorf("failed to send invitation email: %v", err)
	}

	return inv, nil
}

func normalizeRoleName(roleName string) string {
	switch roleName {
	case "Global Administrator", "ADMIN":
		return "ADMIN"
	case "Project Editor", "EDITOR":
		return "EDITOR"
	case "Read-Only Auditor", "VIEWER":
		return "VIEWER"
	default:
		return roleName
	}
}

func (s *userService) UpdateUserAccess(ctx context.Context, userID uuid.UUID, roleName string, projectIDs []string) error {
	if _, err := s.store.UserRepo().GetByID(ctx, userID); err != nil {
		return err
	}

	normalizedRole := normalizeRoleName(roleName)
	role, err := s.store.RoleRepo().GetByName(ctx, normalizedRole)
	if err != nil {
		return err
	}

	return s.store.WithTx(ctx, func(store repository.Store) error {
		if err := store.RoleRepo().RemoveAllUserRoles(ctx, userID); err != nil {
			return err
		}

		if normalizedRole == "ADMIN" {
			return store.RoleRepo().AssignUserRole(ctx, &models.UserRole{
				UserID: userID,
				RoleID: role.ID,
			})
		}

		if len(projectIDs) == 0 {
			return store.RoleRepo().AssignUserRole(ctx, &models.UserRole{
				UserID: userID,
				RoleID: role.ID,
			})
		}

		for _, pIDStr := range projectIDs {
			pID, err := uuid.Parse(pIDStr)
			if err != nil {
				continue
			}
			if err := store.RoleRepo().AssignUserRole(ctx, &models.UserRole{
				UserID:    userID,
				RoleID:    role.ID,
				ProjectID: &pID,
			}); err != nil {
				return err
			}
		}

		return nil
	})
}
