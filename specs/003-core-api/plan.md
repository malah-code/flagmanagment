# Implementation Plan: Core API Service

**Branch**: `003-core-api` | **Date**: 2026-07-20 | **Spec**: [spec.md](file:///home/tarikelmallah/Projects/FlagManagment/specs/003-core-api/spec.md)

**Input**: Feature specification from `/specs/003-core-api/spec.md`

## Summary

Implement the REST API service layer using `go-chi/chi` to provide CRUD operations for Projects, Environments, and Feature Flags, and a high-performance SDK evaluation endpoint with ETag support.

## Technical Context

**Language/Version**: Go 1.22+
**Primary Dependencies**: `go-chi/chi`, `github.com/go-playground/validator/v10`
**Storage**: PostgreSQL (via existing repository layer)
**Testing**: `net/http/httptest`
**Target Platform**: Linux server, Docker
**Project Type**: REST API Backend Service
**Performance Goals**: <200ms p95 for management APIs, <50ms for SDK evaluation
**Constraints**: Follow Google API pagination (`pageToken`), support ETag/If-None-Match caching for SDKs
**Scale/Scope**: Production-grade MVP REST API

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. API-First Contract Design**: PASS. OpenAPI contract written.
- **II. Environment Isolation**: PASS. Endpoints are strictly scoped and SDK API uses isolated environment keys.
- **IV. Local Evaluation Performance**: PASS. Evaluator uses ETags to optimize SDK sync bandwidth.
- **V. Test-First Quality Gates**: PASS. Unit and integration API tests are planned.
- **VII. PII Protection & Compliance**: PASS. API keys are hashed and stored securely.

## Project Structure

### Documentation (this feature)

```text
specs/003-core-api/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── openapi.yaml
└── tasks.md
```

### Source Code (repository root)

```text
backend/
├── cmd/server/
│   └── main.go           # Add routers
├── internal/
│   ├── api/              # HTTP Handlers and routing
│   │   ├── middleware/   # Auth, logger, recoverer
│   │   ├── projects.go
│   │   ├── environments.go
│   │   ├── flags.go
│   │   └── sdk.go
│   ├── dto/              # API request/response structures
│   │   ├── requests.go
│   │   └── responses.go
│   └── repository/       # Existing layer
└── tests/
    └── api/              # API e2e integration tests
```

**Structure Decision**: Add `internal/api/` for the web tier and `internal/dto/` for request/response bodies, hooking them up in `cmd/server/main.go`.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
