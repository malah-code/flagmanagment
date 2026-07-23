# Quickstart: FlagManagment Core Architecture Validation

**Feature**: `001-core-architecture`
**Date**: 2026-07-18

---

## Prerequisites

| Requirement | Minimum Version | Verification Command |
|-------------|----------------|---------------------|
| Docker | 24.0+ | `docker --version` |
| Docker Compose | v2.20+ | `docker compose version` |
| Git | 2.30+ | `git --version` |
| Make (optional) | Any | `make --version` |

> The bootstrap script (`scripts/bootstrap.sh`) validates all prerequisites
> automatically and provides actionable error messages if anything is missing.

---

## Scenario 1: Local Development Bootstrap (US1)

**Goal**: Verify that a single command starts all services and they reach a healthy state.

### Steps

```bash
# 1. Clone the repository
git clone git@github.com:<org>/FlagManagment.git
cd FlagManagment

# 2. (Optional) Copy environment defaults
cp .env.example .env

# 3. Bootstrap all services
make up
# Or directly: ./scripts/bootstrap.sh
```

### Expected Outcome

- All 4 services start: backend, frontend, postgres, redis
- Terminal output shows health check passing within 60 seconds
- No error messages in container logs

### Validation

```bash
# Health check (should return 200 with status: healthy)
curl -s http://localhost:8080/healthz | jq .

# Dashboard (should load in browser)
open http://localhost:3000  # macOS
xdg-open http://localhost:3000  # Linux

# View logs
docker compose logs -f backend
```

**Pass criteria**: Health check returns `{"status": "healthy"}` with both postgres
and redis checks passing. Dashboard loads and displays a green status badge.
Total time from clone to healthy: **< 10 minutes** (SC-001).

---

## Scenario 2: Multi-Architecture Build (US2)

**Goal**: Verify container images build for x86_64 and ARM64.

### Steps

```bash
# Build multi-arch images (requires Docker Buildx)
make build

# Or directly:
docker buildx build --platform linux/amd64,linux/arm64 \
  -t flagmanagment/backend:test \
  -f backend/Dockerfile backend/

docker buildx build --platform linux/amd64,linux/arm64 \
  -t flagmanagment/frontend:test \
  -f frontend/Dockerfile frontend/
```

### Expected Outcome

- Build completes for both platforms without errors
- Image manifest includes both `linux/amd64` and `linux/arm64` entries

### Validation

```bash
# Inspect the manifest
docker buildx imagetools inspect flagmanagment/backend:test

# Verify ARM64 image runs (on ARM machine or via QEMU)
docker run --platform linux/arm64 --rm flagmanagment/backend:test --version
```

**Pass criteria**: Both architecture builds succeed. Images produce a valid
version output when run. On Raspberry Pi 4: services start and health check
passes (SC-002).

---

## Scenario 3: CI Pipeline Verification (US3)

**Goal**: Verify that the CI pipeline catches lint violations and enforces coverage.

### Steps

```bash
# Run linting locally (mirrors CI)
make lint

# Run tests with coverage (mirrors CI)
make test

# Verify coverage threshold
# Backend: go test -coverprofile + coverage threshold check
# Frontend: vitest --coverage with threshold in vitest.config.ts
```

### Expected Outcome

- Clean code: all linters pass, all tests pass, coverage meets thresholds
- Intentional violation: linters report specific errors

### Validation

```bash
# Introduce a lint violation (e.g., unused variable in Go)
echo 'package main; var unused int' >> backend/cmd/server/main.go
make lint
# Expected: golangci-lint reports "unused" error

# Revert
git checkout backend/cmd/server/main.go
```

**Pass criteria**: CI catches 100% of lint violations (SC-003). Coverage report
generated and threshold enforced (80% backend, 70% frontend — FR-009).

---

## Scenario 4: Hot-Reload Verification (US1.4)

**Goal**: Verify code changes reflect without restarting the stack.

### Steps

```bash
# Ensure stack is running
make up

# Modify the health check response (e.g., change version string)
# Edit backend/internal/health/handler.go
# Save the file

# Watch for air to rebuild
docker compose logs -f backend
```

### Expected Outcome

- `air` detects the file change within 2 seconds
- Binary rebuilds and restarts within 3 seconds
- Total: change reflected in < 5 seconds

### Validation

```bash
# Before change
curl -s http://localhost:8080/healthz | jq .version

# After saving the modified file, wait 5 seconds
sleep 5
curl -s http://localhost:8080/healthz | jq .version
# Expected: version string reflects the change
```

**Pass criteria**: Backend hot-reload reflects changes within 5 seconds of file
save (SC-005).

---

## Scenario 5: Resource Consumption Check (SC-007)

**Goal**: Verify the backend engine stays under 250MB RAM.

### Steps

```bash
# Start the stack
make up

# Wait for services to stabilize (30 seconds)
sleep 30

# Check backend container memory usage
docker stats --no-stream --format "table {{.Name}}\t{{.MemUsage}}" | grep backend
```

### Expected Outcome

- Backend memory usage is significantly below 250MB (expected: 20-50MB for the
  skeleton with no business logic)

**Pass criteria**: Backend engine operates under 250MB RAM (SC-007).

---

## Teardown

```bash
# Stop all services
make down

# Remove all data volumes (clean slate)
docker compose down -v
```

---

## Quick Reference

| Command | Purpose |
|---------|---------|
| `make up` | Bootstrap all services |
| `make down` | Stop all services |
| `make test` | Run all tests |
| `make lint` | Run all linters |
| `make build` | Multi-arch Docker build |
| `curl localhost:8080/healthz` | Check backend health |
| `open localhost:3000` | Open dashboard |
