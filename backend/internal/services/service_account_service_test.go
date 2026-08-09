package services_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/services"
)

type mockSARepo struct {
	repository.ServiceAccountRepository
	keys map[string]*models.ServiceAccountKey
}

func (m *mockSARepo) CreateKey(ctx context.Context, key *models.ServiceAccountKey) error {
	if m.keys == nil {
		m.keys = make(map[string]*models.ServiceAccountKey)
	}
	m.keys[key.KeyHash] = key
	return nil
}

func (m *mockSARepo) GetKeyByHash(ctx context.Context, keyHash string) (*models.ServiceAccountKey, error) {
	if k, ok := m.keys[keyHash]; ok {
		return k, nil
	}
	return nil, repository.ErrNotFound
}

type mockSAStore struct {
	repository.Store
	saRepo repository.ServiceAccountRepository
}

func (m *mockSAStore) ServiceAccountRepo() repository.ServiceAccountRepository {
	return m.saRepo
}

func TestServiceAccountService_CreateAndValidateKey(t *testing.T) {
	saRepo := &mockSARepo{keys: make(map[string]*models.ServiceAccountKey)}
	store := &mockSAStore{saRepo: saRepo}
	svc := services.NewServiceAccountService(store)

	ctx := context.Background()
	saID := uuid.New()

	// 1. Create key
	key, plaintextKey, err := svc.CreateKey(ctx, saID, "Test Key", nil)
	if err != nil {
		t.Fatalf("expected no error creating key, got %v", err)
	}
	if !strings.HasPrefix(plaintextKey, "fm_sa_") {
		t.Fatalf("expected key prefix fm_sa_, got %s", plaintextKey)
	}
	if key.ServiceAccountID != saID {
		t.Fatalf("expected SA ID %s, got %s", saID, key.ServiceAccountID)
	}

	// 2. Validate valid key
	validated, err := svc.ValidateKey(ctx, plaintextKey)
	if err != nil {
		t.Fatalf("expected no error validating key, got %v", err)
	}
	if validated.ID != key.ID {
		t.Fatalf("expected key ID %s, got %s", key.ID, validated.ID)
	}

	// 3. Validate non-existent key
	_, err = svc.ValidateKey(ctx, "fm_sa_invalidtokenhere")
	if err == nil {
		t.Fatal("expected error validating non-existent key, got nil")
	}

	// 4. Validate expired key
	expiredTime := time.Now().Add(-1 * time.Hour)
	_, expiredPlaintext, err := svc.CreateKey(ctx, saID, "Expired Key", &expiredTime)
	if err != nil {
		t.Fatalf("failed to create expired key: %v", err)
	}
	_, err = svc.ValidateKey(ctx, expiredPlaintext)
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("expected key expired error, got %v", err)
	}
}
