package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type FeatureFlag struct {
	ID             string `json:"id,omitempty"`
	ProjectID      string `json:"project_id"`
	Key            string `json:"key"`
	Name           string `json:"name"`
	Type           string `json:"type"` // boolean or multivariate
	ParentFlagID   string `json:"parent_flag_id,omitempty"`
}

func (c *Client) CreateFeatureFlag(ctx context.Context, flag *FeatureFlag) (*FeatureFlag, error) {
	buf, err := c.doRequest(ctx, "POST", fmt.Sprintf("/api/v1/projects/%s/flags", flag.ProjectID), flag)
	if err != nil {
		return nil, err
	}
	var created FeatureFlag
	if err := json.Unmarshal(buf, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created flag response: %w", err)
	}
	return &created, nil
}

func (c *Client) GetFeatureFlag(ctx context.Context, id string) (*FeatureFlag, error) {
	buf, err := c.doRequest(ctx, "GET", fmt.Sprintf("/api/v1/flags/%s", id), nil)
	if err != nil {
		return nil, err
	}
	var flag FeatureFlag
	if err := json.Unmarshal(buf, &flag); err != nil {
		return nil, fmt.Errorf("failed to parse flag response: %w", err)
	}
	return &flag, nil
}

func (c *Client) UpdateFeatureFlag(ctx context.Context, flag *FeatureFlag) (*FeatureFlag, error) {
	buf, err := c.doRequest(ctx, "PUT", fmt.Sprintf("/api/v1/flags/%s", flag.ID), flag)
	if err != nil {
		return nil, err
	}
	var updated FeatureFlag
	if err := json.Unmarshal(buf, &updated); err != nil {
		return nil, fmt.Errorf("failed to parse updated flag response: %w", err)
	}
	return &updated, nil
}

func (c *Client) DeleteFeatureFlag(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/api/v1/flags/%s", id), nil)
	return err
}
