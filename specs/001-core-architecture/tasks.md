# Tasks: FlagManagment Core Architecture & Repository Bootstrap

**Input**: Design documents from `specs/001-core-architecture/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story. Each task is written with maximum detail so that a code-generating LLM can execute it without additional context or decision-making.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize the monorepo structure, dependency files, and root-level configuration. No runnable code yet — just the skeleton.

- [x] T001 Create the root-level monorepo directory structure. Create the following empty directories from the repository root: `backend/cmd/server/`, `backend/internal/config/`, `backend/internal/health/`, `backend/internal/logging/`, `frontend/src/pages/`, `frontend/src/services/`, `frontend/public/`, `scripts/`, `.github/workflows/`, `.vscode/`, `.windsurf/`. Use `mkdir -p` for each path. Do NOT create any files yet — only directories.

- [x] T002 [P] Initialize the Go module at `backend/go.mod`. Run `go mod init github.com/flagmanagment/backend` from the `backend/` directory. This creates `backend/go.mod` with module path `github.com/flagmanagment/backend` and Go version `go 1.22`. Do NOT add any dependencies yet — they will be added in later tasks.

- [x] T003 [P] Initialize the frontend project at `frontend/`. Run `npm create vite@latest ./ -- --template react-ts` from the `frontend/` directory to scaffold a React + TypeScript project with Vite. This creates `frontend/package.json`, `frontend/tsconfig.json`, `frontend/vite.config.ts`, `frontend/index.html`, `frontend/src/App.tsx`, `frontend/src/main.tsx`, and related config files. After scaffolding, run `npm install` to generate `frontend/package-lock.json`. Delete the default `frontend/src/App.css` and `frontend/src/index.css` placeholder content (we will replace them later). Keep `frontend/src/vite-env.d.ts` as-is.

- [x] T004 [P] Create the environment variable documentation file at `.env.example` in the repository root. This file documents ALL environment variables used by the project with their default values and comments. Write the following content exactly:

```
# FlagManagment Environment Configuration
# Copy this file to .env and modify as needed.
# All variables have sensible defaults for local development.

# Backend Configuration
FM_BACKEND_PORT=8080          # Port the Go backend listens on
FM_ENV=development            # Environment: development | production
FM_LOG_FORMAT=auto            # Logging format: auto | text | json (auto = text in dev, json in prod)

# Database Configuration
FM_DB_HOST=postgres           # PostgreSQL hostname (Docker service name)
FM_DB_PORT=5432               # PostgreSQL port
FM_DB_USER=flagmgmt           # PostgreSQL username
FM_DB_PASSWORD=flagmgmt_dev   # PostgreSQL password (change in production!)
FM_DB_NAME=flagmanagment      # PostgreSQL database name

# Redis Configuration
FM_REDIS_HOST=redis           # Redis hostname (Docker service name)
FM_REDIS_PORT=6379            # Redis port

