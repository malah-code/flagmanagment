package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flagmanagment/proxy/internal/broadcaster"
	"github.com/flagmanagment/proxy/internal/config"
	"github.com/flagmanagment/proxy/internal/store"
	"github.com/flagmanagment/proxy/internal/upstream"
	pb "github.com/flagmanagment/proxy/pkg/gen/sdk/v1"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	cfg := &config.Config{BackendAddr: "localhost:9090"}
	rStore := store.NewRulesetStore()
	b := broadcaster.NewBroadcaster()
	uState := upstream.NewUpstreamState()

	handler := NewHealthHandler(cfg, rStore, b, uState)

	// 1. Starting Up (no ruleset version)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	var resp HealthResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "starting_up", resp.Status)
	assert.False(t, resp.UpstreamConnected)

	// 2. Healthy (ruleset set & upstream connected)
	rStore.Set(&pb.RulesetSnapshot{Version: "v1.0"})
	uState.SetConnected(true)

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", resp.Status)
	assert.True(t, resp.UpstreamConnected)
	assert.Equal(t, "v1.0", resp.RulesetVersion)

	// 3. Degraded (ruleset set & upstream disconnected)
	uState.SetConnected(false)

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "degraded", resp.Status)
	assert.False(t, resp.UpstreamConnected)
	assert.Equal(t, "v1.0", resp.RulesetVersion)
}
