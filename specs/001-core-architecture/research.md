# Research: FlagManagment Core Architecture & Repository Bootstrap

**Feature**: `001-core-architecture`
**Date**: 2026-07-18

---

## 1. Go HTTP Router Selection

**Decision**: `chi` (go-chi/chi v5)

**Rationale**: chi is a lightweight, idiomatic HTTP router built on the standard `net/http` stack. It provides middleware chaining, route grouping, and URL parameter extraction without deviating from Go's stdlib interfaces. This means health check handlers and future API handlers remain testable with `httptest` directly — no framework lock-in.

**Alternatives considered**:
- **echo**: Full framework with built-in middleware, but introduces a custom `Context` interface that diverges from `net/http`. Adds unnecessary abstraction for a project that will later add gRPC alongside REST.
- **gin**: Popular but uses a custom handler signature (`gin.Context`) that makes stdlib-based testing and middleware reuse harder.
- **stdlib only (Go 1.22 ServeMux)**: Go 1.22 added method-based routing to the stdlib. Viable, but lacks middleware chaining and route grouping — chi adds this with zero overhead while staying stdlib-compatible.

---

## 2. Structured Logging Library

**Decision**: `zerolog` (rs/zerolog)

**Rationale**: zerolog provides zero-allocation structured logging with both human-readable (console) and JSON output modes, directly supporting the clarified requirement of auto-detecting environment (text locally, JSON in production). It has excellent performance characteristics — critical given the <250MB RAM constraint for Raspberry Pi 4 deployment. The `ConsoleWriter` for dev mode and default JSON encoder for production can be selected via a single environment variable.

**Alternatives considered**:
- **zap (uber-go/zap)**: Excellent performance but heavier API surface. `SugaredLogger` is convenient but allocates. zap's production/development presets map well to our needs, but zerolog is simpler to configure for the text/JSON auto-detect pattern.
- **slog (Go 1.21+ stdlib)**: Built into the standard library. Suitable for simple use cases but lacks the zero-allocation guarantees and the polished console writer needed for dev-mode readability.
- **logrus**: Widely used but in maintenance mode. Slower than zerolog/zap due to reflection-based field serialization.

---

## 3. Docker Multi-Architecture Build Strategy

**Decision**: Docker Buildx with `docker buildx build --platform linux/amd64,linux/arm64`

**Rationale**: Docker Buildx is the standard tool for multi-platform image builds. Using multi-stage Dockerfiles with Go's built-in cross-compilation (`GOOS=linux GOARCH=$TARGETARCH`) eliminates the need for QEMU emulation during the Go compilation step — only the final `COPY` stage uses the target platform. This produces native binaries for each architecture with no performance penalty. For the frontend, Node.js base images already support multi-arch via manifest lists.

**Alternatives considered**:
- **QEMU-based full emulation**: Simpler configuration but 10-50x slower builds for Go compilation. Unacceptable for CI.
- **Separate Dockerfiles per architecture**: Duplicates maintenance burden. A single Dockerfile with `TARGETARCH` build arg handles both.
- **Nix/Bazel cross-compilation**: Overkill for this project size; introduces steep learning curves for contributors.

**Build targets**:
| Platform | Use Case | Verification |
|----------|----------|-------------|
| `linux/amd64` | Servers, Linux workstations, CI runners | Primary CI target |
| `linux/arm64` | Apple Silicon Macs, Raspberry Pi 4, ARM servers | Secondary CI target |

---

## 4. Go Backend Container Image Base

**Decision**: Multi-stage build — `golang:1.22-bookworm` (build), `gcr.io/distroless/static-debian12` (runtime)

**Rationale**: Distroless images contain only the application binary and its runtime dependencies — no shell, package manager, or OS utilities. This minimizes attack surface and image size (~2MB base vs ~140MB for `golang:1.22`). Go produces statically-linked binaries when built with `CGO_ENABLED=0`, so no libc is needed at runtime.

**Alternatives considered**:
- **`alpine` runtime**: Small (~5MB) but includes musl libc, shell, and package manager. Slightly larger attack surface than distroless.
- **`scratch`**: Even smaller than distroless but lacks CA certificates and timezone data needed for HTTPS and log timestamps.
- **`debian-slim` runtime**: ~80MB; unnecessary bloat for a static Go binary.

---

## 5. Frontend Container Image & Dev Mode

**Decision**: Multi-stage build — `node:20-bookworm-slim` (build), `nginx:stable-alpine` (serve). Development mode uses Vite dev server with HMR via volume mounts.

**Rationale**: Vite's dev server provides instant Hot Module Replacement (HMR) for the React frontend. In production, static assets are built via `vite build` and served by nginx. The nginx-alpine base is ~40MB and handles static file serving, gzip compression, and reverse proxying efficiently.

**Dev mode hot-reload implementation**:
- `docker-compose.override.yml` mounts `./frontend/src` into the container
- Vite dev server watches for filesystem changes via polling (necessary inside Docker)
- Changes reflect in the browser in <2 seconds via HMR