# Frontend Configuration
FM_FRONTEND_PORT=3000         # Port the frontend dev server / nginx listens on
```

- [x] T005 [P] Create the root `.gitignore` file at `.gitignore` in the repository root. Include the following patterns: `.env` (not `.env.example`), `backend/tmp/` (air build output), `frontend/node_modules/`, `frontend/dist/`, `*.log`, `.DS_Store`, `coverage/`, `*.coverprofile`. Do NOT ignore `.vscode/` or `.windsurf/` — these are committed workspace configs.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented. These tasks create the Go backend skeleton, configuration loading, structured logging, and the Docker Compose orchestration.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T006 [P] Add backend dependencies. Run the following commands from the `backend/` directory:, run the following commands one by one:
  - `go get github.com/go-chi/chi/v5@latest` — This is the HTTP router. chi v5 is a lightweight, idiomatic Go HTTP router built on `net/http`. We chose it over echo/gin because it uses standard `http.Handler` interfaces (see research.md §1).
  - `go get github.com/rs/zerolog@latest` — This is the structured logging library. We chose zerolog over zap/slog because it provides zero-allocation logging with built-in ConsoleWriter for dev mode and JSON for production (see research.md §2).
  - `go get github.com/lib/pq@latest` — This is the PostgreSQL driver. It provides the `database/sql` driver for PostgreSQL. We only need it for the health check ping in this phase.
  - `go get github.com/redis/go-redis/v9@latest` — This is the Redis client library. We only need it for the health check ping in this phase.
  After running all four commands, run `go mod tidy` to clean up `backend/go.sum`.

- [x] T007 Create the environment configuration loader at `backend/internal/config/config.go`. This file defines a `Config` struct and a `Load()` function that reads ALL environment variables. The struct MUST have the following fields with these exact types and `env` tags. Use `os.Getenv` with manual defaults — do NOT add a third-party env parsing library (keep dependencies minimal for Phase 1).

  **Exact implementation**:
  ```go
  package config

  import "os"

  type Config struct {
      BackendPort string
      Env         string // "development" or "production"
      LogFormat   string // "auto", "text", or "json"
      DBHost      string
      DBPort      string
      DBUser      string
      DBPassword  string
      DBName      string
      RedisHost   string
      RedisPort   string
  }

  func Load() *Config {
      return &Config{
          BackendPort: getEnv("FM_BACKEND_PORT", "8080"),
          Env:         getEnv("FM_ENV", "development"),
          LogFormat:   getEnv("FM_LOG_FORMAT", "auto"),
          DBHost:      getEnv("FM_DB_HOST", "localhost"),
          DBPort:      getEnv("FM_DB_PORT", "5432"),
          DBUser:      getEnv("FM_DB_USER", "flagmgmt"),
          DBPassword:  getEnv("FM_DB_PASSWORD", "flagmgmt_dev"),
          DBName:      getEnv("FM_DB_NAME", "flagmanagment"),
          RedisHost:   getEnv("FM_REDIS_HOST", "localhost"),
          RedisPort:   getEnv("FM_REDIS_PORT", "6379"),
      }
  }

  func getEnv(key, fallback string) string {
      if val := os.LookupEnv(key); val != "" {
          return val
      }
      return fallback
  }
  ```
  Note: Use `os.LookupEnv` not `os.Getenv` so that empty string set explicitly is distinguishable, but for simplicity in this phase, treat empty as "use fallback". The `getEnv` helper handles this.

- [x] T008 Create the structured logger at `backend/internal/logging/logger.go`. This file provides a `NewLogger` function that returns a configured `zerolog.Logger`. It MUST auto-detect the environment:
  - If `FM_LOG_FORMAT` is `"text"` → use `zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}` (human-readable colored output).
  - If `FM_LOG_FORMAT` is `"json"` → use `zerolog.New(os.Stdout).With().Timestamp().Logger()` (machine-readable JSON).
  - If `FM_LOG_FORMAT` is `"auto"` (the default) → check `FM_ENV`: if `"development"` use text, if `"production"` use JSON.

  **Exact implementation**:
  ```go
  package logging

  import (
      "os"
      "time"
      "github.com/rs/zerolog"
  )

  func NewLogger(logFormat, env string) zerolog.Logger {
      format := logFormat
      if format == "auto" {
          if env == "production" {
              format = "json"
          } else {
              format = "text"
          }
      }

      if format == "text" {
          writer := zerolog.ConsoleWriter{
              Out:        os.Stdout,
              TimeFormat: time.RFC3339,
          }
          return zerolog.New(writer).With().Timestamp().Logger()
      }

      return zerolog.New(os.Stdout).With().Timestamp().Logger()
  }
  ```

- [x] T009 Create the health check handler at `backend/internal/health/handler.go`. This is the ONLY API endpoint in Phase 1. It implements the contract defined in `specs/001-core-architecture/contracts/health-check.md`.

  The handler MUST:
  1. Accept `*sql.DB` (PostgreSQL) and `*redis.Client` (Redis) as dependencies via a struct.
  2. Ping PostgreSQL using `db.PingContext(ctx)` with a 2-second timeout.
  3. Ping Redis using `rdb.Ping(ctx)` with a 2-second timeout.
  4. Return HTTP 200 with `{"status": "healthy", ...}` if BOTH pings succeed.
  5. Return HTTP 503 with `{"status": "unhealthy", ...}` if ANY ping fails.
  6. Track uptime by recording `time.Now()` at initialization and computing `uptime_seconds` on each request.
  7. Include individual latency in milliseconds for each dependency check.
  8. The response JSON schema MUST match the contract exactly:
     ```json
     {
       "status": "healthy|unhealthy",
       "version": "0.1.0",
       "uptime_seconds": 1234,
       "checks": {
         "postgres": { "status": "healthy|unhealthy", "latency_ms": 2, "error": "optional" },
         "redis": { "status": "healthy|unhealthy", "latency_ms": 1, "error": "optional" }
       }
     }
     ```
  9. The `version` field should be a package-level `var Version = "0.1.0"` that can be overridden at build time via `ldflags`.

  **Key design decisions**:
  - Use `context.WithTimeout` derived from `r.Context()` for each ping — do NOT use `context.Background()`.
  - Use `encoding/json` from the stdlib for JSON marshaling — do NOT add a third-party JSON library.
  - The handler function signature MUST be `http.HandlerFunc` (compatible with chi router and stdlib).
  - Use a `Handler` struct with a `ServeHTTP` method or a constructor `NewHandler(db, rdb, logger) http.HandlerFunc` — either pattern works, but constructor returning `http.HandlerFunc` is preferred for simplicity.

- [x] T010 Create the application entry point at `backend/cmd/server/main.go`. This file wires everything together and starts the HTTP server.

  It MUST do the following in this exact order:
  1. Call `config.Load()` to get the configuration.
  2. Call `logging.NewLogger(cfg.LogFormat, cfg.Env)` to create the logger.
  3. Log a startup message: `logger.Info().Str("port", cfg.BackendPort).Msg("starting FlagManagment backend")`.
  4. Open a PostgreSQL connection using `sql.Open("postgres", connStr)` where `connStr` is built from config fields: `fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)`. Import `_ "github.com/lib/pq"` for the driver side effect.
  5. Verify the DB connection with `db.Ping()` — log a warning (not fatal) if it fails, because Docker Compose health checks handle retry.
  6. Open a Redis connection using `redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)})`.
  7. Create a `chi.NewRouter()`.
  8. Register `GET /healthz` using the health handler from T009.
  9. Start the HTTP server with `http.ListenAndServe(":"+cfg.BackendPort, router)`.
  10. Log a fatal error if the server fails to start.

  **Important**: Do NOT add graceful shutdown in this phase — it will be added in a later feature. A simple `log.Fatal(http.ListenAndServe(...))` is sufficient.

- [x] T011 Create the golangci-lint configuration at `backend/.golangci.yml`. This file configures the Go linter that CI will run. Use the following configuration:
  ```yaml
  run:
    timeout: 5m
    go: "1.22"

  linters:
    enable:
      - errcheck
      - gosimple
      - govet
      - ineffassign
      - staticcheck
      - unused
      - gofmt
      - goimports

  linters-settings:
    goimports:
      local-prefixes: github.com/flagmanagment

  issues:
    max-issues-per-linter: 0
    max-same-issues: 0
  ```
  This configuration enables the essential linters without being overly aggressive for Phase 1. The `goimports` setting ensures local package imports are grouped separately.

- [x] T012 Create the Air hot-reload configuration at `backend/.air.toml`. This file configures the `cosmtrek/air` tool for Go live reloading during development (see research.md §6). Write the following content:
  ```toml
  root = "."
  tmp_dir = "tmp"

  [build]
    bin = "./tmp/main"
    cmd = "go build -o ./tmp/main ./cmd/server"
    delay = 1000
    exclude_dir = ["tmp", "vendor"]
    exclude_regex = ["_test\\.go$"]
    include_ext = ["go", "toml", "yaml"]
    kill_delay = "0s"
    send_interrupt = false
    stop_on_error = true

  [misc]
    clean_on_exit = true

  [log]
    time = false
  ```
  Also create the `backend/tmp/` directory and add `backend/tmp/` to `.gitignore` (it should already be there from T005).

**Checkpoint**: Foundation ready — all Go code compiles, config loads from env vars, logger auto-detects format, health handler pings DB/Redis. User story implementation can now begin.

---

## Phase 3: User Story 1 - Local Development Bootstrap (Priority: P1) 🎯 MVP

**Goal**: A developer clones the repo, runs one command, and has all 4 services running locally with health checks passing.

**Independent Test**: Clone repo → `make up` → `curl localhost:8080/healthz` returns `{"status": "healthy"}` → dashboard loads at `localhost:3000`.

### Implementation for User Story 1

- [x] T013 [US1] Create the backend Dockerfile at `backend/Dockerfile`. This is a multi-stage Docker build that produces a minimal, multi-arch-ready container image (see research.md §3 and §4).

  **Stage 1 — Build** (named `builder`):
  - Base image: `golang:1.22-bookworm`
  - Set `WORKDIR /app`
  - Copy `go.mod` and `go.sum` first, then run `go mod download` (Docker layer caching for dependencies)
  - Copy the rest of the source code
  - Build with: `CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags="-s -w -X github.com/flagmanagment/backend/internal/health.Version=0.1.0" -o /app/server ./cmd/server`
  - **IMPORTANT**: Use the `$TARGETARCH` build argument which Docker Buildx provides automatically. Do NOT hardcode `amd64`. The `CGO_ENABLED=0` flag is critical — it produces a statically-linked binary that runs on distroless.

  **Stage 2 — Runtime**:
  - Base image: `gcr.io/distroless/static-debian12`
  - Copy the binary from the builder stage: `COPY --from=builder /app/server /server`
  - Set `EXPOSE 8080`
  - Set `ENTRYPOINT ["/server"]`

  The final image should be approximately 10-15MB total.

- [x] T014 [US1] Create the frontend Dockerfile at `frontend/Dockerfile`. This is a multi-stage Docker build for the React dashboard (see research.md §5).

  **Stage 1 — Build** (named `builder`):
  - Base image: `node:20-bookworm-slim`
  - Set `WORKDIR /app`
  - Copy `package.json` and `package-lock.json` first, then run `npm ci` (Docker layer caching)
  - Copy the rest of the source code
  - Run `npm run build` to produce static assets in `/app/dist/`

  **Stage 2 — Serve**:
  - Base image: `nginx:stable-alpine`
  - Copy static assets: `COPY --from=builder /app/dist /usr/share/nginx/html`
  - Copy a custom nginx config (created in T015) to `/etc/nginx/conf.d/default.conf`
  - Set `EXPOSE 3000`

  **IMPORTANT**: The nginx config must listen on port 3000 (not the default 80) to match the `FM_FRONTEND_PORT` default.

- [x] T015 [US1] Create the nginx configuration at `frontend/nginx.conf`. This file configures nginx to serve the React SPA correctly. Key requirements:
  - Listen on port 3000
  - Serve static files from `/usr/share/nginx/html`
  - For any route that doesn't match a file, return `index.html` (SPA client-side routing support via `try_files $uri $uri/ /index.html`)
  - Enable gzip compression for text, CSS, JavaScript, and JSON
  - Set `server_tokens off` for security (don't expose nginx version)
  ```nginx
  server {
      listen 3000;
      server_name _;
      server_tokens off;
      root /usr/share/nginx/html;
      index index.html;

      gzip on;
      gzip_types text/plain text/css application/json application/javascript text/xml;

      location / {
          try_files $uri $uri/ /index.html;
      }
  }
  ```

- [x] T016 [US1] Create the Docker Compose orchestration file at `docker-compose.yml` in the repository root. This file defines ALL 4 services and their relationships. Use Docker Compose v3.8+ syntax (`version` key is optional in modern compose).

  **Services to define**:

  1. **`postgres`**:
     - Image: `postgres:16-alpine`
     - Environment: `POSTGRES_USER=${FM_DB_USER:-flagmgmt}`, `POSTGRES_PASSWORD=${FM_DB_PASSWORD:-flagmgmt_dev}`, `POSTGRES_DB=${FM_DB_NAME:-flagmanagment}`
     - Ports: `${FM_DB_PORT:-5432}:5432`
     - Volumes: `pgdata:/var/lib/postgresql/data` (named volume for persistence)
     - Healthcheck: `test: ["CMD-SHELL", "pg_isready -U ${FM_DB_USER:-flagmgmt} -d ${FM_DB_NAME:-flagmanagment}"]`, interval 10s, timeout 5s, retries 5

  2. **`redis`**:
     - Image: `redis:7-alpine`
     - Ports: `${FM_REDIS_PORT:-6379}:6379`
     - Volumes: `redisdata:/data` (named volume)
     - Healthcheck: `test: ["CMD", "redis-cli", "ping"]`, interval 10s, timeout 5s, retries 5

  3. **`backend`**:
     - Build context: `./backend` with Dockerfile `./backend/Dockerfile`
     - Ports: `${FM_BACKEND_PORT:-8080}:8080`
     - Environment: pass ALL `FM_*` variables (DB, Redis, logging config). IMPORTANT: override `FM_DB_HOST=postgres` and `FM_REDIS_HOST=redis` so the backend connects to Docker's internal DNS names, not localhost.
     - `depends_on`: postgres (condition: service_healthy), redis (condition: service_healthy)
     - Healthcheck: `test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/healthz"]`, interval 10s, timeout 5s, retries 3. Note: `wget` is available in distroless — if not, use a different approach like a Go-based health probe or `curl` added to the build. **Decision**: Since distroless has no shell or wget, use an alternative: change the backend to include a small `/healthz` self-check, OR use `docker compose exec` based check. **Simplest solution**: Skip the container-internal healthcheck for the backend and use `depends_on` only for postgres/redis. The external health check (`curl localhost:8080/healthz`) verifies the backend from the host.

  4. **`frontend`**:
     - Build context: `./frontend` with Dockerfile `./frontend/Dockerfile`
     - Ports: `${FM_FRONTEND_PORT:-3000}:3000`
     - `depends_on`: backend (no health condition needed — frontend works independently)

  **Volumes** (at bottom of file):
  ```yaml
  volumes:
    pgdata:
    redisdata:
  ```

  **Networks**: Use the default Docker Compose network — all services on the same bridge network can communicate by service name.

