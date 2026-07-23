package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// PublishRulesetUpdate notifies connected instances of a ruleset delta.
func (c *Client) PublishRulesetUpdate(ctx context.Context, envID, payload string) error {
	channel := fmt.Sprintf("ruleset:updates:%s", envID)
	return c.rdb.Publish(ctx, channel, payload).Err()
}

// SubscribeRulesetUpdates returns a Redis PubSub channel for an environment.
func (c *Client) SubscribeRulesetUpdates(ctx context.Context, envID string) *redis.PubSub {
	channel := fmt.Sprintf("ruleset:updates:%s", envID)
	return c.rdb.Subscribe(ctx, channel)
}
