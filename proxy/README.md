# FlagManagment Edge Proxy / Relay Node

The **FlagManagment Edge Proxy** is a lightweight, stateless microservice designed for deployment inside private corporate subnets or edge networks (such as air-gapped environments or local Kubernetes clusters).

It maintains a single persistent gRPC connection to the central FlagManagment backend, bootstraps an in-memory ruleset snapshot, receives real-time delta updates, and serves internal microservice SDK clients locally.

---

## Key Benefits

- **Zero Outbound Internet per Microservice**: Only the Edge Proxy instance requires connectivity to the FlagManagment backend. Internal SDK clients point directly to the proxy.
- **In-Memory Local Evaluation**: Downstream SDK clients evaluate flags with sub-millisecond latency.
- **Outage Resilience**: If the upstream connection to the backend drops, the proxy continues serving evaluation requests from its last known good in-memory state.
- **High Concurrency & Fanout**: Sustains 500+ concurrent SDK client connections per instance with non-blocking delta fanout (< 500ms).
- **Stateless & Minimal Footprint**: Single Go binary with zero database or Redis dependencies. Compressed Docker image < 30MB.

---

## Configuration

Driven entirely by environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `FM_PROXY_BACKEND_ADDR` | gRPC host:port of FlagManagment backend | `localhost:9090` |
| `FM_PROXY_ENV_TOKEN` | Secret SDK authentication token for the target environment | `""` |
| `FM_PROXY_GRPC_PORT` | Port exposed for downstream SDK client connections | `9091` |
| `FM_PROXY_HEALTH_PORT` | Port exposed for HTTP health checks | `8081` |
| `FM_PROXY_UPSTREAM_TLS` | Enable TLS for upstream connection to backend | `false` |
| `FM_PROXY_TLS_CERT` | File path to TLS cert for downstream SDK connections (optional) | `""` |
| `FM_PROXY_TLS_KEY` | File path to TLS key for downstream SDK connections (optional) | `""` |
| `FM_PROXY_LOG_FORMAT` | Logging output format (`auto`, `text`, `json`) | `auto` |

---

## Running Locally

### Direct Execution

```bash
cd proxy
export FM_PROXY_BACKEND_ADDR="localhost:9090"
export FM_PROXY_ENV_TOKEN="env_abc123..."
go run cmd/proxy/main.go
```

### Docker Compose

```bash
docker-compose -f docker-compose.yml -f docker-compose.edge-proxy.yml up --build
```

---

## Monitoring & Health Checks

The HTTP `/healthz` endpoint on port `8081` returns the proxy status:

```bash
curl http://localhost:8081/healthz
```

Example output:
```json
{
  "status": "healthy",
  "upstream_connected": true,
  "upstream_addr": "localhost:9090",
  "connected_since": "2026-08-08T12:00:00Z",
  "last_delta_at": "2026-08-08T12:05:30Z",
  "downstream_clients": 42,
  "ruleset_version": "v1720123456"
}
```

Status levels:
- `healthy`: Upstream stream connected and serving live ruleset.
- `degraded`: Upstream stream disconnected; proxy is serving from its last known good in-memory state.
- `starting_up`: Initial ruleset snapshot has not been fetched yet (HTTP 503).

---

## SDK Configuration

SDK clients connect to the Edge Proxy identically to how they connect to the central backend. Simply point your SDK's gRPC endpoint to the proxy address:

```go
// Example Go SDK / OpenFeature Provider setup
provider := flagmanagement.NewProvider(flagmanagement.Config{
    GRPCEndpoint: "edge-proxy.internal:9091", // Point to Edge Proxy instead of central backend
    APIKey:       "env_abc123...",
})
```
