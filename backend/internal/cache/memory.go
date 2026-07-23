package cache

import (
	"context"
	"sync"
)

// MemoryCache provides a fast thread-safe in-memory cache to shadow Redis for ultra low latency.
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]string),
	}
}

func (m *MemoryCache) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.items[key]
	return val, ok
}

func (m *MemoryCache) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[key] = value
}

func (m *MemoryCache) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, key)
}

// SyncWithRedis listens to Redis PubSub and invalidates the local memory cache.
func (m *MemoryCache) SyncWithRedis(ctx context.Context, client *Client, envID string) {
	pubsub := client.SubscribeRulesetUpdates(ctx, envID)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			m.Delete(envID)
			_ = msg
		}
	}
}
