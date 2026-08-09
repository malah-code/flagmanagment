package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type Project struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (c *Client) CreateProject(ctx context.Context, p *Project) (*Project, error) {
	buf, err := c.doRequest(ctx, "POST", "/api/v1/projects", p)
	if err != nil {
		return nil, err
	}
	var created Project
	if err := json.Unmarshal(buf, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created project response: %w", err)
	}
	return &created, nil
}

func (c *Client) GetProject(ctx context.Context, id string) (*Project, error) {
	buf, err := c.doRequest(ctx, "GET", fmt.Sprintf("/api/v1/projects/%s", id), nil)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := json.Unmarshal(buf, &project); err != nil {
		return nil, fmt.Errorf("failed to parse project response: %w", err)
	}
	return &project, nil
}

func (c *Client) UpdateProject(ctx context.Context, p *Project) (*Project, error) {
	buf, err := c.doRequest(ctx, "PUT", fmt.Sprintf("/api/v1/projects/%s", p.ID), p)
	if err != nil {
		return nil, err
	}
	var updated Project
	if err := json.Unmarshal(buf, &updated); err != nil {
		return nil, fmt.Errorf("failed to parse updated project response: %w", err)
	}
	return &updated, nil
}

func (c *Client) DeleteProject(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/api/v1/projects/%s", id), nil)
	return err
}
