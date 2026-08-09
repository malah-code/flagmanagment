package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type Environment struct {
	ID          string `json:"id,omitempty"`
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	IsProtected bool   `json:"is_protected"`
}

func (c *Client) CreateEnvironment(ctx context.Context, env *Environment) (*Environment, error) {
	buf, err := c.doRequest(ctx, "POST", fmt.Sprintf("/api/v1/projects/%s/environments", env.ProjectID), env)
	if err != nil {
		return nil, err
	}
	var created Environment
	if err := json.Unmarshal(buf, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created environment response: %w", err)
	}
	return &created, nil
}

func (c *Client) GetEnvironment(ctx context.Context, id string) (*Environment, error) {
	buf, err := c.doRequest(ctx, "GET", fmt.Sprintf("/api/v1/environments/%s", id), nil)
	if err != nil {
		return nil, err
	}
	var env Environment
	if err := json.Unmarshal(buf, &env); err != nil {
		return nil, fmt.Errorf("failed to parse environment response: %w", err)
	}
	return &env, nil
}

func (c *Client) UpdateEnvironment(ctx context.Context, env *Environment) (*Environment, error) {
	buf, err := c.doRequest(ctx, "PUT", fmt.Sprintf("/api/v1/environments/%s", env.ID), env)
	if err != nil {
		return nil, err
	}
	var updated Environment
	if err := json.Unmarshal(buf, &updated); err != nil {
		return nil, fmt.Errorf("failed to parse updated environment response: %w", err)
	}
	return &updated, nil
}

func (c *Client) DeleteEnvironment(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/api/v1/environments/%s", id), nil)
	return err
}
