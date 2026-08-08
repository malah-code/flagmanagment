package services

import (
	"context"
	"sync"
	"time"

	"github.com/flagmanagment/backend/internal/repository"
	"github.com/google/uuid"
)

type MetricAggregationService interface {
	RecordEvaluation(flagStateID uuid.UUID, evaluatedAt time.Time)
	Start(ctx context.Context, flushInterval time.Duration)
	Flush(ctx context.Context) error
}

type metricAggregationService struct {
	repo repository.FlagStateRepository
	mu   sync.Mutex
	hits map[uuid.UUID]time.Time
}

func NewMetricAggregationService(repo repository.FlagStateRepository) MetricAggregationService {
	return &metricAggregationService{
		repo: repo,
		hits: make(map[uuid.UUID]time.Time),
	}
}

func (s *metricAggregationService) RecordEvaluation(flagStateID uuid.UUID, evaluatedAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.hits[flagStateID]; !ok || evaluatedAt.After(existing) {
		s.hits[flagStateID] = evaluatedAt
	}
}

func (s *metricAggregationService) Start(ctx context.Context, flushInterval time.Duration) {
	ticker := time.NewTicker(flushInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				_ = s.Flush(context.Background())
				return
			case <-ticker.C:
				_ = s.Flush(ctx)
			}
		}
	}()
}

func (s *metricAggregationService) Flush(ctx context.Context) error {
	s.mu.Lock()
	if len(s.hits) == 0 {
		s.mu.Unlock()
		return nil
	}
	batch := s.hits
	s.hits = make(map[uuid.UUID]time.Time)
	s.mu.Unlock()

	return s.repo.UpdateLastEvaluatedAtBatch(ctx, batch)
}
