package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// StreamData represents a streaming event payload
type StreamData struct {
	Version  int                    `json:"version"`
	Flags    map[string]interface{} `json:"flags,omitempty"`
	Segments map[string]interface{} `json:"segments,omitempty"`
	FlagKey  string                 `json:"flagKey,omitempty"`
	Flag     interface{}            `json:"flag,omitempty"`
}

// SDKStreamHandler handles SSE connections from SDKs
func SDKStreamHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Ensure the connection supports flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Send initial bootstrap payload
	bootstrapData := StreamData{
		Version: 1,
		Flags: map[string]interface{}{
			"sample-flag": map[string]interface{}{
				"key":            "sample-flag",
				"enabled":        true,
				"type":           "BOOLEAN",
				"defaultVariant": "on",
				"variants":       map[string]interface{}{"on": true, "off": false},
				"rules":          []interface{}{},
			},
		},
		Segments: make(map[string]interface{}),
	}
	sendSSE(w, flusher, "bootstrap", bootstrapData)

	// Keep connection open and send pings
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Wait for connection to close
	notify := r.Context().Done()

	for {
		select {
		case <-notify:
			return
		case <-ticker.C:
			sendSSEPing(w, flusher)
		}
	}
}

func sendSSE(w http.ResponseWriter, flusher http.Flusher, eventType string, data StreamData) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	
	// Replace newlines in JSON to ensure SSE format is respected
	dataStr := strings.ReplaceAll(string(jsonData), "\n", "")
	
	fmt.Fprintf(w, "event: %s\n", eventType)
	fmt.Fprintf(w, "data: %s\n\n", dataStr)
	flusher.Flush()
}

func sendSSEPing(w http.ResponseWriter, flusher http.Flusher) {
	fmt.Fprintf(w, "event: ping\n")
	fmt.Fprintf(w, "data: {\"time\": %d}\n\n", time.Now().Unix())
	flusher.Flush()
}
