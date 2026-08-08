package health

import (
	"encoding/json"
	"net/http"

	"github.com/flagmanagment/proxy/internal/broadcaster"
	"github.com/flagmanagment/proxy/internal/config"
	"github.com/flagmanagment/proxy/internal/store"
	"github.com/flagmanagment/proxy/internal/upstream"
)

type HealthResponse struct {
	Status            string  `json:"status"`
	UpstreamConnected bool    `json:"upstream_connected"`
	UpstreamAddr      string  `json:"upstream_addr"`
	ConnectedSince    *string `json:"connected_since"`
	LastDeltaAt       *string `json:"last_delta_at"`
	DownstreamClients int     `json:"downstream_clients"`
	RulesetVersion    string  `json:"ruleset_version"`
}

type HealthHandler struct {
	cfg         *config.Config
	store       *store.RulesetStore
	broadcaster *broadcaster.Broadcaster
	state       *upstream.UpstreamState
}

func NewHealthHandler(cfg *config.Config, store *store.RulesetStore, broadcaster *broadcaster.Broadcaster, state *upstream.UpstreamState) *HealthHandler {
	return &HealthHandler{
		cfg:         cfg,
		store:       store,
		broadcaster: broadcaster,
		state:       state,
	}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stateSnap := h.state.Snapshot()
	version := h.store.Version()

	var statusStr string
	var statusCode int

	if version == "" {
		statusStr = "starting_up"
		statusCode = http.StatusServiceUnavailable
	} else if stateSnap.Connected {
		statusStr = "healthy"
		statusCode = http.StatusOK
	} else {
		statusStr = "degraded"
		statusCode = http.StatusOK // 200 OK because it is still serving last-known-good in-memory state
	}

	var connSinceStr *string
	if stateSnap.ConnectedSince != nil {
		formatted := stateSnap.ConnectedSince.Format("2006-01-02T15:04:05Z07:00")
		connSinceStr = &formatted
	}

	var lastDeltaStr *string
	if stateSnap.LastDeltaAt != nil {
		formatted := stateSnap.LastDeltaAt.Format("2006-01-02T15:04:05Z07:00")
		lastDeltaStr = &formatted
	}

	resp := HealthResponse{
		Status:            statusStr,
		UpstreamConnected: stateSnap.Connected,
		UpstreamAddr:      h.cfg.BackendAddr,
		ConnectedSince:    connSinceStr,
		LastDeltaAt:       lastDeltaStr,
		DownstreamClients: h.broadcaster.Count(),
		RulesetVersion:    version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}
