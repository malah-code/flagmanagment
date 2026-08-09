package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
)

type mockUserRepo struct {
	repository.UserRepository
	usersByExt map[string]*models.User
	usersByEmail map[string]*models.User
}

func (m *mockUserRepo) Create(ctx context.Context, u *models.User) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if m.usersByEmail == nil {
		m.usersByEmail = make(map[string]*models.User)
	}
	m.usersByEmail[u.Email] = u

	if u.AuthProvider != "" && u.ExternalID != nil {
		if m.usersByExt == nil {
			m.usersByExt = make(map[string]*models.User)
		}
		m.usersByExt[u.AuthProvider+":"+*u.ExternalID] = u
	}
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if u, ok := m.usersByEmail[email]; ok {
		return u, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockUserRepo) GetByExternalID(ctx context.Context, provider, externalID string) (*models.User, error) {
	if u, ok := m.usersByExt[provider+":"+externalID]; ok {
		return u, nil
	}
	return nil, repository.ErrNotFound
}

type mockAuthStore struct {
	repository.Store
	userRepo repository.UserRepository
}

func (m *mockAuthStore) UserRepo() repository.UserRepository {
	return m.userRepo
}

func (m *mockAuthStore) RoleRepo() repository.RoleRepository {
	return nil
}

func TestAuthService_HandleSSOLogin_JITProvisioning(t *testing.T) {
	userRepo := &mockUserRepo{
		usersByExt:   make(map[string]*models.User),
		usersByEmail: make(map[string]*models.User),
	}
	store := &mockAuthStore{userRepo: userRepo}
	svc := services.NewAuthService(store)

	ctx := context.Background()

	// 1. Provision new SSO user
	user, err := svc.HandleSSOLogin(ctx, "oidc", "sso_user@example.com", "sub_12345")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil || user.Email != "sso_user@example.com" {
		t.Fatalf("expected email sso_user@example.com, got %v", user)
	}
	if user.AuthProvider != "oidc" {
		t.Fatalf("expected auth provider oidc, got %s", user.AuthProvider)
	}

	// 2. Login again with same external ID -> returns existing user
	user2, err := svc.HandleSSOLogin(ctx, "oidc", "sso_user@example.com", "sub_12345")
	if err != nil {
		t.Fatalf("expected no error on subsequent login, got %v", err)
	}
	if user2.ID != user.ID {
		t.Fatalf("expected same user ID, got %s vs %s", user.ID, user2.ID)
	}

	// 3. Attempt login with invalid provider
	_, err = svc.HandleSSOLogin(ctx, "invalid", "sso_user@example.com", "sub_12345")
	if err == nil {
		t.Fatal("expected error for invalid provider, got nil")
	}
}

func TestAuthService_HandleSSOLogin_Conflict(t *testing.T) {
	userRepo := &mockUserRepo{
		usersByExt:   make(map[string]*models.User),
		usersByEmail: make(map[string]*models.User),
	}
	// Existing local user
	userRepo.usersByEmail["local@example.com"] = &models.User{
		ID:           uuid.New(),
		Email:        "local@example.com",
		AuthProvider: "local",
	}

	store := &mockAuthStore{userRepo: userRepo}
	svc := services.NewAuthService(store)

	ctx := context.Background()

	// Try to login via SSO with email belonging to local user
	_, err := svc.HandleSSOLogin(ctx, "oidc", "local@example.com", "sub_999")
	if err == nil {
		t.Fatal("expected conflict error when SSO email belongs to local user, got nil")
	}
}
