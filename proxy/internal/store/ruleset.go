package store

import (
	"sync"

	pb "github.com/flagmanagment/proxy/pkg/gen/sdk/v1"
)

type RulesetStore struct {
	mu       sync.RWMutex
	snapshot *pb.RulesetSnapshot
	version  string
}

func NewRulesetStore() *RulesetStore {
	return &RulesetStore{}
}

func (s *RulesetStore) Get() *pb.RulesetSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshot
}

func (s *RulesetStore) Set(snapshot *pb.RulesetSnapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshot = snapshot
	if snapshot != nil {
		s.version = snapshot.Version
	} else {
		s.version = ""
	}
}

func (s *RulesetStore) Version() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.version
}
