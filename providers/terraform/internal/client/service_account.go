package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type ServiceAccount struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	ProjectID string `json:"project_id,omitempty"`
	Role      string `json:"role"`
	Token     string `json:"token,omitempty"`
}

func (c *Client) CreateServiceAccount(ctx context.Context, sa *ServiceAccount) (*ServiceAccount, error) {
	buf, err := c.doRequest(ctx, "POST", "/api/v1/service-accounts", sa)
	if err != nil {
		return nil, err
	}
	var created ServiceAccount
	if err := json.Unmarshal(buf, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created service account response: %w", err)
	}
	return &created, nil
}

func (c *Client) GetServiceAccount(ctx context.Context, id string) (*ServiceAccount, error) {
	buf, err := c.doRequest(ctx, "GET", fmt.Sprintf("/api/v1/service-accounts/%s", id), nil)
	if err != nil {
		return nil, err
	}
	var sa ServiceAccount
	if err := json.Unmarshal(buf, &sa); err != nil {
		return nil, fmt.Errorf("failed to parse service account response: %w", err)
	}
	return &sa, nil
}

func (c *Client) DeleteServiceAccount(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/api/v1/service-accounts/%s", id), nil)
	return err
}
