# Research: Edge Proxy / Relay Node

**Feature**: `019-edge-proxy-relay`
**Date**: 2026-08-08

---

## Decision 1: Architecture Pattern — Proxy-as-gRPC-Server vs. HTTP Forward Proxy

**Decision**: The Edge Proxy implements the exact same `SDKService` gRPC interface as the FlagManagment backend. It acts as a transparent gRPC server downstream, while being a gRPC client upstream.

**Rationale**: SDK clients require zero code changes — only the endpoint address changes. This is the cleanest transparency model. An HTTP forward proxy would require significant SDK plumbing and break the OpenFeature contract.

**Alternatives considered**:
- HTTP CONNECT tunnel — rejected; would expose raw gRPC frames over HTTP/1.1, breaking streaming.
- Sidecar (Envoy/NGINX) — rejected; adds operational complexity, not a single stateless binary, and cannot implement in-memory fan-out semantics.

---

## Decision 2: In-Memory State — No Redis Required

**Decision**: The proxy maintains its own `sync.RWMutex`-protected in-memory ruleset. No Redis instance is required for the proxy itself.

**Rationale**: Adding Redis as a proxy dependency defeats the purpose of the relay node (simplicity, minimal deps). The proxy bootstraps from the backend's `FetchSnapshot` RPC, then applies deltas from `StreamRulesets`. It becomes a self-sufficient caching layer for its downstream clients.

**Alternatives considered**:
- Redis for ruleset cache — rejected; unnecessary complexity, adds infrastructure requirement, and Redis connectivity from the proxy still requires outbound access from the subnet.

---

## Decision 3: Reconnection Strategy — Exponential Backoff with Jitter

**Decision**: On upstream connection loss, the proxy uses exponential backoff starting at 1s, doubling to a maximum of 60s, with ±25% random jitter to avoid thundering herd.

**Rationale**: Standard industry pattern for gRPC streaming reconnection. During the backoff period, the proxy continues serving its in-memory state (last known good). The health endpoint reports `degraded` immediately upon loss.

**Alternatives considered**:
- Fixed retry interval — rejected; causes thundering herd under backend restarts.
- Immediate fail-fast — rejected; violates the constitution requirement that SDKs must continue evaluating using last known good state.

---

## Decision 4: Fan-out Broadcaster — Channel-Per-Client vs. Shared Channel

**Decision**: The broadcaster maintains a `map[clientID]*chan Delta` (one write channel per connected downstream client). When a delta arrives from upstream, it is written to all client channels concurrently using goroutines.

**Rationale**: Non-blocking fan-out. A slow consumer cannot stall the broadcaster or other clients. Each client has its own buffered channel (buffer size: 10 deltas). If a client's channel is full (backpressure), the broadcaster logs a warning and drops the message for that client (the client will receive the next delta or can re-bootstrap).

**Alternatives considered**:
- Single shared broadcast channel — rejected; slow consumer stalls all others.
- Redis pub/sub at proxy level — rejected; reintroduces Redis dependency.

---

## Decision 5: Docker Image Strategy — Multi-Stage Scratch Build

**Decision**: Multi-stage Dockerfile using `golang:1.21-alpine` for build → `scratch` (or `gcr.io/distroless/static`) for runtime. `CGO_ENABLED=0`, `GOOS=linux` for fully static binary.

**Rationale**: Achieves < 30MB compressed image target while supporting multi-arch (`linux/amd64`, `linux/arm64`). A static Go binary with no libc dependency runs on any Linux kernel including Raspberry Pi 4 (ARM64).

**Alternatives considered**:
- Alpine runtime image — produces ~15MB image but adds busybox/libc surface area unnecessarily.
- Debian slim — produces ~80MB compressed, exceeds the 30MB target.

---

## Decision 6: Single Environment Token per Instance

**Decision**: Each proxy instance is configured with exactly one `FM_PROXY_ENV_TOKEN`. Multiple environments require multiple proxy instances.

**Rationale**: Simplest operational model. Environment isolation is maintained structurally (no multiplexing code needed). Operators can run multiple proxy instances behind a load balancer for the same environment, or separate instances for separate environments.

**Alternatives considered**:
- Multi-tenant proxy (multiple tokens) — rejected; dramatically increases complexity, memory footprint, and operational debugging surface.

---

## Decision 7: TLS Configuration

**Decision**: Both upstream (proxy → backend) and downstream (SDK → proxy) TLS are configurable via environment variables. TLS is enabled by default for upstream; downstream TLS is opt-in via `FM_PROXY_TLS_CERT` / `FM_PROXY_TLS_KEY`.

**Rationale**: Many private network deployments terminate TLS at the load balancer level (mTLS not universally supported). Forcing TLS downstream would block quick internal deployments. Making it configurable preserves security for those who need it without blocking simpler deployments.
