# Tasks: Edge Proxy / Relay Node

**Feature**: `019-edge-proxy-relay`
**Branch**: `019-edge-proxy-relay`
**Plan**: [plan.md](plan.md) | **Spec**: [spec.md](spec.md)

---

## Phase 1: Setup (Project Initialization)

**Purpose**: Create the `proxy/` Go module and directory skeleton so all subsequent tasks have a compile target.

- [ ] T001 Create top-level `proxy/` directory and initialize Go module (`proxy/go.mod`) with module path `github.com/flagmanagment/proxy`
- [ ] T002 [P] Create directory skeleton: `proxy/cmd/proxy/`, `proxy/internal/config/`, `proxy/internal/upstream/`, `proxy/internal/store/`, `proxy/internal/broadcaster/`, `proxy/internal/server/`, `proxy/internal/health/`
- [ ] T003 [P] Add Go dependencies to `proxy/go.mod`: `google.golang.org/grpc`, `github.com/go-chi/chi/v5`, `github.com/rs/zerolog`, `github.com/google/uuid`
- [ ] T004 Copy `backend/pkg/gen/sdk/v1/` generated proto types to `proxy/pkg/gen/sdk/v1/` (or add as a local `replace` directive pointing at the backend module — no duplication of proto source)
- [ ] T005 [P] Create `proxy/Dockerfile` (multi-stage, multi-arch: `golang:1.21-alpine` build → `gcr.io/distroless/static` runtime, `CGO_ENABLED=0`)
- [ ] T006 [P] Create `proxy/Dockerfile.dev` (hot-reload with `air`)

**Checkpoint**: `cd proxy && go build ./...` succeeds (compiles empty stubs).

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core config and in-memory store that every user story depends on. Must complete before Phase 3.

- [ ] T007 Implement `proxy/internal/config/config.go` — `ProxyConfig` struct and `Load()` function reading all env vars (`FM_PROXY_BACKEND_ADDR`, `FM_PROXY_ENV_TOKEN`, `FM_PROXY_GRPC_PORT`, `FM_PROXY_HEALTH_PORT`, `FM_PROXY_UPSTREAM_TLS`, `FM_PROXY_TLS_CERT`, `FM_PROXY_TLS_KEY`, `FM_PROXY_LOG_FORMAT`) with defaults
- [ ] T008 Implement `proxy/internal/store/ruleset.go` — thread-safe `RulesetStore` with `Get()`, `Set()`, `Version()` methods using `sync.RWMutex` protecting `*pb.RulesetSnapshot` and `version string`
- [ ] T009 Implement `proxy/internal/broadcaster/broadcaster.go` — `Broadcaster` with `Register(client)`, `Deregister(clientID)`, `Broadcast(delta)` (non-blocking, per-client buffered channels of size 10), `Count() int`
- [ ] T010 Write unit tests for `RulesetStore` in `proxy/internal/store/ruleset_test.go` — concurrent Get/Set, version updates
- [ ] T011 Write unit tests for `Broadcaster` in `proxy/internal/broadcaster/broadcaster_test.go` — fan-out to N clients, slow consumer drop, deregister during broadcast

**Checkpoint**: `go test -race ./internal/store/... ./internal/broadcaster/...` passes with no race conditions.

---

## Phase 3: User Story 1 — Private Network Isolation (Priority: P1) 🎯 MVP

**Goal**: A proxy instance boots, connects to the FlagManagment backend, loads its ruleset, and serves SDK clients entirely within a private subnet with zero client outbound internet access required.

**Independent Test**: SDK client pointed at `localhost:9091` (proxy) fetches a `FetchSnapshot` response matching what the backend returns directly. No direct route from SDK client process to the backend needed.

### Implementation for User Story 1

