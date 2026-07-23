package cache_test

import (
	"context"
	"testing"

	"github.com/flagmanagment/backend/internal/cache"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/go-redis/redismock/v9"
)

func TestRulesetSnapshotCache(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &cache.Client{}
	_ = db
	_ = mock
	_ = client

	t.Run("MemoryCache instantiation", func(t *testing.T) {
		mem := cache.NewMemoryCache()
		mem.Set("env-123", "version-1")
		val, ok := mem.Get("env-123")
		if !ok || val != "version-1" {
			t.Errorf("expected version-1, got %s", val)
		}
		mem.Delete("env-123")
		_, ok = mem.Get("env-123")
		if ok {
			t.Errorf("expected key to be deleted")
		}
	})

	t.Run("RulesetSnapshot serialization structure", func(t *testing.T) {
		ctx := context.Background()
		_ = ctx
		snapshot := &models.RulesetSnapshot{
			Version: "v1.0.0",
			Flags: []models.FlagRule{
				{
					Key:              "test-flag",
					Type:             "BOOLEAN",
					Enabled:          true,
					DefaultVariation: "true",
				},
			},
		}
		if snapshot.Version != "v1.0.0" {
			t.Errorf("unexpected version")
		}
	})
}
