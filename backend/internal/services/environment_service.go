package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrProtectedEnvironment = errors.New("cannot delete protected environment")
)

type EnvironmentService struct {
	store repository.Store
	audit *AuditService
}

func NewEnvironmentService(store repository.Store, audit *AuditService) *EnvironmentService {
	return &EnvironmentService{
		store: store,
		audit: audit,
	}
}

func generateSecureAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateKeyAndHash() (apiKey, hashHex string, err error) {
	rawKey, err := generateSecureAPIKey()
	if err != nil {
		return "", "", err
	}
	apiKey = "env_" + rawKey
	hash := sha256.Sum256([]byte(apiKey))
	hashHex = hex.EncodeToString(hash[:])
	return apiKey, hashHex, nil
}

func generateRandomSalt() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *EnvironmentService) CreateEnvironment(ctx context.Context, projectID uuid.UUID, name string, isProtected bool, actorID uuid.UUID, actorIP string) (*models.Environment, string, error) {
	apiKey, hashHex, err := generateKeyAndHash()
	if err != nil {
		return nil, "", err
	}

	salt, err := generateRandomSalt()
	if err != nil {
		return nil, "", err
	}

	now := time.Now().UTC()
	env := &models.Environment{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        name,
		Key:         strings.ToLower(strings.ReplaceAll(name, " ", "-")),
		APIKeyHash:  hashHex,
		Salt:        salt,
		IsProtected: isProtected,
		CreatedAt:   now,
		UpdatedAt:   now,
		SDKSettings: models.JSONB{},
	}

	if err := s.store.EnvironmentRepo().Create(ctx, env); err != nil {
		return nil, "", err
	}

	// Auto-initialize flag states for this new environment
	flags, _, err := s.store.FlagRepo().ListByProject(ctx, projectID, 10000, 0)
	if err == nil {
		for _, f := range flags {
			state := &models.EnvironmentFlagState{
				ID:               uuid.New(),
				EnvironmentID:    env.ID,
				FeatureFlagID:    f.ID,
				Enabled:          false,
				TargetingRules:   models.JSONB{"rules": []interface{}{}},
				RemoteConfig:     models.JSONB{},
				CreatedAt:        now,
				UpdatedAt:        now,
			}
			_ = s.store.FlagStateRepo().Create(ctx, state)
		}
	}

	bNew, _ := json.Marshal(env)
	var newState models.JSONB
	json.Unmarshal(bNew, &newState)

	s.audit.LogAction(ctx, &models.AuditLog{
		ID:            uuid.New(),
		ProjectID:     &projectID,
		EnvironmentID: &env.ID,
		ActorID:       actorID,
		Action:        "CREATE",
		TargetType:    "ENVIRONMENT",
		TargetID:      env.ID,
		NewState:      newState,
		ActorIP:       actorIP,
		CreatedAt:     now,
	})

	return env, apiKey, nil
}


