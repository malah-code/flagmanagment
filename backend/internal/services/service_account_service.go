package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
)

type ServiceAccountService interface {
	CreateKey(ctx context.Context, saID uuid.UUID, name string, expiresAt *time.Time) (*models.ServiceAccountKey, string, error)
	ValidateKey(ctx context.Context, plaintextKey string) (*models.ServiceAccountKey, error)
}

type serviceAccountService struct {
	store repository.Store
}

func NewServiceAccountService(store repository.Store) ServiceAccountService {
	return &serviceAccountService{store: store}
}

func (s *serviceAccountService) CreateKey(ctx context.Context, saID uuid.UUID, name string, expiresAt *time.Time) (*models.ServiceAccountKey, string, error) {
	// Generate random 32 byte key
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, "", err
	}
	plaintextKey := "fm_sa_" + base64.RawURLEncoding.EncodeToString(b)

	// Hash the key for storage
	hash := sha256.Sum256([]byte(plaintextKey))
	keyHash := hex.EncodeToString(hash[:])

	key := &models.ServiceAccountKey{
		ID:               uuid.New(),
		ServiceAccountID: saID,
		Name:             name,
		KeyHash:          keyHash,
		ExpiresAt:        expiresAt,
	}

	if err := s.store.ServiceAccountRepo().CreateKey(ctx, key); err != nil {
		return nil, "", err
	}

	return key, plaintextKey, nil
}

func (s *serviceAccountService) ValidateKey(ctx context.Context, plaintextKey string) (*models.ServiceAccountKey, error) {
	hash := sha256.Sum256([]byte(plaintextKey))
	keyHash := hex.EncodeToString(hash[:])

	key, err := s.store.ServiceAccountRepo().GetKeyByHash(ctx, keyHash)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, errors.New("invalid key")
		}
		return nil, err
	}

	if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("key expired")
	}

	return key, nil
}
