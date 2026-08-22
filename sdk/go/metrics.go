package sdk

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type SDKEvaluationMetric struct {
	FlagKey   string `json:"flagKey"`
	Timestamp string `json:"timestamp"`
}

type MetricsSync struct {
	apiKey      string
	endpointURL string
	buffer      []SDKEvaluationMetric
	mu          sync.Mutex
	done        chan struct{}
}

func NewMetricsSync(apiKey string, endpointURL string) *MetricsSync {
	return &MetricsSync{
		apiKey:      apiKey,
		endpointURL: endpointURL,
		buffer:      make([]SDKEvaluationMetric, 0),
		done:        make(chan struct{}),
	}
}

func (m *MetricsSync) Start(interval time.Duration) {
	if interval == 0 {
		interval = 15 * time.Second
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-m.done:
				m.flush()
				return
			case <-ticker.C:
				m.flush()
			}
		}
	}()
}

func (m *MetricsSync) Stop() {
	close(m.done)
}

func (m *MetricsSync) RecordEvaluation(flagKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.buffer = append(m.buffer, SDKEvaluationMetric{
		FlagKey:   flagKey,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func (m *MetricsSync) flush() {
	m.mu.Lock()
	if len(m.buffer) == 0 {
		m.mu.Unlock()
		return
	}
	
	// Copy buffer and clear it
	toSend := make([]SDKEvaluationMetric, len(m.buffer))
	copy(toSend, m.buffer)
	m.buffer = m.buffer[:0]
	m.mu.Unlock()

	payload := map[string]interface{}{
		"evaluations": toSend,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[flagmanagment-sdk-metrics] failed to marshal metrics: %v", err)
		return
	}

	url := m.endpointURL
	if url == "" {
		url = "http://localhost:8080"
	}
	url = url + "/api/v1/sdk/metrics"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.apiKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[flagmanagment-sdk-metrics] failed to sync metrics: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		log.Printf("[flagmanagment-sdk-metrics] unexpected status syncing metrics: %d", resp.StatusCode)
	}
}
