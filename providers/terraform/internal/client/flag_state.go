package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type Variant struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Condition struct {
	Attribute string   `json:"attribute"`
	Operator  string   `json:"operator"`
	Values    []string `json:"values"`
}

type Rollout struct {
	Percentages map[string]int `json:"percentages"`
}

type TargetingRule struct {
	Name       string      `json:"name"`
	Variant    string      `json:"variant,omitempty"`
	Rollout    *Rollout    `json:"rollout,omitempty"`
	Conditions []Condition `json:"conditions"`
}

type TelemetryTrigger struct {
	MetricName string  `json:"metric_name"`
	Threshold  float64 `json:"threshold"`
	Action     string  `json:"action"` // e.g., KILL_SWITCH
}

type FlagState struct {
	EnvironmentID    string            `json:"environment_id"`
	FlagID           string            `json:"flag_id"`
	Enabled          bool              `json:"enabled"`
	DefaultVariant   string            `json:"default_variant"`
	Variants         []Variant         `json:"variants,omitempty"`
	TargetingRules   []TargetingRule   `json:"targeting_rules,omitempty"`
	TelemetryTrigger *TelemetryTrigger `json:"telemetry_trigger,omitempty"`
}

func (c *Client) GetFlagState(ctx context.Context, envID, flagID string) (*FlagState, error) {
	buf, err := c.doRequest(ctx, "GET", fmt.Sprintf("/api/v1/environments/%s/flags/%s", envID, flagID), nil)
	if err != nil {
		return nil, err
	}
	var fs FlagState
	if err := json.Unmarshal(buf, &fs); err != nil {
		return nil, fmt.Errorf("failed to parse flag state response: %w", err)
	}
	return &fs, nil
}

func (c *Client) UpdateFlagState(ctx context.Context, fs *FlagState) (*FlagState, error) {
	buf, err := c.doRequest(ctx, "PUT", fmt.Sprintf("/api/v1/environments/%s/flags/%s", fs.EnvironmentID, fs.FlagID), fs)
	if err != nil {
		return nil, err
	}
	var updated FlagState
	if err := json.Unmarshal(buf, &updated); err != nil {
		return nil, fmt.Errorf("failed to parse updated flag state response: %w", err)
	}
	return &updated, nil
}

func (c *Client) DeleteFlagState(ctx context.Context, envID, flagID string) error {
	_, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/api/v1/environments/%s/flags/%s", envID, flagID), nil)
	return err
}
