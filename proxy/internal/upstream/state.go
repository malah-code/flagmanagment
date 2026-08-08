package upstream

import (
	"sync"
	"time"
)

type UpstreamStateSnapshot struct {
	Connected      bool       `json:"connected"`
	ConnectedSince *time.Time `json:"connected_since"`
	LastDeltaAt    *time.Time `json:"last_delta_at"`
}

type UpstreamState struct {
	mu             sync.RWMutex
	connected      bool
	connectedSince *time.Time
	lastDeltaAt    *time.Time
}

func NewUpstreamState() *UpstreamState {
	return &UpstreamState{}
}

func (s *UpstreamState) SetConnected(connected bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.connected = connected
	if connected {
		now := time.Now().UTC()
		s.connectedSince = &now
	}
}

func (s *UpstreamState) RecordDelta() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	s.lastDeltaAt = &now
}

func (s *UpstreamState) Snapshot() UpstreamStateSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return UpstreamStateSnapshot{
		Connected:      s.connected,
		ConnectedSince: s.connectedSince,
		LastDeltaAt:    s.lastDeltaAt,
	}
}
