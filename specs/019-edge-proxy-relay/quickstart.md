# Quickstart Validation Guide: Edge Proxy / Relay Node

**Feature**: `019-edge-proxy-relay`

---

## Prerequisites

- Docker and Docker Compose installed
- FlagManagment backend running (local dev or staging)
- A valid environment SDK token from FlagManagment

---

## Setup

### 1. Configure the proxy

```bash
export FM_PROXY_BACKEND_ADDR=localhost:9090     # gRPC address of FlagManagment backend
export FM_PROXY_ENV_TOKEN=env_<your_token>      # SDK token for the environment
export FM_PROXY_GRPC_PORT=9091                  # Port for downstream SDK connections
export FM_PROXY_HEALTH_PORT=8081                # Port for health HTTP endpoint
export FM_PROXY_UPSTREAM_TLS=false              # Disable TLS for local dev
```

### 2. Build and run the proxy

```bash
cd proxy/
go run cmd/proxy/main.go
```

Or with Docker:

```bash
docker build -t flagmanagment-proxy:dev .
docker run \
  -e FM_PROXY_BACKEND_ADDR=host.docker.internal:9090 \
  -e FM_PROXY_ENV_TOKEN=env_<your_token> \
  -e FM_PROXY_UPSTREAM_TLS=false \
  -p 9091:9091 -p 8081:8081 \
  flagmanagment-proxy:dev
```

---

## Validation Scenarios

### Scenario 1: Health Check — Proxy Boots and Connects

**Goal**: Verify the proxy connects to the backend and loads its ruleset.

```bash
# Wait ~2s for startup, then:
curl -s http://localhost:8081/healthz | jq .
```

**Expected output**:
```json
{
  "status": "healthy",
  "upstream_connected": true,
  "upstream_addr": "localhost:9090",
  "downstream_clients": 0,
  "ruleset_version": "v172..."
}
```

✅ Pass if `status == "healthy"` and `upstream_connected == true`.

---

### Scenario 2: SDK Client Connects via Proxy

**Goal**: Verify an SDK client can fetch the ruleset snapshot through the proxy.

Using `grpcurl`:

```bash
grpcurl -plaintext \
  -H "authorization: Bearer env_<your_token>" \
  -d '{"environment_token": "env_<your_token>"}' \
  localhost:9091 \
  flagmanagement.sdk.v1.SDKService/FetchSnapshot
```

**Expected**: Response containing `version` and `flags` array matching the backend's current ruleset.

✅ Pass if response is non-empty and matches direct backend query.

---

### Scenario 3: Delta Propagation

**Goal**: Verify a flag change in the dashboard is received by an SDK client connected to the proxy within 500ms.

1. Open a streaming subscription:
```bash
grpcurl -plaintext \
  -H "authorization: Bearer env_<your_token>" \
  -d '{"environment_token": "env_<your_token>"}' \
  localhost:9091 \
  flagmanagement.sdk.v1.SDKService/StreamRulesets
```

2. In the FlagManagment dashboard, toggle any flag in the configured environment.

**Expected**: A `RulesetDelta` message appears in the stream output within 500ms.

✅ Pass if delta arrives within 500ms.

---

### Scenario 4: Upstream Outage Resilience

**Goal**: Verify the proxy continues serving during upstream disconnect.

1. Verify proxy is healthy (Scenario 1).
2. Stop the FlagManagment backend.
3. Check health immediately:

```bash
curl -s http://localhost:8081/healthz | jq .status
```
**Expected**: `"degraded"` (not `"starting_up"` — ruleset is still in memory).

4. Verify SDK client can still fetch snapshot:
```bash
grpcurl -plaintext \
  -H "authorization: Bearer env_<your_token>" \
  -d '{"environment_token": "env_<your_token>"}' \
  localhost:9091 \
  flagmanagement.sdk.v1.SDKService/FetchSnapshot
```
**Expected**: Same successful response as before.

5. Restart the backend. Wait for reconnect. Check health again.
**Expected**: `"healthy"` within 60s (max backoff).

✅ Pass if all steps succeed.

---

### Scenario 5: Docker Image Size

**Goal**: Validate the compressed Docker image is < 30MB.

```bash
docker build -t flagmanagment-proxy:release \
  --build-arg TARGETOS=linux --build-arg TARGETARCH=amd64 \
  proxy/
docker image ls flagmanagment-proxy:release --format "{{.Size}}"
```

✅ Pass if reported size is < 30MB.

---

### Scenario 6: ARM64 Build (Raspberry Pi 4)

**Goal**: Verify multi-arch build succeeds.

```bash
docker buildx build \
  --platform linux/arm64 \
  -t flagmanagment-proxy:arm64 \
  proxy/
```

✅ Pass if build succeeds without errors.

---

## Running Unit Tests

```bash
cd proxy/
go test -race ./...
```

Expected: all tests pass, no race conditions detected.