- [x] T017 [US1] Create the Docker Compose override file for local development at `docker-compose.override.yml` in the repository root. This file is automatically loaded by `docker compose up` and adds development-specific configuration on top of `docker-compose.yml`.

  **Overrides for `backend`**:
  - Replace the `build` with a development image that includes `air` for hot-reload
  - Mount the source code: `volumes: ["./backend:/app"]`
  - Override the entrypoint to use `air` instead of the compiled binary: `command: ["air", "-c", ".air.toml"]`
  - Use `golang:1.22-bookworm` as the base image instead of building the Dockerfile (for development, we need the full Go toolchain with air)
  - **Decision**: Instead of modifying the Dockerfile, use a separate `backend/Dockerfile.dev` for development:
    ```dockerfile
    FROM golang:1.22-bookworm
    WORKDIR /app
    RUN go install github.com/air-verse/air@latest
    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .
    CMD ["air", "-c", ".air.toml"]
    ```

  **Overrides for `frontend`**:
  - Mount the source code: `volumes: ["./frontend:/app", "/app/node_modules"]` (the anonymous volume `/app/node_modules` prevents the host's node_modules from overwriting the container's)
  - Override the command to use Vite dev server instead of nginx: `command: ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "3000"]`
  - Use `node:20-bookworm-slim` as the base image for dev mode

  **Important**: The override file MUST use `image:` or `build: dockerfile:` to select the dev variants. The exact approach:
  ```yaml
  services:
    backend:
      build:
        context: ./backend
        dockerfile: Dockerfile.dev
      volumes:
        - ./backend:/app
      command: ["air", "-c", ".air.toml"]

    frontend:
      build:
        context: ./frontend
        dockerfile: Dockerfile
      image: node:20-bookworm-slim
      volumes:
        - ./frontend:/app
        - /app/node_modules
      command: ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "3000"]
  ```