func (s *EnvironmentService) CloneEnvironment(ctx context.Context, projectID, sourceEnvID uuid.UUID, name string, actorID uuid.UUID, actorIP string) (*models.Environment, string, error) {
	sourceEnv, err := s.store.EnvironmentRepo().GetByID(ctx, sourceEnvID)
	if err != nil {
		return nil, "", err
	}
	if sourceEnv.ProjectID != projectID {
		return nil, "", repository.ErrNotFound
	}

	apiKey, hashHex, err := generateKeyAndHash()
	if err != nil {
		return nil, "", err
	}

	salt, err := generateRandomSalt()
	if err != nil {
		return nil, "", err
	}

	now := time.Now().UTC()
	newEnv := &models.Environment{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        name,
		Key:         strings.ToLower(strings.ReplaceAll(name, " ", "-")),
		APIKeyHash:  hashHex,
		Salt:        salt,
		IsProtected: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.store.WithTx(ctx, func(txStore repository.Store) error {
		if err := txStore.EnvironmentRepo().Create(ctx, newEnv); err != nil {
			return err
		}
		return txStore.FlagStateRepo().CloneEnvironmentState(ctx, sourceEnvID, newEnv.ID)
	})
	if err != nil {
		return nil, "", err
	}

	bNew, _ := json.Marshal(newEnv)
	var newState models.JSONB
	json.Unmarshal(bNew, &newState)

	prevState := models.JSONB{
		"source_environment_id": sourceEnvID.String(),
	}

	s.audit.LogAction(ctx, &models.AuditLog{
		ID:            uuid.New(),
		ProjectID:     &projectID,
		EnvironmentID: &newEnv.ID,
		ActorID:       actorID,
		Action:        "CLONE",
		TargetType:    "ENVIRONMENT",
		TargetID:      newEnv.ID,
		PreviousState: prevState,
		NewState:      newState,
		ActorIP:       actorIP,
		CreatedAt:     now,
	})

	return newEnv, apiKey, nil
}

func (s *EnvironmentService) DeleteEnvironment(ctx context.Context, envID uuid.UUID, actorID uuid.UUID, actorIP string) error {
	env, err := s.store.EnvironmentRepo().GetByID(ctx, envID)
	if err != nil {
		return err
	}
	if env.IsProtected {
		return ErrProtectedEnvironment
	}

	if err := s.store.EnvironmentRepo().Delete(ctx, envID); err != nil {
		return err
	}

	s.audit.LogAction(ctx, &models.AuditLog{
		ID:            uuid.New(),
		ProjectID:     &env.ProjectID,
		EnvironmentID: &env.ID,
		ActorID:       actorID,
		Action:        "DELETE",
		TargetType:    "ENVIRONMENT",
		TargetID:      envID,
		ActorIP:       actorIP,
		CreatedAt:     time.Now().UTC(),
	})

	return nil
}

func (s *EnvironmentService) CreateServerKey(ctx context.Context, envID uuid.UUID, name string, actorID uuid.UUID, actorIP string) (*models.EnvironmentServerKey, string, error) {
	env, err := s.store.EnvironmentRepo().GetByID(ctx, envID)
	if err != nil {
		return nil, "", err
	}

	apiKey, hashHex, err := generateKeyAndHash()
	if err != nil {
		return nil, "", err
	}

	now := time.Now().UTC()
	serverKey := &models.EnvironmentServerKey{
		ID:            uuid.New(),
		EnvironmentID: envID,
		Name:          name,
		KeyHash:       hashHex,
		CreatedAt:     now,
	}

	if err := s.store.EnvironmentRepo().CreateServerKey(ctx, serverKey); err != nil {
		return nil, "", err
	}

	s.audit.LogAction(ctx, &models.AuditLog{
		ID:            uuid.New(),
		ProjectID:     &env.ProjectID,
		EnvironmentID: &envID,
		ActorID:       actorID,
		Action:        "CREATE",
		TargetType:    "SERVER_KEY",
		TargetID:      serverKey.ID,
		NewState: models.JSONB{
			"name": name,
		},
		ActorIP:   actorIP,
		CreatedAt: now,
	})

	return serverKey, apiKey, nil
}

func (s *EnvironmentService) ListServerKeys(ctx context.Context, envID uuid.UUID) ([]*models.EnvironmentServerKey, error) {
	return s.store.EnvironmentRepo().ListServerKeys(ctx, envID)
}

func (s *EnvironmentService) DeleteServerKey(ctx context.Context, envID, keyID uuid.UUID, actorID uuid.UUID, actorIP string) error {
	env, err := s.store.EnvironmentRepo().GetByID(ctx, envID)
	if err != nil {
		return err
	}

	if err := s.store.EnvironmentRepo().DeleteServerKey(ctx, keyID); err != nil {
		return err
	}

	s.audit.LogAction(ctx, &models.AuditLog{
		ID:            uuid.New(),
		ProjectID:     &env.ProjectID,
		EnvironmentID: &envID,
		ActorID:       actorID,
		Action:        "DELETE",
		TargetType:    "SERVER_KEY",
		TargetID:      keyID,
		ActorIP:       actorIP,
		CreatedAt:     time.Now().UTC(),
	})

	return nil
}
