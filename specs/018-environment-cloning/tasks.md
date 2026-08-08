# Tasks: API-Driven Environment Cloning (Ephemeral Environments)

**Input**: Design documents from `/specs/018-environment-cloning/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

*(No setup tasks needed. Backend architecture is established)*

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

*(No foundational DB schema changes needed. The environment and environment_flag_states tables already exist)*

---

## Phase 3: User Story 1 - CI/CD Ephemeral Environment Creation (Priority: P1) 🎯 MVP

**Goal**: Programmatically clone an existing environment into a new, temporary environment so automated integration tests can run against isolated flag states.

**Independent Test**: Hit the clone API endpoint and verify a new SDK token and environment are created with copied flag states.

### Implementation for User Story 1

- [ ] T001 [US1] Implement `CloneEnvironmentState` in `backend/internal/repository/environment_repo.go` using a single `pgx` transaction.
- [ ] T002 [US1] Implement `CloneEnvironment` logic in `backend/internal/services/environment_service.go` (generate token, call repo clone, log audit event).
- [ ] T003 [US1] Expose `POST /api/v1/projects/{projectId}/environments/{sourceEnvId}/clone` endpoint in `backend/internal/api/environment.go`.
- [ ] T004 [P] [US1] Add unit tests for `CloneEnvironment` in `backend/internal/services/environment_service_test.go`.

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently.

---

## Phase 4: User Story 2 - Ephemeral Environment Teardown (Priority: P1)

**Goal**: Programmatically delete ephemeral environments once CI/CD pipeline completes to prevent clutter.

**Independent Test**: Hit the delete API endpoint and verify the environment and its states are permanently removed (and protected environments reject deletion).

### Implementation for User Story 2

- [ ] T005 [US2] Implement `DeleteEnvironment` in `backend/internal/repository/environment_repo.go`.
- [ ] T006 [US2] Implement `DeleteEnvironment` business logic in `backend/internal/services/environment_service.go` (prevent deletion if `IsProtected`).
- [ ] T007 [US2] Expose `DELETE /api/v1/projects/{projectId}/environments/{envId}` endpoint in `backend/internal/api/environment.go`.
- [ ] T008 [P] [US2] Add unit tests for `DeleteEnvironment` in `backend/internal/services/environment_service_test.go`.

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T009 [P] Run `make test-backend` to ensure all tests pass.
- [ ] T010 [P] Execute curl scenarios from `quickstart.md` against local running server to manually validate.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: N/A
- **Foundational (Phase 2)**: N/A
- **User Stories (Phase 3+)**: 
- **Polish (Final Phase)**: Depends on all user stories being complete.

### User Story Dependencies

- **User Story 1 (P1)**: No dependencies.
- **User Story 2 (P2)**: Depends on User Story 1 for end-to-end testing (clone an env, then delete it).

### Parallel Opportunities

- Unit tests for US1 and US2 can be run/written in parallel once the business logic is in place.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 3: User Story 1
2. **STOP and VALIDATE**: Test User Story 1 independently

### Incremental Delivery

1. Add User Story 1 → Test independently
2. Add User Story 2 → Test independently (create and delete)
