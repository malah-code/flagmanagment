package services

import (
	"context"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
)

type AuditService struct {
	store repository.Store
}

func NewAuditService(store repository.Store) *AuditService {
	return &AuditService{store: store}
}

func (s *AuditService) LogAction(ctx context.Context, log *models.AuditLog) error {
	log.NewState = s.sanitize(log.NewState)
	log.PreviousState = s.sanitize(log.PreviousState)
	
	// Asynchronously log to prevent blocking hot paths
	go func(l models.AuditLog) {
		_ = s.store.AuditRepo().Create(context.Background(), &l)
	}(*log)

	return nil
}

func (s *AuditService) sanitize(data models.JSONB) models.JSONB {
	if data == nil {
		return data
	}

	sanitized := make(models.JSONB)
	for k, v := range data {
		sanitized[k] = v
	}

	// Redact sensitive keys
	sensitiveKeys := []string{"api_key", "password", "token", "secret"}
	for _, key := range sensitiveKeys {
		if _, exists := sanitized[key]; exists {
			sanitized[key] = "***"
		}
	}

	return sanitized
}
