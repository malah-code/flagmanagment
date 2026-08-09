package services

import (
	"context"
	"errors"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
)

type AuthService interface {
	HandleSSOLogin(ctx context.Context, provider string, email string, externalID string) (*models.User, error)
}

type authService struct {
	store repository.Store
}

func NewAuthService(store repository.Store) AuthService {
	return &authService{
		store: store,
	}
}

func (s *authService) HandleSSOLogin(ctx context.Context, provider string, email string, externalID string) (*models.User, error) {
	if provider != "oidc" && provider != "saml" {
		return nil, errors.New("invalid provider")
	}

	// Try to get by external ID first
	user, err := s.store.UserRepo().GetByExternalID(ctx, provider, externalID)
	if err == nil {
		// Found existing SSO user
		return user, nil
	}

	if err != repository.ErrNotFound {
		return nil, err
	}

	// If we didn't find them by external ID, check if they exist by email but as a local user.
	// If they do, we might want to link them, or return an error saying the email is already registered locally.
	// For security, if they exist locally, we will NOT link automatically in this MVP, we return an error.
	existingUser, err := s.store.UserRepo().GetByEmail(ctx, email)
	if err == nil {
		if existingUser.AuthProvider == "local" {
			return nil, errors.New("user already exists as local user. please login with password")
		}
		// If they exist as another SSO provider, prevent login
		return nil, errors.New("user exists with different auth provider")
	}

	if err != repository.ErrNotFound {
		return nil, err
	}

	// JIT Provisioning (Task T010)
	newUser := &models.User{
		Email:        email,
		AuthProvider: provider,
		ExternalID:   &externalID,
		PasswordHash: nil, // No password for SSO
	}

	err = s.store.UserRepo().Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// In the future, we would assign a default role here.
	// For the MVP, we assume any created user has access.
	// Or we can add them to a default role if RoleRepo is available.
	if roleRepo := s.store.RoleRepo(); roleRepo != nil {
		defaultRole, err := roleRepo.GetByName(ctx, "Viewer")
		if err == nil && defaultRole != nil {
			_ = roleRepo.AssignUserRole(ctx, &models.UserRole{
				UserID: newUser.ID,
				RoleID: defaultRole.ID,
			})
		}
	}

	return newUser, nil
}
