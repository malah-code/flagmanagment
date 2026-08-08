# Implementation Plan: Edge Proxy / Relay Node

**Branch**: `019-edge-proxy-relay` | **Date**: 2026-08-08 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/019-edge-proxy-relay/spec.md`

---

## Summary

The Edge Proxy / Relay Node is a lightweight, stateless Go microservice that sits inside a private corporate network, maintains a single persistent gRPC upstream connection to the FlagManagment backend, bootstraps a full in-memory flag ruleset, and serves that ruleset to all internal microservice SDKs. This eliminates the requirement for every internal service to have individual outbound internet access to the FlagManagment cloud. The proxy receives real-time delta updates from the backend and fans them out to all connected downstream SDK clients, ensuring sub-500ms propagation. It continues serving evaluations from its last known good in-memory state if the upstream connection drops.

---

## Technical Context

**Language/Version**: Go 1.21+ (multi-arch: `linux/amd64`, `linux/arm64` for Raspberry Pi 4)

**Primary Dependencies**:
- `google.golang.org/grpc` — upstream connection to FlagManagment backend + downstream SDK server
- `github.com/redis/go-redis/v9` — NOT needed; proxy is Redis-free (self-contained in-memory state)
- `github.com/go-chi/chi/v5` — HTTP health check endpoint
- `zerolog` — structured logging (consistent with backend)

**Storage**: None — entirely in-memory. Ruleset snapshot is fetched at startup and updated via streaming deltas. No persistence between restarts.

**Testing**: `go test ./...` — unit tests for in-memory store, fan-out broadcaster, reconnection logic; integration test for end-to-end snapshot + delta propagation

**Target Platform**: Docker multi-arch image (`linux/amd64`, `linux/arm64`). Target image size < 30MB compressed using scratch/distroless base.

**Project Type**: Standalone microservice — lives as `proxy/` at the repository root

**Performance Goals**:
- Sustain 500 concurrent downstream SDK connections per instance
- Delta fan-out latency: < 500ms under 500 clients
- Memory footprint: < 128MB for a typical ruleset (10,000 flags)

**Constraints**:
- Zero Redis dependency (proxy is the sole cache layer between backend and SDK clients)
- Configuration via environment variables only — no config files
- Single environment token per proxy instance
- Docker image < 30MB compressed

**Scale/Scope**:
- Single binary, single process
- Horizontally scalable: multiple proxy instances can run for the same environment if needed (each maintains its own upstream connection)

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| **I. API-First Contract Design** | ✅ PASS | Proxy exposes the same `SDKService` gRPC contract as the backend. Contract already exists. Health HTTP endpoint will have its OpenAPI schema documented. |
| **II. Environment Isolation** | ✅ PASS | Each proxy instance binds to exactly one environment token. Cross-environment access is structurally impossible — the proxy only holds the ruleset for its configured token. |
| **III. Governance by Default** | ✅ PASS | No governance impact. Proxy is an infrastructure component; all governance happens at the backend. |
| **IV. Local Evaluation Performance** | ✅ PASS | Proxy holds flag state in-memory. SDK clients connected to the proxy continue local in-memory evaluation. No network call introduced in the evaluation hot path. |
| **V. Test-First Quality Gates** | ✅ PASS | Unit tests for in-memory broadcaster and reconnection logic; integration test for end-to-end flow required before marking complete. |
| **VI. OpenFeature Interoperability** | ✅ PASS | Proxy implements the same `SDKService` gRPC interface. SDK clients require only an endpoint configuration change — no code changes. |
| **VII. PII Protection & Compliance** | ✅ PASS | Proxy is a transparent relay — no PII enters or is stored. Ruleset data that passes through is already sanitized at the backend. |
| **VIII. Cloud-Native Portability** | ✅ PASS | Multi-arch Docker image, environment-variable-driven config. No cloud-provider-specific dependencies. |
| **Technology Stack** | ✅ PASS | Go, gRPC, Chi HTTP — all within the approved technology stack. |

**Post-Phase-1 Constitution Re-check**: No violations introduced by Phase 1 design.

---

## Project Structure

### Documentation (this feature)

```text
specs/019-edge-proxy-relay/
├── plan.md              ← this file
├── research.md          ← Phase 0 output
├── data-model.md        ← Phase 1 output
├── quickstart.md        ← Phase 1 output
├── contracts/           ← Phase 1 output
│   ├── health-api.yaml  ← OpenAPI for health endpoint
│   └── proxy-grpc.md    ← gRPC interface contract
└── tasks.md             ← Phase 2 (/speckit-tasks output)
```

### Source Code

```text
proxy/                           ← NEW top-level service directory
├── cmd/
│   └── proxy/
│       └── main.go              ← entrypoint, config loading, wiring
├── internal/
│   ├── config/
│   │   └── config.go            ← env-var config struct
│   ├── upstream/
│   │   └── client.go            ← gRPC client to FlagManagment backend
│   ├── store/
│   │   └── ruleset.go           ← in-memory ruleset store (thread-safe)
│   ├── broadcaster/
│   │   └── broadcaster.go       ← fan-out delta push to downstream clients
│   ├── server/
│   │   └── service.go           ← downstream gRPC server (SDKService impl)
│   └── health/
│       └── handler.go           ← HTTP health check handler
├── Dockerfile                   ← multi-stage, multi-arch build
├── Dockerfile.dev               ← dev hot-reload
└── go.mod                       ← separate Go module (proxy/go.mod)

docker-compose.yml               ← MODIFIED: add `edge-proxy` service (optional, off by default)
```

**Structure Decision**: New top-level `proxy/` service, independent Go module to keep the binary minimal and free of backend/database dependencies.

---

## Complexity Tracking

No constitution violations. No complexity justifications required.
