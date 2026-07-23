package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/flagmanagment/backend/internal/models"
)

// GetRulesetSnapshot retrieves a ruleset snapshot for an environment from Redis.
func (c *Client) GetRulesetSnapshot(ctx context.Context, envID string) (*models.RulesetSnapshot, error) {
	key := fmt.Sprintf("ruleset:%s", envID)
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var snapshot models.RulesetSnapshot
	if err := json.Unmarshal([]byte(val), &snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}

// SetRulesetSnapshot caches a ruleset snapshot for an environment in Redis.
func (c *Client) SetRulesetSnapshot(ctx context.Context, envID string, snapshot *models.RulesetSnapshot) error {
	key := fmt.Sprintf("ruleset:%s", envID)
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, data, 0).Err()
}
