package sdk

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Client manages the SSE connection to the FlagManagment backend and
// maintains a thread-safe in-memory cache of flag definitions.
type Client struct {
	apiKey    string
	streamURL string
	flags     map[string]interface{}
	mu        sync.RWMutex
	cancel    context.CancelFunc
	done      chan struct{}
}

func NewClient(apiKey string, streamURL string) *Client {
	return &Client{
		apiKey:    apiKey,
		streamURL: streamURL,
		flags:     make(map[string]interface{}),
		done:      make(chan struct{}),
	}
}

// Connect starts the background SSE streaming goroutine.
func (c *Client) Connect() {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	go c.stream(ctx)
}

// Shutdown gracefully stops the background streaming goroutine.
func (c *Client) Shutdown() {
	if c.cancel != nil {
		c.cancel()
		<-c.done // wait for goroutine to exit
	}
}

func (c *Client) stream(ctx context.Context) {
	defer close(c.done)

	attempt := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		req, err := http.NewRequestWithContext(ctx, "GET", c.streamURL, nil)
		if err != nil {
			log.Printf("[flagmanagment-sdk] failed to create request: %v", err)
			attempt++
			c.backoff(ctx, attempt)
			continue
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "text/event-stream")

		client := &http.Client{Timeout: 0}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[flagmanagment-sdk] connection failed: %v", err)
			attempt++
			c.backoff(ctx, attempt)
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			log.Printf("[flagmanagment-sdk] unexpected status %d", resp.StatusCode)
			attempt++
			c.backoff(ctx, attempt)
			continue
		}

		// Connection established — reset backoff counter
		attempt = 0
		log.Printf("[flagmanagment-sdk] connected to %s", c.streamURL)

		// Read SSE events, tracking the current event type
		reader := bufio.NewScanner(resp.Body)
		currentEvent := ""
		for reader.Scan() {
			line := reader.Text()

			if strings.HasPrefix(line, "event:") {
				currentEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			} else if strings.HasPrefix(line, "data:") {
				dataStr := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				c.handleEvent(currentEvent, dataStr)
			} else if line == "" {
				// End of an SSE message block — reset event type
				currentEvent = ""
			}
		}

		resp.Body.Close()
		log.Printf("[flagmanagment-sdk] connection lost, reconnecting...")
		attempt++
		c.backoff(ctx, attempt)
	}
}

// handleEvent processes a single SSE event based on its type.
func (c *Client) handleEvent(eventType string, dataStr string) {
	switch eventType {
	case "bootstrap":
		var payload struct {
			Flags map[string]interface{} `json:"flags"`
		}
		if err := json.Unmarshal([]byte(dataStr), &payload); err != nil {
			log.Printf("[flagmanagment-sdk] failed to parse bootstrap: %v", err)
			return
		}
		if payload.Flags != nil {
			c.mu.Lock()
			c.flags = payload.Flags
			c.mu.Unlock()
			log.Printf("[flagmanagment-sdk] bootstrapped %d flags", len(payload.Flags))
		}

	case "flag_updated":
		var payload struct {
			FlagKey string      `json:"flagKey"`
			Flag    interface{} `json:"flag"`
		}
		if err := json.Unmarshal([]byte(dataStr), &payload); err != nil {
			log.Printf("[flagmanagment-sdk] failed to parse flag_updated: %v", err)
			return
		}
		if payload.FlagKey != "" && payload.Flag != nil {
			c.mu.Lock()
			c.flags[payload.FlagKey] = payload.Flag
			c.mu.Unlock()
			log.Printf("[flagmanagment-sdk] updated flag: %s", payload.FlagKey)
		}

	case "ping":
		// Heartbeat — no action needed
	default:
		// Unknown event type — ignore gracefully
	}
}

// backoff waits with exponential backoff (capped at 60s).
func (c *Client) backoff(ctx context.Context, attempt int) {
	delay := time.Duration(math.Min(float64(attempt*attempt), 60)) * time.Second
	if delay < 1*time.Second {
		delay = 1 * time.Second
	}
	select {
	case <-ctx.Done():
	case <-time.After(delay):
	}
}

// GetFlag returns a flag definition by key. Thread-safe.
func (c *Client) GetFlag(flagKey string) (map[string]interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	flagObj, exists := c.flags[flagKey]
	if !exists {
		return nil, fmt.Errorf("flag not found: %s", flagKey)
	}

	flag, ok := flagObj.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid flag format for: %s", flagKey)
	}

	return flag, nil
}