- [x] T018 [US1] Create the backend development Dockerfile at `backend/Dockerfile.dev`. This Dockerfile is used ONLY during local development (via `docker-compose.override.yml`) to enable hot-reload with `air`. It is NOT used in production or CI.

  ```dockerfile
  FROM golang:1.22-bookworm

  WORKDIR /app

  # Install air for hot-reload
  RUN go install github.com/air-verse/air@latest

  # Pre-download dependencies (cached layer)
  COPY go.mod go.sum ./
  RUN go mod download

  # Copy source (will be overridden by volume mount in docker-compose.override.yml)
  COPY . .

  # Air will build and run the server
  CMD ["air", "-c", ".air.toml"]
  ```

- [x] T019 [US1] Create the bootstrap script at `scripts/bootstrap.sh`. This script validates prerequisites and then runs `docker compose up` (see research.md §8). It MUST be executable (`chmod +x`).

  **The script must perform these checks in order**:
  1. Check if `docker` command exists → if not, print: `"ERROR: Docker is not installed. Install it from https://docs.docker.com/get-docker/"` and exit 1.
  2. Check if `docker compose` subcommand works (run `docker compose version`) → if not, print: `"ERROR: Docker Compose v2 is required. Update Docker Desktop or install the compose plugin."` and exit 1.
  3. Check if Docker daemon is running (run `docker info > /dev/null 2>&1`) → if not, print: `"ERROR: Docker daemon is not running. Start Docker Desktop or run 'sudo systemctl start docker'."` and exit 1.
  4. Check Docker version is >= 24.0 (parse `docker version --format '{{.Server.Version}}'`) → if below, print a WARNING (not error): `"WARNING: Docker 24.0+ recommended. You have X.Y.Z. Some features may not work."`. Continue anyway.
  5. Check port availability for ports 8080, 3000, 5432, 6379 using `lsof -i :PORT` or `ss -tlnp | grep :PORT` → if any port is in use, print: `"ERROR: Port XXXX is already in use. Either stop the conflicting service or set FM_BACKEND_PORT/FM_FRONTEND_PORT/FM_DB_PORT/FM_REDIS_PORT in your .env file."` and exit 1.
  6. Check architecture (run `uname -m`) → if it returns something other than `x86_64` or `aarch64`, print a WARNING: `"WARNING: Architecture $(uname -m) is not officially supported. Supported: x86_64, aarch64 (ARM64)."`.
  7. If `.env` file does not exist, copy `.env.example` to `.env` and print: `"INFO: Created .env from .env.example with default values."`.
  8. Print: `"Starting FlagManagment services..."`.
  9. Run: `docker compose up --build -d`.
  10. Print: `"Waiting for services to become healthy..."`.
  11. Wait up to 120 seconds for the backend health check (loop: `curl -sf http://localhost:${FM_BACKEND_PORT:-8080}/healthz > /dev/null 2>&1` every 5 seconds). If it succeeds, print the health check JSON response and a success message. If timeout, print: `"WARNING: Backend health check did not pass within 120 seconds. Check logs with: docker compose logs backend"`.
  12. Print the final summary:
      ```
      ✅ FlagManagment is running!
         Backend:    http://localhost:8080/healthz
         Dashboard:  http://localhost:3000
         PostgreSQL: localhost:5432
         Redis:      localhost:6379
         Logs:       docker compose logs -f
         Stop:       make down
      ```