- [ ] T012 [US1] Implement `proxy/internal/upstream/client.go` — `UpstreamClient` struct that wraps `pb.SDKServiceClient`; `Connect(ctx, cfg)` creates a gRPC client connection (with optional TLS from config); `Bootstrap(ctx) (*pb.RulesetSnapshot, error)` calls `FetchSnapshot`; `StreamDeltas(ctx) (<-chan *pb.RulesetDelta, error)` calls `StreamRulesets` and returns a channel
- [ ] T013 [US1] Implement reconnection loop in `proxy/internal/upstream/client.go` — `Run(ctx, store, broadcaster)` goroutine that: bootstraps snapshot → populates `RulesetStore` → streams deltas → on disconnect applies exponential backoff (1s start, 2x factor, 60s max, ±25% jitter) → reconnects
- [ ] T014 [US1] Implement `proxy/internal/server/service.go` — `ProxyServer` implementing `pb.SDKServiceServer`; `FetchSnapshot`: validates token matches `FM_PROXY_ENV_TOKEN`, returns `RulesetStore.Get()`; `StreamRulesets`: validates token, registers client with `Broadcaster`, reads delta channel and streams to client, deregisters on disconnect or context cancel
- [ ] T015 [US1] Implement `proxy/cmd/proxy/main.go` — load config, init zerolog logger, create `RulesetStore`, `Broadcaster`, `UpstreamClient`; start upstream `Run()` goroutine; start downstream gRPC server on `FM_PROXY_GRPC_PORT`; graceful shutdown on SIGTERM/SIGINT
- [ ] T016 [US1] Implement `proxy/internal/upstream/state.go` — `UpstreamState` struct with `connected bool`, `connectedSince *time.Time`, `lastDeltaAt *time.Time`, `mu sync.RWMutex`; `SetConnected(bool)`, `RecordDelta()`, `Snapshot() UpstreamStateSnapshot` methods (used by health handler)
- [ ] T017 [US1] Write integration test `proxy/internal/upstream/client_test.go` — spin up a mock `SDKServiceServer` in-process, verify `Bootstrap()` populates store, verify `StreamDeltas()` delivers a delta when mock server sends one

**Checkpoint**: Run Quickstart Scenario 1 and Scenario 2. `grpcurl FetchSnapshot` returns correct flags. Health reports `healthy`.

---

## Phase 4: User Story 2 — Operational Observability (Priority: P2)

**Goal**: Operators can monitor proxy health via the `/healthz` HTTP endpoint and detect upstream disconnection.

**Independent Test**: Start the proxy in healthy state. `curl /healthz` returns `status: healthy`. Stop the backend. `curl /healthz` returns `status: degraded`. Restart the backend. `curl /healthz` returns `status: healthy`.

### Implementation for User Story 2

- [ ] T018 [US2] Implement `proxy/internal/health/handler.go` — `HealthHandler` accepting `*UpstreamState` and `*Broadcaster`; `ServeHTTP`: builds `HealthResponse` from state snapshot and client count; returns HTTP 200 for `healthy`/`degraded`, HTTP 503 for `starting_up`
- [ ] T019 [US2] Wire health HTTP server into `proxy/cmd/proxy/main.go` — start Chi router on `FM_PROXY_HEALTH_PORT` with `GET /healthz` route served by `HealthHandler`; graceful shutdown on SIGTERM
- [ ] T020 [P] [US2] Write unit test `proxy/internal/health/handler_test.go` — mock `UpstreamState` and `Broadcaster`; assert HTTP 200 + correct JSON for healthy, HTTP 200 + `degraded` status when upstream disconnected, HTTP 503 for `starting_up`

**Checkpoint**: Run Quickstart Scenario 4 (upstream outage). Health endpoint correctly reports `degraded` when backend is stopped.

---

## Phase 5: User Story 3 — Multi-SDK Client Fan-out (Priority: P2)

**Goal**: The proxy sustains 500 concurrent downstream SDK client connections and propagates flag deltas to all clients within 500ms.

**Independent Test**: Spawn 100 goroutines each holding a `StreamRulesets` connection. Trigger one delta from the upstream mock. Assert all 100 streams receive the delta within 500ms.

### Implementation for User Story 3

- [ ] T021 [US3] Add concurrent connection tracking to `Broadcaster` — ensure `Register` and `Deregister` are safe under high concurrency (`go test -race -count=100`); add `Count()` method used by health handler
- [ ] T022 [US3] Write stress test `proxy/internal/broadcaster/broadcaster_stress_test.go` — 500 registered clients, one broadcast, assert all 500 channels received delta within 500ms using `time.After` assertion; verify no goroutine leaks after `Deregister`
- [ ] T023 [US3] Implement `StreamRulesets` bootstrap delivery in `proxy/internal/server/service.go` — on client connect, immediately send a synthetic `RulesetDelta{FullReset: true, Version: store.Version()}` so the reconnecting client gets the current version without re-bootstrapping against the backend

