package services

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
)

type AuditService struct {
	store       repository.Store
	subscribers map[chan *models.AuditLog]bool
	mu          sync.RWMutex
}

func NewAuditService(store repository.Store) *AuditService {
	return &AuditService{
		store:       store,
		subscribers: make(map[chan *models.AuditLog]bool),
	}
}

func (s *AuditService) LogAction(ctx context.Context, log *models.AuditLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	log.NewState = ScrubJSONB(log.NewState)
	log.PreviousState = ScrubJSONB(log.PreviousState)
	
	// Asynchronously log to prevent blocking hot paths
	go func(l models.AuditLog) {
		_ = s.store.AuditRepo().Create(context.Background(), &l)
	}(*log)

	s.mu.RLock()
	defer s.mu.RUnlock()
	for ch := range s.subscribers {
		select {
		case ch <- log:
		default:
		}
	}

	return nil
}

func (s *AuditService) Subscribe() <-chan *models.AuditLog {
	ch := make(chan *models.AuditLog, 100)
	s.mu.Lock()
	s.subscribers[ch] = true
	s.mu.Unlock()
	return ch
}

func (s *AuditService) CleanupOldLogs(ctx context.Context, days int) error {
	cutoff := time.Now().UTC().AddDate(0, 0, -days)
	return s.store.AuditRepo().DeleteOlderThan(ctx, cutoff)
}

var sensitiveKeys = map[string]bool{
	"api_key":       true,
	"apikey":        true,
	"token":         true,
	"secret":        true,
	"secret_key":    true,
	"password":      true,
	"authorization": true,
}

// ScrubJSONB recursively redacts sensitive fields in models.JSONB.
func ScrubJSONB(j models.JSONB) models.JSONB {
	if j == nil {
		return nil
	}
	res := scrubValue(j)
	if m, ok := res.(map[string]interface{}); ok {
		return models.JSONB(m)
	}
	return j
}

func scrubValue(v interface{}) interface{} {
	switch val := v.(type) {
	case models.JSONB:
		res := make(map[string]interface{}, len(val))
		for k, vChild := range val {
			kLower := strings.ToLower(k)
			if sensitiveKeys[kLower] {
				res[k] = "[REDACTED]"
			} else {
				res[k] = scrubValue(vChild)
			}
		}
		return res
	case map[string]interface{}:
		res := make(map[string]interface{}, len(val))
		for k, vChild := range val {
			kLower := strings.ToLower(k)
			if sensitiveKeys[kLower] {
				res[k] = "[REDACTED]"
			} else {
				res[k] = scrubValue(vChild)
			}
		}
		return res
	case []interface{}:
		res := make([]interface{}, len(val))
		for i, item := range val {
			res[i] = scrubValue(item)
		}
		return res
	default:
		return val
	}
}