- [x] T020 [US1] Create the Makefile at `Makefile` in the repository root. This provides convenience wrappers for common development commands.

  ```makefile
  .PHONY: up down test lint build logs clean

  # Start all services (with prerequisite validation)
  up:
  	@./scripts/bootstrap.sh

  # Stop all services
  down:
  	docker compose down

  # Stop and remove all data volumes (clean slate)
  clean:
  	docker compose down -v

  # Run all tests
  test: test-backend test-frontend

  test-backend:
  	cd backend && go test -race -cover -coverprofile=coverage.out ./...

  test-frontend:
  	cd frontend && npm test -- --coverage

  # Run all linters
  lint: lint-backend lint-frontend

  lint-backend:
  	cd backend && golangci-lint run ./...

  lint-frontend:
  	cd frontend && npx eslint src/ --ext .ts,.tsx
  	cd frontend && npx prettier --check "src/**/*.{ts,tsx,css}"

  # Multi-architecture Docker build (does not push)
  build:
  	docker buildx build --platform linux/amd64,linux/arm64 -t flagmanagment/backend:latest -f backend/Dockerfile backend/
  	docker buildx build --platform linux/amd64,linux/arm64 -t flagmanagment/frontend:latest -f frontend/Dockerfile frontend/

  # View logs
  logs:
  	docker compose logs -f

  # Check health
  health:
  	@curl -s http://localhost:8080/healthz | python3 -m json.tool 2>/dev/null || echo "Backend not running"
  ```

  **IMPORTANT**: Makefile rules MUST use TAB characters for indentation, not spaces. This is a Makefile syntax requirement — using spaces will cause `make` to fail with a confusing error.

- [x] T021 Create the health check handler tests at `backend/internal/health/handler_test.go`. Test both healthy and unhealthy scenarios using mock database and redis connections. You can use `github.com/DATA-DOG/go-sqlmock` and `github.com/go-redis/redismock/v9`.

- [x] T022 Setup the initial React entry points.
  - Overwrite `frontend/src/main.tsx` to just import and render `App`.
  - Overwrite `frontend/src/App.tsx` to mount the `Dashboard` component.

- [x] T023 Create the frontend API service at `frontend/src/services/api.ts`.
  - Export TypeScript interfaces matching the backend health check JSON (`HealthResponse`, `CheckResult`).
  - Create a `fetchHealth()` async function that calls `/api/healthz`. (Nginx or Vite proxy will handle the `/api` prefix).

- [x] T024 Create the dashboard layout at `frontend/src/pages/Dashboard.tsx`.
  - A simple React component with a header ("FlagManagment Dashboard") and a main content area.
  - Include the `HealthStatus` component in the main area.

- [x] T025 Create the health status component at `frontend/src/components/HealthStatus.tsx`.
  - Use `useEffect` and `useState` to call `fetchHealth()` on mount and every 10 seconds.
  - Render a visual indicator (green/red dot) for overall status.
  - Render a grid of cards showing Postgres and Redis status, latencies, and backend uptime.
  - Show the overall status as a large colored badge (green for healthy, red for unhealthy)
  - Show each dependency (PostgreSQL, Redis) with its status and latency
  - Show the backend version and uptime
  - Show a "Last checked" timestamp
  - Use inline styles or a small CSS module — do NOT install a UI library in this phase

  **File 3 — Modify `frontend/src/App.tsx`**: Replace the default Vite scaffold content with a simple wrapper that renders `<HealthDashboard />`. Remove all Vite demo code (counter, logos, etc.).

  **File 4 — Modify `frontend/src/main.tsx`**: Keep this as-is from the Vite scaffold (it renders `<App />` into `#root`). Just ensure `import './index.css'` is present if you have global styles, or remove it if not.

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. Running `make up` should bootstrap all services, `curl localhost:8080/healthz` should return healthy, and `localhost:3000` should show the health dashboard.

---

## Phase 4: User Story 2 - Multi-Architecture Container Builds (Priority: P2)

**Goal**: Container images build and run on x86_64 and ARM64 (including Raspberry Pi 4).

**Independent Test**: Run `make build` → verify images exist for both architectures → run on ARM64 machine if available.

### Implementation for User Story 2

- [x] T022 [P] [US2] Verify and update the backend Dockerfile at `backend/Dockerfile` to ensure multi-arch compatibility. The Dockerfile created in T013 already uses `$TARGETARCH`, but verify these specific points:
  1. The `FROM` line for the builder stage does NOT pin a specific platform — Docker Buildx handles this.
  2. The `go build` command uses `GOARCH=$TARGETARCH` (not hardcoded `amd64`).
  3. The runtime stage `FROM gcr.io/distroless/static-debian12` also does NOT pin a platform.
  4. Add a build arg declaration at the top of the builder stage: `ARG TARGETARCH` — this is automatically provided by Buildx but must be declared if used in RUN commands.
  5. Verify the binary is statically linked: `CGO_ENABLED=0` must be set.

  If all points are already correct from T013, mark this task as verified/complete with no changes needed.

- [x] T023 [P] [US2] Verify and update the frontend Dockerfile at `frontend/Dockerfile` to ensure multi-arch compatibility. Verify:
  1. `node:20-bookworm-slim` supports ARM64 natively (it does — Node.js publishes multi-arch manifests).
  2. `nginx:stable-alpine` supports ARM64 natively (it does).
  3. No architecture-specific commands in the build steps (e.g., no `--arch` flags in npm commands).
  4. The `npm ci` command works identically on both architectures.

  If all points are already correct from T014, mark this task as verified/complete with no changes needed.

- [x] T024 [US2] Update the `Makefile` `build` target (already created in T020) to ensure it creates a Buildx builder instance if one doesn't exist. Add a `build-setup` target:
  ```makefile
  build-setup:
  	@docker buildx inspect flagmanagment-builder > /dev/null 2>&1 || \
  	  docker buildx create --name flagmanagment-builder --use --driver docker-container
  	@docker buildx use flagmanagment-builder

  build: build-setup
  	docker buildx build --platform linux/amd64,linux/arm64 -t flagmanagment/backend:latest -f backend/Dockerfile backend/
  	docker buildx build --platform linux/amd64,linux/arm64 -t flagmanagment/frontend:latest -f frontend/Dockerfile frontend/
  ```
  The `docker-container` driver is required for multi-platform builds — the default `docker` driver only builds for the host platform.

**Checkpoint**: Multi-arch images build successfully for both linux/amd64 and linux/arm64. `docker buildx imagetools inspect` shows both platforms in the manifest.