**Checkpoint**: Run Quickstart Scenario 3 (100 concurrent streaming clients, one delta broadcast). All clients receive within 500ms. `go test -race ./...` passes.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Production hardening, Docker image validation, and docker-compose integration.

- [ ] T024 [P] Add `docker-compose.edge-proxy.yml` (override file) adding an optional `edge-proxy` service: uses `proxy/Dockerfile`, mounts env vars, exposes ports `9091` and `8081`, `depends_on: backend`
- [ ] T025 [P] Add `proxy/README.md` documenting all environment variables, quickstart steps, and SDK configuration change required (endpoint only)
- [ ] T026 Validate Docker image size: run `docker build -t flagmanagment-proxy:ci proxy/` and assert `docker image inspect --format='{{.Size}}'` < 30MB (add as CI step or note in quickstart)
- [ ] T027 [P] Configure multi-arch build in `proxy/Dockerfile` using `--platform linux/amd64,linux/arm64` via `docker buildx`; test ARM64 build completes without errors
- [ ] T028 Run full test suite `cd proxy && go test -race ./...` and confirm all tests pass with zero race conditions
- [ ] T029 Run all 6 Quickstart validation scenarios from `specs/019-edge-proxy-relay/quickstart.md` and confirm each passes
- [ ] T030 [P] Update root `README.md` (or `docs/`) to document the Edge Proxy component, when to use it, and link to `proxy/README.md`

**Checkpoint**: All Quickstart scenarios pass. `go test -race ./...` green. Docker image < 30MB. ARM64 build succeeds.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 — BLOCKS all user stories
- **US1 / Phase 3**: Depends on Phase 2
- **US2 / Phase 4**: Depends on Phase 2 (and uses `UpstreamState` from Phase 3 T016 — implement T016 first)
- **US3 / Phase 5**: Depends on Phase 2, depends on Phase 3 broadcaster foundation
- **Polish (Phase 6)**: Depends on Phases 3 + 4 + 5

### User Story Dependencies

| Story | Depends On | Notes |
|-------|-----------|-------|
| US1 (P1) | Phase 2 | Core upstream connection — MVP |
| US2 (P2) | Phase 2, T016 from US1 | Needs `UpstreamState`; can start after T016 completes |
| US3 (P2) | Phase 2, Phase 3 | Needs broadcaster fully implemented |

### Parallel Opportunities

- T002, T003, T005, T006 (Phase 1) — all parallel
- T010, T011 (Phase 2 tests) — parallel with each other
- T012, T016 (Phase 3 setup tasks) — parallel
- T020 (Phase 4 test), T021 (Phase 5 broadcaster) — parallel after Phase 2

---

## Parallel Execution Example: Phase 1 Setup

```
Task: T002 — Create directory skeleton
Task: T003 — Add Go dependencies
Task: T005 — Create Dockerfile
Task: T006 — Create Dockerfile.dev
```

## Parallel Execution Example: Phase 3 (US1)

```
Task: T012 — UpstreamClient.Connect + Bootstrap
Task: T016 — UpstreamState struct
  ↓ (after T012 + T016 done)
Task: T013 — Reconnection loop (depends on T012)
Task: T014 — ProxyServer gRPC service (depends on T012, T016)
  ↓
Task: T015 — main.go wiring
Task: T017 — Integration test
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: US1 (T012 → T016 → T013 → T014 → T015)
4. **STOP and VALIDATE**: Run Quickstart Scenarios 1, 2, 4 (health, snapshot, outage resilience)
5. Demo: SDK client pointed at proxy, backend firewalled — works end-to-end

### Incremental Delivery

1. Setup + Foundational → module compiles
2. US1 → private network isolation works (MVP)
3. US2 → health endpoint for monitoring
4. US3 → validated 500-client fan-out performance
5. Polish → production Docker image, docker-compose integration

---

## Notes

- `[P]` = parallelizable (no shared file dependencies)
- `[USn]` = maps to user story n from spec.md
- The proxy module is entirely independent — no imports from `backend/` except shared proto types
- The proto types can be shared via a `replace` directive or by copying generated files — prefer `replace` to avoid drift
- No database migrations, no Redis, no PostgreSQL — this is a pure in-memory relay
- Commit after each phase checkpoint
