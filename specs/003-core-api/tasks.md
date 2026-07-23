# Implementation Tasks: Core API Service

**Feature**: `003-core-api`
**Status**: Draft

## Dependencies

- No feature dependencies. Can be built independently.

## Parallel Execution

- API handlers for different entities (Projects, Environments, Flags, SDK) can be built in parallel once Phase 1 (Setup) is complete.

## Implementation Strategy

We will build the DTOs and API middleware first, then implement the route handlers per user story, and finally mount everything in the main router.

---

## Phase 1: Setup & Foundational APIs

- [ ] T001 Create incoming request DTOs with validation tags in `backend/internal/dto/requests.go`
- [ ] T002 Create outgoing response DTOs in `backend/internal/dto/responses.go`
- [ ] T003 Create standardized JSON error response utility in `backend/internal/api/errors.go`
- [ ] T004 Implement Auth and Pagination middleware in `backend/internal/api/middleware.go`

## Phase 2: Workspace & Hierarchy Management (US1)

**Goal**: Manage projects and environments via REST API.

- [ ] T005 [P] [US1] Implement Project Handlers (Create, List, Get, Update, Delete) in `backend/internal/api/projects.go`
- [ ] T006 [P] [US1] Implement Environment Handlers (Create with secure API Key generation, List, Get) in `backend/internal/api/environments.go`

## Phase 3: Feature Flag Configuration (US2)

**Goal**: Create, update, and toggle feature flags via REST API.

- [ ] T007 [P] [US2] Implement Feature Flag Handlers (Create globally per project) in `backend/internal/api/flags.go`
- [ ] T008 [P] [US2] Implement Flag State Handlers (Update targeting rules and state per environment) in `backend/internal/api/flags.go`

## Phase 4: SDK Evaluation Retrieval (US3)

**Goal**: High-performance SDK evaluation endpoint with ETag support.

- [ ] T009 [P] [US3] Implement SDK Evaluation endpoint caching logic and route handler in `backend/internal/api/sdk.go`

## Phase 5: Integration & Polish

- [ ] T010 Mount all routers in `backend/cmd/server/main.go`
- [ ] T011 Write API integration tests in `backend/tests/api/api_test.go`

## Phase 6: Convergence

- [x] T012 Remount the SDK evaluation endpoint to match the contract path `/api/v1/evaluate/flags` per FR-005 (partial)
- [x] T013 Update `projects.go` and `flags.go` list endpoints to use `GetPagination` limits and offsets rather than hardcoded values per FR-008 (partial)
- [x] T014 Expand `api_test.go` to contain integration tests for Projects, Environments, Flags, and SDK endpoints per T011 (partial)
