# Implementation Plan: FlagManagment Core Architecture & Repository Bootstrap

**Branch**: `001-core-architecture` | **Date**: 2026-07-18 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/001-core-architecture/spec.md`

## Summary

Bootstrap the FlagManagment platform as a monorepo containing four independent services — a Go backend engine, a React+TypeScript dashboard, PostgreSQL 16+ datastore, and Redis 7+ cache — orchestrated via Docker Compose with multi-architecture builds (x86_64, ARM64, Raspberry Pi 4). Establish CI/CD pipelines on GitHub Actions for automated linting, testing, coverage enforcement, and image publishing to a private GitHub Container Registry. Configure IDE workspaces for VS Code and Windsurf. All services bootstrap with a single `docker compose up` command using environment-driven configuration with auto-detecting structured logging (text locally, JSON in production).

## Technical Context

**Language/Version**: Go 1.22+ (backend), TypeScript 5.4+ / Node 20+ (frontend)

**Primary Dependencies**: chi or echo (Go HTTP router), zerolog or zap (structured logging), Vite 5+ (frontend build), React 18+ (frontend UI)

**Storage**: PostgreSQL 16+ (primary datastore), Redis 7+ (caching & pub/sub)

**Testing**: `go test` with `-race -cover` (backend), Vitest + React Testing Library (frontend), golangci-lint (Go linting), ESLint + Prettier (TypeScript linting/formatting)

**Target Platform**: Linux server (x86_64), macOS (ARM64 Apple Silicon), Raspberry Pi 4 (ARM64), containerized via Docker

**Project Type**: Web service (monorepo: backend API + frontend SPA + infrastructure)

**Performance Goals**: Health check response <1s, hot-reload <5s, backend engine <250MB RAM under standard load

**Constraints**: <250MB RAM for backend engine on Raspberry Pi 4, single-command bootstrap (<10 min), multi-arch Docker images, 80% backend / 70% frontend test coverage

**Scale/Scope**: Foundation for a platform serving enterprise feature flag management — Phase 1 establishes skeleton only, no business logic

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Status | Notes |
|---|-----------|--------|-------|
| I | API-First Contract Design | ✅ PASS | Health check endpoint contract defined. No business APIs at this phase — deferred to Feature 3 (API Contracts). |
| II | Environment Isolation | ⏭️ N/A | No environments, flags, or SDK tokens at this phase. |
| III | Governance by Default | ⏭️ N/A | No RBAC, change requests, or audit logs at this phase. |
| IV | Local Evaluation Performance | ⏭️ N/A | No SDK evaluation engine at this phase. |
| V | Test-First Quality Gates | ✅ PASS | CI pipeline enforces 80% backend / 70% frontend coverage thresholds (FR-009). golangci-lint + ESLint configured. |
| VI | OpenFeature Interoperability | ⏭️ N/A | No SDKs at this phase. |
| VII | PII Protection & Compliance | ⏭️ N/A | No user data, identity, or audit logs at this phase. |
| VIII | Cloud-Native Portability | ✅ PASS | Docker Compose single-command bootstrap, multi-arch builds (x86_64 + ARM64 + RPi4), environment-driven config, Kubernetes manifests planned. |

**Gate Result**: ✅ PASS — All applicable principles satisfied. 3 principles apply to this foundation phase; 5 are deferred to subsequent features.

## Project Structure

### Documentation (this feature)

```text
specs/001-core-architecture/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── health-check.md  # Health check endpoint contract
└── tasks.md             # Phase 2 output (/speckit-tasks)
```

### Source Code (repository root)

```text
backend/
├── cmd/
│   └── server/
│       └── main.go               # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go             # Environment variable loading
│   ├── health/
│   │   └── handler.go            # Health check endpoint handler
│   └── logging/
│       └── logger.go             # Structured logging (text/JSON auto-detect)
├── go.mod
├── go.sum
├── Dockerfile                    # Multi-stage, multi-arch build
└── .golangci.yml                 # Linter configuration

frontend/
├── src/
│   ├── App.tsx                   # Root component
│   ├── main.tsx                  # Entry point
│   ├── pages/
│   │   └── HealthDashboard.tsx   # Minimal landing page with health status
│   └── services/
│       └── api.ts                # Backend API client (health check)
├── public/
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
├── vitest.config.ts
├── Dockerfile                    # Multi-stage, multi-arch build (Nginx serve)
└── .eslintrc.cjs                 # Linter configuration

docker-compose.yml                # Full stack orchestration
docker-compose.override.yml       # Local dev overrides (hot-reload, ports)
.env.example                      # Documented environment variable defaults

.github/
├── workflows/
│   ├── ci.yml                    # PR checks: lint, test, coverage, build
│   └── publish.yml               # Main branch: multi-arch build + ghcr.io push
└── CODEOWNERS                    # Review enforcement

.vscode/
├── extensions.json               # Recommended extensions
├── settings.json                 # Formatting, linting, editor config
└── launch.json                   # Debug configurations

.windsurf/
└── settings.json                 # Windsurf workspace config

Makefile                          # Common development commands
scripts/
└── bootstrap.sh                  # Prerequisite validation + docker compose up
```

**Structure Decision**: Monorepo with top-level `backend/` and `frontend/` directories. This is a web application architecture with a Go API backend and React SPA frontend, each independently containerized but orchestrated together. Infrastructure configs (`docker-compose.yml`, `.github/`, `.vscode/`) live at the repository root.

## Complexity Tracking

> No constitution violations to justify. All applicable principles pass.