---

## Phase 5: User Story 3 - CI/CD Pipeline Scaffolding (Priority: P3)

**Goal**: Every pull request automatically gets linted, tested, and built. Main branch merges publish images to private ghcr.io.

**Independent Test**: Open a test PR with clean code → CI passes. Open a PR with lint violations → CI fails with specific errors.

### Implementation for User Story 3

- [x] T025 [P] [US3] Create the CI workflow at `.github/workflows/ci.yml`. This GitHub Actions workflow runs on every pull request to the `main` branch.

  **Workflow structure**:
  ```yaml
  name: CI

  on:
    pull_request:
      branches: [main]

  jobs:
    lint-backend:
      runs-on: ubuntu-latest
      steps:
        - uses: actions/checkout@v4
        - uses: actions/setup-go@v5
          with:
            go-version: '1.22'
        - name: golangci-lint
          uses: golangci/golangci-lint-action@v6
          with:
            version: latest
            working-directory: backend

    lint-frontend:
      runs-on: ubuntu-latest
      steps:
        - uses: actions/checkout@v4
        - uses: actions/setup-node@v4
          with:
            node-version: '20'
            cache: 'npm'
            cache-dependency-path: frontend/package-lock.json
        - run: cd frontend && npm ci
        - run: cd frontend && npx eslint src/ --ext .ts,.tsx
        - run: cd frontend && npx prettier --check "src/**/*.{ts,tsx,css}"

    test-backend:
      runs-on: ubuntu-latest
      needs: lint-backend
      steps:
        - uses: actions/checkout@v4
        - uses: actions/setup-go@v5
          with:
            go-version: '1.22'
        - name: Run tests with coverage
          run: cd backend && go test -race -coverprofile=coverage.out -covermode=atomic ./...
        - name: Check coverage threshold
          run: |
            cd backend
            COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
            echo "Backend coverage: ${COVERAGE}%"
            if (( $(echo "$COVERAGE < 80" | bc -l) )); then
              echo "::error::Backend coverage ${COVERAGE}% is below the 80% threshold"
              exit 1
            fi

    test-frontend:
      runs-on: ubuntu-latest
      needs: lint-frontend
      steps:
        - uses: actions/checkout@v4
        - uses: actions/setup-node@v4
          with:
            node-version: '20'
            cache: 'npm'
            cache-dependency-path: frontend/package-lock.json
        - run: cd frontend && npm ci
        - name: Run tests with coverage
          run: cd frontend && npx vitest run --coverage
          # Note: Coverage threshold is enforced in vitest.config.ts (see T027)

    build-verify:
      runs-on: ubuntu-latest
      needs: [test-backend, test-frontend]
      steps:
        - uses: actions/checkout@v4
        - name: Build backend image
          run: docker build -t flagmanagment/backend:ci -f backend/Dockerfile backend/
        - name: Build frontend image
          run: docker build -t flagmanagment/frontend:ci -f frontend/Dockerfile frontend/
  ```

  **Key decisions**:
  - Jobs are structured as a pipeline: lint → test → build. This ensures fast failure (lint is fastest).
  - Backend lint and frontend lint run in parallel (separate jobs).
  - The `build-verify` job only does a single-arch build (host arch) to verify the Dockerfile works. Multi-arch is handled by the publish workflow.
  - Coverage threshold is checked in the CI script (80% for backend) and in vitest config (70% for frontend).

- [x] T026 [P] [US3] Create the publish workflow at `.github/workflows/publish.yml`. This workflow runs when code is pushed to `main` (after a PR merge) and publishes multi-arch images to the private GitHub Container Registry.

  ```yaml
  name: Publish

  on:
    push:
      branches: [main]

  env:
    REGISTRY: ghcr.io
    BACKEND_IMAGE: ghcr.io/${{ github.repository_owner }}/flagmanagment-backend
    FRONTEND_IMAGE: ghcr.io/${{ github.repository_owner }}/flagmanagment-frontend

  jobs:
    publish:
      runs-on: ubuntu-latest
      permissions:
        contents: read
        packages: write

      steps:
        - uses: actions/checkout@v4

        - name: Log in to GitHub Container Registry
          uses: docker/login-action@v3
          with:
            registry: ${{ env.REGISTRY }}
            username: ${{ github.actor }}
            password: ${{ secrets.GITHUB_TOKEN }}

        - name: Set up Docker Buildx
          uses: docker/setup-buildx-action@v3

        - name: Build and push backend
          uses: docker/build-push-action@v6
          with:
            context: ./backend
            file: ./backend/Dockerfile
            platforms: linux/amd64,linux/arm64
            push: true
            tags: |
              ${{ env.BACKEND_IMAGE }}:latest
              ${{ env.BACKEND_IMAGE }}:${{ github.sha }}

        - name: Build and push frontend
          uses: docker/build-push-action@v6
          with:
            context: ./frontend
            file: ./frontend/Dockerfile
            platforms: linux/amd64,linux/arm64
            push: true
            tags: |
              ${{ env.FRONTEND_IMAGE }}:latest
              ${{ env.FRONTEND_IMAGE }}:${{ github.sha }}
  ```

  **Key decisions**:
  - Uses `GITHUB_TOKEN` (automatically available) — no need to create a PAT.
  - `permissions.packages: write` is required to push to ghcr.io.
  - Images are tagged with both `latest` and the commit SHA for traceability.
  - Uses the official `docker/build-push-action` which handles Buildx setup automatically.

- [x] T027 [US3] Configure Vitest with coverage threshold enforcement at `frontend/vitest.config.ts`. This file configures the test runner and enforces the 70% frontend coverage minimum. If the file already exists from T003 (Vite scaffold), update it. Otherwise create it:

  ```typescript
  import { defineConfig } from 'vitest/config';
  import react from '@vitejs/plugin-react';

  export default defineConfig({
    plugins: [react()],
    test: {
      globals: true,
      environment: 'jsdom',
      setupFiles: './src/test-setup.ts',
      coverage: {
        provider: 'v8',
        reporter: ['text', 'lcov'],
        thresholds: {
          statements: 70,
          branches: 70,
          functions: 70,
          lines: 70,
        },
      },
    },
  });
  ```

  Also create `frontend/src/test-setup.ts` with:
  ```typescript
  import '@testing-library/jest-dom';
  ```

  Install required dev dependencies from the `frontend/` directory:
  ```bash
  npm install --save-dev vitest @vitest/coverage-v8 jsdom @testing-library/react @testing-library/jest-dom @testing-library/user-event
  ```