---

## 6. Backend Hot-Reload Strategy

**Decision**: `air` (cosmtrek/air) for Go live reloading in development

**Rationale**: `air` watches Go source files and automatically rebuilds/restarts the binary on save. It's the most widely adopted Go hot-reload tool with Docker support. The rebuild cycle for a minimal Go binary is typically <3 seconds, well within the SC-005 target of <5 seconds.

**Alternatives considered**:
- **`reflex`**: Similar functionality but less actively maintained.
- **`CompileDaemon`**: Older, fewer features, less Docker integration documentation.
- **Manual rebuild**: Unacceptable developer experience for a project targeting contributor onboarding.

**Configuration**: `.air.toml` at `backend/` root, configured to:
- Watch `*.go` files recursively
- Exclude `vendor/`, `tmp/`, test files
- Build with `go build -o ./tmp/main ./cmd/server`
- Use polling mode (required for Docker volume mounts)

---

## 7. CI/CD Platform & Pipeline Architecture

**Decision**: GitHub Actions with reusable workflows

**Rationale**: The repository is hosted on GitHub (clarified: private for Phase 1). GitHub Actions provides native integration with GitHub Container Registry (ghcr.io), branch protection rules, and PR status checks. The free tier for private repos provides 2,000 minutes/month — sufficient for Phase 1 development.

**Pipeline architecture**:

| Workflow | Trigger | Steps |
|----------|---------|-------|
| `ci.yml` | Pull request | 1. Lint (golangci-lint, ESLint, Prettier) → 2. Test (go test, vitest) → 3. Coverage report → 4. Coverage threshold check → 5. Docker build (verify) |
| `publish.yml` | Push to `main` | 1. Multi-arch build (amd64 + arm64) → 2. Push to ghcr.io |

**Coverage enforcement**: The CI pipeline will use `go test -coverprofile` and parse the output to enforce the 80% threshold. Frontend coverage via Vitest's built-in coverage reporter with a threshold configuration in `vitest.config.ts`.

---

## 8. Bootstrap Script Design

**Decision**: `scripts/bootstrap.sh` — a Bash prerequisite validation script wrapping `docker compose up`

**Rationale**: A thin wrapper script validates the development environment before invoking Docker Compose. This addresses edge cases (FR-014) without introducing additional dependencies.

**Validation checks**:
1. Docker (or Podman) installed → clear error message with install URL if missing
2. Docker Compose v2 available → error with upgrade instructions if missing
3. Docker daemon running → prompt to start if stopped
4. Minimum Docker version check (24.0+) → warning if outdated
5. Port availability scan (default ports: 8080 backend, 3000 frontend, 5432 PostgreSQL, 6379 Redis) → error listing conflicts with suggestion to configure via `.env`
6. Architecture detection → warning if unsupported (e.g., 32-bit ARM)

**Makefile targets** (convenience wrappers):
- `make up` → `scripts/bootstrap.sh`
- `make down` → `docker compose down`
- `make test` → run all tests
- `make lint` → run all linters
- `make build` → multi-arch Docker build

---

## 9. Environment Configuration Strategy

**Decision**: `.env.example` file with documented defaults, loaded by Docker Compose

**Rationale**: Docker Compose natively loads `.env` files. Developers copy `.env.example` to `.env` (or use defaults directly). All configuration is driven by environment variables — no hardcoded values in source code (FR-002).

**Key environment variables**:

| Variable | Default | Service |
|----------|---------|---------|
| `FM_BACKEND_PORT` | `8080` | Backend |
| `FM_FRONTEND_PORT` | `3000` | Frontend |
| `FM_DB_HOST` | `postgres` | Backend |
| `FM_DB_PORT` | `5432` | PostgreSQL |
| `FM_DB_USER` | `flagmgmt` | PostgreSQL |
| `FM_DB_PASSWORD` | `flagmgmt_dev` | PostgreSQL |
| `FM_DB_NAME` | `flagmanagment` | PostgreSQL |
| `FM_REDIS_HOST` | `redis` | Backend |
| `FM_REDIS_PORT` | `6379` | Redis |
| `FM_LOG_FORMAT` | `auto` | Backend (auto/text/json) |
| `FM_ENV` | `development` | Backend (development/production) |

---

## 10. PostgreSQL & Redis Container Configuration

**Decision**: Use official images `postgres:16-alpine` and `redis:7-alpine` with named volumes

**Rationale**: Alpine variants minimize image size (~80MB Postgres, ~30MB Redis). Named Docker volumes persist data across container restarts without cluttering the repository with data directories.

**Health checks** (in Docker Compose):
- PostgreSQL: `pg_isready -U ${FM_DB_USER} -d ${FM_DB_NAME}`
- Redis: `redis-cli ping`
- Backend depends_on with `condition: service_healthy` to ensure database readiness before API starts
