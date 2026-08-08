# Data Model: Edge Proxy / Relay Node

**Feature**: `019-edge-proxy-relay`

---

## Core In-Memory Entities

### ProxyConfig

Loaded once at startup from environment variables.

| Field | Type | Source | Description |
|-------|------|--------|-------------|
| `BackendAddr` | `string` | `FM_PROXY_BACKEND_ADDR` | gRPC address of FlagManagment backend (e.g., `flagmanagment.internal:9090`) |
| `EnvToken` | `string` | `FM_PROXY_ENV_TOKEN` | Environment SDK token (secret) |
| `GRPCPort` | `string` | `FM_PROXY_GRPC_PORT` (default: `9091`) | Port proxy listens on for downstream SDK connections |
| `HealthPort` | `string` | `FM_PROXY_HEALTH_PORT` (default: `8081`) | Port for HTTP health endpoint |
| `TLSCertFile` | `string` | `FM_PROXY_TLS_CERT` | Path to TLS cert for downstream (optional) |
| `TLSKeyFile` | `string` | `FM_PROXY_TLS_KEY` | Path to TLS key for downstream (optional) |
| `UpstreamTLS` | `bool` | `FM_PROXY_UPSTREAM_TLS` (default: `true`) | Enable TLS for upstream gRPC connection |
| `LogFormat` | `string` | `FM_PROXY_LOG_FORMAT` (default: `auto`) | `json`, `text`, or `auto` |

---

### RulesetStore

Thread-safe in-memory store holding the current flag ruleset for the configured environment.

| Field | Type | Description |
|-------|------|-------------|
| `snapshot` | `*RulesetSnapshot` | Current full ruleset (reused from `pkg/gen/sdk/v1`) |
| `version` | `string` | Monotonically increasing version string |
| `mu` | `sync.RWMutex` | Guards concurrent reads/writes |

**Operations**:
- `Get() *RulesetSnapshot` — thread-safe read
- `Set(snapshot *RulesetSnapshot)` — thread-safe write (called on bootstrap and full reset)
- `Version() string` — returns current version

---

### Client (Downstream SDK Connection)

Represents a single connected downstream SDK client.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | UUID assigned on connection |
| `ch` | `chan *pb.RulesetDelta` | Buffered write channel (buffer: 10) |
| `connectedAt` | `time.Time` | Connection timestamp |

---

### Broadcaster

Manages the registry of connected downstream clients and pushes deltas.

| Field | Type | Description |
|-------|------|-------------|
| `clients` | `map[string]*Client` | Active downstream connections |
| `mu` | `sync.RWMutex` | Guards clients map |

**Operations**:
- `Register(client *Client)` — add a new downstream connection
- `Deregister(clientID string)` — remove on disconnect
- `Broadcast(delta *pb.RulesetDelta)` — non-blocking fan-out to all registered clients
- `Count() int` — current connection count (for health check)

---

### UpstreamState

Tracks the state of the connection to the FlagManagment backend.

| Field | Type | Description |
|-------|------|-------------|
| `connected` | `bool` | Whether the upstream gRPC stream is active |
| `connectedSince` | `*time.Time` | Timestamp of last successful connection |
| `lastDeltaAt` | `*time.Time` | Timestamp of last delta received |
| `mu` | `sync.RWMutex` | Guards fields |

---

## State Transitions

```
Startup
  └─▶ [BOOTSTRAPPING] — FetchSnapshot from backend
        ├─▶ Success → [CONNECTED] — StreamRulesets active, serving clients
        │     ├─▶ Delta received → update RulesetStore → Broadcast to clients
        │     └─▶ Stream broken → [DEGRADED] — serving from last known good state
        │           └─▶ Reconnect (backoff) → [BOOTSTRAPPING]
        └─▶ Failure → [DEGRADED] — no ruleset yet, health: starting_up
              └─▶ Reconnect (backoff) → [BOOTSTRAPPING]
```

---

## Health Response Schema

The HTTP `/healthz` endpoint returns:

```json
{
  "status": "healthy | degraded | starting_up",
  "upstream_connected": true,
  "upstream_addr": "flagmanagment.internal:9090",
  "connected_since": "2026-08-08T12:00:00Z",
  "last_delta_at": "2026-08-08T12:05:30Z",
  "downstream_clients": 47,
  "ruleset_version": "v1720123456"
}
```

| Field | Description |
|-------|-------------|
| `status` | `healthy` = connected and serving; `degraded` = upstream disconnected, serving last-known-good; `starting_up` = no ruleset yet |
| `upstream_connected` | Boolean upstream stream state |
| `connected_since` | ISO8601 timestamp of last successful upstream connect |
| `last_delta_at` | ISO8601 timestamp of last received delta |
| `downstream_clients` | Current count of connected SDK clients |
| `ruleset_version` | Current in-memory ruleset version string |