- [x] T028 [US3] Configure ESLint for the frontend at `frontend/.eslintrc.cjs` (or update if it exists from Vite scaffold). Ensure it includes TypeScript and React rules:

  ```javascript
  module.exports = {
    root: true,
    env: { browser: true, es2020: true },
    extends: [
      'eslint:recommended',
      'plugin:@typescript-eslint/recommended',
      'plugin:react-hooks/recommended',
    ],
    ignorePatterns: ['dist', '.eslintrc.cjs'],
    parser: '@typescript-eslint/parser',
    plugins: ['react-refresh'],
    rules: {
      'react-refresh/only-export-components': ['warn', { allowConstantExport: true }],
      '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
    },
  };
  ```

  Install Prettier as a dev dependency and create `frontend/.prettierrc`:
  ```json
  {
    "semi": true,
    "singleQuote": true,
    "trailingComma": "all",
    "printWidth": 100,
    "tabWidth": 2
  }
  ```

  Run: `cd frontend && npm install --save-dev prettier`

- [x] T029 [US3] Create the CODEOWNERS file at `.github/CODEOWNERS`. This enforces that at least one code owner reviews every PR. For Phase 1, assign the repository owner:
  ```
  # Global code owners — all PRs require review from at least one owner
  * @flagmanagment/core-team
  ```
  Note: Update `@flagmanagment/core-team` with the actual GitHub team or user handle.

**Checkpoint**: CI pipeline runs on PRs (lint, test, coverage check, build verify). Publish workflow pushes multi-arch images to ghcr.io on main merges.

---

## Phase 6: User Story 4 - IDE Workspace Configuration (Priority: P4)

**Goal**: Contributors open the project in VS Code or Windsurf and immediately have linting, formatting, and debugging configured.

**Independent Test**: Open the project in a fresh VS Code → it suggests recommended extensions → formatting works on save → Go and TypeScript linting errors appear inline.

### Implementation for User Story 4

- [x] T030 [P] [US4] Create VS Code recommended extensions at `.vscode/extensions.json`:
  ```json
  {
    "recommendations": [
      "golang.go",
      "dbaeumer.vscode-eslint",
      "esbenp.prettier-vscode",
      "ms-azuretools.vscode-docker",
      "redhat.vscode-yaml",
      "bradlc.vscode-tailwindcss",
      "streetsidesoftware.code-spell-checker",
      "eamodio.gitlens"
    ]
  }
  ```
  These extensions cover: Go development, ESLint, Prettier, Docker, YAML editing, spell checking, and Git history.

- [x] T031 [P] [US4] Create VS Code workspace settings at `.vscode/settings.json`:
  ```json
  {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode",
    "editor.codeActionsOnSave": {
      "source.fixAll.eslint": "explicit"
    },
    "[go]": {
      "editor.defaultFormatter": "golang.go",
      "editor.codeActionsOnSave": {
        "source.organizeImports": "explicit"
      }
    },
    "[typescript]": {
      "editor.defaultFormatter": "esbenp.prettier-vscode"
    },
    "[typescriptreact]": {
      "editor.defaultFormatter": "esbenp.prettier-vscode"
    },
    "go.lintTool": "golangci-lint",
    "go.lintFlags": ["--config=backend/.golangci.yml"],
    "go.testFlags": ["-race", "-cover"],
    "typescript.tsdk": "frontend/node_modules/typescript/lib",
    "files.trimTrailingWhitespace": true,
    "files.insertFinalNewline": true,
    "files.exclude": {
      "**/node_modules": true,
      "backend/tmp": true
    }
  }
  ```

- [x] T032 [P] [US4] Create VS Code debug configurations at `.vscode/launch.json`:
  ```json
  {
    "version": "0.2.0",
    "configurations": [
      {
        "name": "Debug Backend (Go)",
        "type": "go",
        "request": "launch",
        "mode": "auto",
        "program": "${workspaceFolder}/backend/cmd/server",
        "env": {
          "FM_DB_HOST": "localhost",
          "FM_REDIS_HOST": "localhost",
          "FM_ENV": "development",
          "FM_LOG_FORMAT": "text"
        },
        "cwd": "${workspaceFolder}/backend"
      },
      {
        "name": "Debug Backend Tests (Go)",
        "type": "go",
        "request": "launch",
        "mode": "test",
        "program": "${workspaceFolder}/backend/...",
        "cwd": "${workspaceFolder}/backend"
      }
    ]
  }
  ```
  Note: The debug configurations assume PostgreSQL and Redis are running locally (e.g., via `docker compose up postgres redis`). The `FM_DB_HOST` and `FM_REDIS_HOST` are set to `localhost` to connect directly rather than through Docker networking.

- [x] T033 [P] [US4] Create Windsurf workspace configuration at `.windsurf/settings.json`:
  ```json
  {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode",
    "[go]": {
      "editor.defaultFormatter": "golang.go"
    },
    "go.lintTool": "golangci-lint",
    "files.trimTrailingWhitespace": true,
    "files.insertFinalNewline": true
  }
  ```
  Note: Windsurf settings follow the same VS Code settings format. Only include the essential settings — Windsurf does not support all VS Code extension IDs.

**Checkpoint**: All user stories should now be independently functional. Contributors get full IDE support out of the box.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, cleanup, and final validation across all user stories.

