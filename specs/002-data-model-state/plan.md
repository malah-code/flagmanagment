# Implementation Plan: Data Model & State Management

**Branch**: `002-data-model-state` | **Date**: 2026-07-20 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/002-data-model-state/spec.md`

## Summary

Design and implement the complete PostgreSQL schema for FlagManagment covering core entities (projects, environments, feature flags, environment flag states), governance tables (change requests, approvals, audit logs), and RBAC tables (roles, user roles). The schema is managed via versioned SQL migration files using `golang-migrate`. Includes JSONB schema definitions for targeting rules and remote config payloads, performance-oriented indexing, and a Go repository layer for database access.

## Technical Context

**Language/Version**: Go 1.22+ (backend), PostgreSQL 16+ (database)

**Primary Dependencies**:
- `golang-migrate/migrate` v4 — versioned SQL migration runner (CLI + Go library)
- `lib/pq` — PostgreSQL driver (already in `go.mod` from Feature 001)
- `database/sql` — Go standard library database abstraction
- `google/uuid` — UUID v4/v7 generation

**Storage**: PostgreSQL 16+ (containerized via Docker Compose, already provisioned in Feature 001)

**Testing**: Go standard `testing` package + `DATA-DOG/go-sqlmock` (already in `go.mod`). Migration tests run against the real PostgreSQL container.

**Target Platform**: Linux containers (Docker), multi-arch (amd64/arm64)

**Project Type**: Web service backend (monorepo with `backend/` and `frontend/`)

**Performance Goals**: Flag state lookups by `api_key_hash` < 2ms for 500+ flags per environment. Migration execution < 5 seconds. No deadlocks under 1,000 concurrent reads/second.

**Constraints**: All entity IDs use UUID. JSONB columns for targeting rules and remote config. `last_evaluated_at` updated asynchronously to avoid write contention. Audit logs are append-only (no UPDATE/DELETE).

**Scale/Scope**: Initial schema supports hundreds of projects with thousands of flags and millions of evaluation records.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. API-First Contract Design | ✅ PASS | This feature defines data models; API contracts are Feature 003's scope. Migration contracts and Go repository interfaces are defined here. |
| II. Environment Isolation | ✅ PASS | `environment_flag_states` is keyed by `(environment_id, feature_flag_id)`. Environments have unique `api_key_hash`. Schema enforces isolation at the database level. |
| III. Governance by Default | ✅ PASS | `change_requests`, `change_request_approvals`, `audit_logs`, `roles`, and `user_roles` tables are all core schema — not behind a paywall or optional feature gate. |
| IV. Local Evaluation Performance | ✅ PASS | `idx_env_flag_states_env_id` and `idx_environments_api_key_hash` indexes support sub-millisecond lookups. `last_evaluated_at` updates are asynchronous. |
| V. Test-First Quality Gates | ✅ PASS | Migration up/down tests and repository unit tests are planned. 80% backend coverage maintained. |
| VI. OpenFeature Interoperability | ✅ PASS (N/A) | Data model is internal; SDK/API conformance is handled in Features 003–004. |
| VII. PII Protection & Compliance | ✅ PASS | API keys stored as SHA-256 hashes. Audit logs sanitize PII. `actor_ip` stored for compliance but no plaintext tokens. |
| VIII. Cloud-Native Portability | ✅ PASS | PostgreSQL runs in Docker. Migrations are SQL files executable in any PostgreSQL environment. No cloud-vendor-specific features. |

**Gate Result**: ✅ ALL PASS — proceed to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/002-data-model-state/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (Go repository interfaces)
└── tasks.md             # Phase 2 output (/speckit-tasks command)
```

### Source Code (repository root)

```text
backend/
├── cmd/server/main.go                  # Entry point (updated: add migration runner)
├── internal/
│   ├── config/config.go                # Env config (updated: add migration path)
│   ├── health/handler.go               # Health check (existing)
│   ├── logging/logger.go               # Logger (existing)
│   ├── models/                         # NEW: Go struct definitions
│   │   ├── project.go
│   │   ├── environment.go
│   │   ├── feature_flag.go
│   │   ├── environment_flag_state.go
│   │   ├── change_request.go
│   │   ├── audit_log.go
│   │   ├── role.go
│   │   └── types.go                    # Shared types: enums, JSONB structs
│   └── repository/                     # NEW: Database access layer
│       ├── repository.go               # Common DB interface + transaction helpers
│       ├── project_repo.go
│       ├── environment_repo.go
│       ├── flag_repo.go
│       ├── flag_state_repo.go
│       └── audit_repo.go
├── migrations/                         # NEW: golang-migrate SQL files
│   ├── 000001_create_projects.up.sql
│   ├── 000001_create_projects.down.sql
│   ├── 000002_create_environments.up.sql
│   ├── 000002_create_environments.down.sql
│   ├── 000003_create_feature_flags.up.sql
│   ├── 000003_create_feature_flags.down.sql
│   ├── 000004_create_environment_flag_states.up.sql
│   ├── 000004_create_environment_flag_states.down.sql
│   ├── 000005_create_change_requests.up.sql
│   ├── 000005_create_change_requests.down.sql
│   ├── 000006_create_audit_logs.up.sql
│   ├── 000006_create_audit_logs.down.sql
│   ├── 000007_create_roles.up.sql
│   ├── 000007_create_roles.down.sql
│   ├── 000008_create_user_roles.up.sql
│   └── 000008_create_user_roles.down.sql
└── go.mod                              # Updated: add golang-migrate, google/uuid
```

**Structure Decision**: Follows the existing `backend/internal/` convention from Feature 001. New `models/` and `repository/` packages are added under `internal/` to maintain clean separation. SQL migrations live in `backend/migrations/` following `golang-migrate` file naming conventions (`NNNNNN_description.up.sql` / `.down.sql`).

## Complexity Tracking

No constitution violations — this table is empty.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| — | — | — |