- [x] T034 Create the project README at `README.md` in the repository root. Include:
  1. Project name and one-line description: "FlagManagment — Cloud-native feature flag management and remote configuration platform."
  2. Prerequisites section (Docker 24+, Git)
  3. Quick Start section with `make up` and health check verification
  4. Development section (hot-reload, running tests, running linters)
  5. Architecture overview (4 services: backend Go, frontend React, PostgreSQL, Redis)
  6. Link to the spec docs: `specs/001-core-architecture/`
  7. License placeholder: "BSL / Fair-Source (TBD)"

- [x] T035 [P] Add a basic unit test for the health check handler at `backend/internal/health/handler_test.go`. This test verifies the handler returns correct JSON when dependencies are healthy and unhealthy.

  **Test 1 — Healthy response**: Mock a healthy `*sql.DB` (use `sqlmock` or simply create a test database connection) and a healthy Redis client. Call the handler with `httptest.NewRecorder`. Assert: HTTP 200, body contains `"status":"healthy"`, both checks have `"status":"healthy"`.

  **Test 2 — Unhealthy response**: Use a `*sql.DB` pointing to a non-existent host and a Redis client with bad address. Call the handler. Assert: HTTP 503, body contains `"status":"unhealthy"`, at least one check has `"status":"unhealthy"` with an error message.

  **Decision**: For Phase 1, use `DATA-DOG/go-sqlmock` to mock the database connection. Install it: `cd backend && go get github.com/DATA-DOG/go-sqlmock@latest`. This avoids needing a real PostgreSQL for unit tests.

  For Redis mocking, use `go-redis/redismock`. Install: `cd backend && go get github.com/go-redis/redismock/v9@latest`.

- [x] T036 [P] Add a basic unit test for the config loader at `backend/internal/config/config_test.go`. Test:
  1. Default values: Call `Load()` without setting any env vars → assert all fields match the defaults from T007.
  2. Override values: Set `FM_BACKEND_PORT=9090` and `FM_DB_NAME=testdb` via `t.Setenv()` → call `Load()` → assert those fields are overridden while others keep defaults.

- [x] T037 [P] Add a basic unit test for the logger at `backend/internal/logging/logger_test.go`. Test:
  1. Text mode: Call `NewLogger("text", "development")` → assert the logger is created without panic (basic smoke test).
  2. JSON mode: Call `NewLogger("json", "production")` → assert the logger writes valid JSON to a buffer.
  3. Auto mode: Call `NewLogger("auto", "development")` → assert it resolves to text. Call `NewLogger("auto", "production")` → assert it resolves to JSON.

  **Implementation tip**: To test logger output, pass a `bytes.Buffer` as the writer. However, since our `NewLogger` currently writes to `os.Stdout`, you'll need to modify the function signature to accept an `io.Writer` parameter (or modify the test to capture stdout). **Decision**: Update `NewLogger` to accept an optional `io.Writer` parameter: `func NewLogger(logFormat, env string, writers ...io.Writer) zerolog.Logger`. If no writer is provided, use `os.Stdout`. This makes it testable without modifying stdout.

- [x] T038 Run the full quickstart validation guide at `specs/001-core-architecture/quickstart.md`. Execute all 5 scenarios manually:
  1. Local Development Bootstrap → verify `make up` works end-to-end
  2. Multi-Architecture Build → verify `make build` succeeds
  3. CI Pipeline Verification → verify `make lint` and `make test` work locally
  4. Hot-Reload Verification → modify a Go file and verify the change appears within 5 seconds
  5. Resource Consumption Check → verify backend uses <250MB RAM via `docker stats`
  Document any issues found and fix them before marking complete.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 completion — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Phase 2 completion — P1 MVP
- **User Story 2 (Phase 4)**: Depends on T013 and T014 (Dockerfiles from US1)
- **User Story 3 (Phase 5)**: Depends on US1 and US2 (needs working Dockerfiles and tests)
- **User Story 4 (Phase 6)**: Can start after Phase 2 (independent of other stories)
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Requires Phase 2 complete. No dependencies on other stories.
- **User Story 2 (P2)**: Requires T013 (backend Dockerfile) and T014 (frontend Dockerfile) from US1. Can otherwise run in parallel.
- **User Story 3 (P3)**: Requires T013, T014 (Dockerfiles), and working tests. Best done after US1+US2.
- **User Story 4 (P4)**: Fully independent. Can start after Phase 2.

### Within Each User Story

- Dockerfiles before Docker Compose
- Docker Compose before bootstrap script
- Bootstrap script before Makefile convenience wrappers
- All infrastructure before the frontend landing page

### Parallel Opportunities

- T002, T003, T004, T005 can all run in parallel (Phase 1 — independent files)
- T022, T023 can run in parallel (Phase 4 — different Dockerfiles)
- T025, T026 can run in parallel (Phase 5 — different workflow files)
- T030, T031, T032, T033 can all run in parallel (Phase 6 — different config files)
- T035, T036, T037 can all run in parallel (Phase 7 — different test files)

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T005)
2. Complete Phase 2: Foundational (T006-T012)
3. Complete Phase 3: User Story 1 (T013-T021)
4. **STOP and VALIDATE**: Run `make up` → `curl localhost:8080/healthz` → open `localhost:3000`
5. This is a deployable MVP

### Incremental Delivery

1. Setup + Foundational → Infrastructure ready
2. Add User Story 1 → Test independently → MVP! (working local dev)
3. Add User Story 2 → Test independently → Multi-arch builds work
4. Add User Story 3 → Test independently → CI/CD pipeline active
5. Add User Story 4 → Test independently → IDE configs committed
6. Polish → Tests, README, final validation

### Parallel Team Strategy

With multiple developers:

1. Everyone completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Docker + Compose + Bootstrap)
   - Developer B: User Story 4 (IDE configs — fully independent)
3. After US1 complete:
   - Developer A: User Story 2 (multi-arch verification)
   - Developer C: User Story 3 (CI/CD workflows)
4. Polish phase: All developers contribute tests and docs

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Every task includes the exact file path and implementation details so a code-generating LLM can execute without asking clarifying questions
- Technology decisions are embedded directly in tasks (router=chi, logger=zerolog, images=distroless, hot-reload=air) — do not second-guess these choices during implementation
