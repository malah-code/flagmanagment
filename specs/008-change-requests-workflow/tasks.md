# Tasks: Change Requests Workflow

**Input**: Design documents from `/specs/008-change-requests-workflow/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Initialize database migration file for `ChangeRequest` table and `Environment` changes in `backend/internal/db/migrations/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T002 [P] Create `Environment` DB model changes (add `is_protected`) in `backend/internal/models/environment.go`
- [x] T003 [P] Create `ChangeRequest` DB model in `backend/internal/models/change_request.go`
- [x] T004 Define `ChangeRequest` repository interface in `backend/internal/repository/change_request_repo.go`
- [x] T005 Seed `Release Manager` role in the DB (if not exists) via `backend/internal/db/seeds.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Create Change Request in Protected Environment (Priority: P1) 🎯 MVP

**Goal**: As a team member, I want to propose a flag configuration change in a protected environment, so that it can be safely reviewed before affecting production users.

**Independent Test**: Can be tested by toggling a flag in a protected environment and verifying a 202 is returned and a pending Change Request is saved.

### Implementation for User Story 1

- [x] T006 [P] [US1] Implement `Environment` update API to set `is_protected` in `backend/internal/api/environment_handler.go`
- [x] T007 [US1] Implement `ChangeRequestService.Create` in `backend/internal/services/change_request_service.go`
- [x] T008 [US1] Update `FlagService.UpdateFlag` to intercept mutations for protected environments and call `ChangeRequestService` in `backend/internal/services/flag_service.go`
- [x] T009 [P] [US1] Update frontend API client for flag toggling to handle `202 Accepted` in `frontend/src/services/flagApi.ts`
- [x] T010 [P] [US1] Update frontend `EnvironmentSettings.tsx` to allow toggling the `is_protected` flag in `frontend/src/pages/EnvironmentSettings.tsx`
- [x] T011 [US1] Add backend unit tests for flag interception in `backend/internal/services/flag_service_test.go`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Review and Approve Change Request (Priority: P1)

**Goal**: As a Release Manager, I want to review a visual diff of proposed changes and approve them, so that safe configurations are atomically applied to the protected environment.

**Independent Test**: Can be tested by logging in as a Release Manager, viewing the pending request diff, and clicking Approve to verify the changes go live.

### Implementation for User Story 2

- [x] T012 [P] [US2] Implement `ChangeRequestService.List` and `Approve` in `backend/internal/services/change_request_service.go`
- [x] T013 [US2] Implement API endpoints for listing and approving Change Requests in `backend/internal/api/change_request_handler.go`
- [x] T014 [US2] Enforce `Release Manager` role and self-approval restriction in `backend/internal/api/change_request_handler.go`
- [x] T015 [US2] Ensure atomic application of flag state inside `ChangeRequestService.Approve` via DB transaction.
- [x] T016 [P] [US2] Implement frontend `changeRequestApi.ts` for listing and approving in `frontend/src/services/changeRequestApi.ts`
- [x] T017 [US2] Create visual diff component `ChangeRequestDiff.tsx` in `frontend/src/components/ChangeRequestDiff.tsx` (using `react-diff-viewer` or similar)
- [x] T018 [US2] Create frontend page `ChangeRequestsPage.tsx` to list and approve requests in `frontend/src/pages/ChangeRequestsPage.tsx`
- [x] T019 [US2] Add backend unit tests for approval logic and RBAC in `backend/internal/services/change_request_service_test.go`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Reject Change Request (Priority: P2)

**Goal**: As a Release Manager, I want to reject an unsafe or incorrect Change Request, so that bad configurations do not reach production.

**Independent Test**: Can be tested by rejecting a Pending request and verifying it is marked as Rejected without altering the live environment.

### Implementation for User Story 3

- [x] T020 [P] [US3] Implement `ChangeRequestService.Reject` in `backend/internal/services/change_request_service.go`
- [x] T021 [US3] Implement API endpoint for rejecting Change Requests in `backend/internal/api/change_request_handler.go`
- [x] T022 [P] [US3] Add reject functionality to frontend API client in `frontend/src/services/changeRequestApi.ts`
- [x] T023 [US3] Add Reject button and reason modal to `ChangeRequestsPage.tsx` in `frontend/src/pages/ChangeRequestsPage.tsx`
- [x] T024 [US3] Add backend unit tests for rejection logic in `backend/internal/services/change_request_service_test.go`

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T025 [P] Audit Log integration for all Change Request transitions in `backend/internal/services/change_request_service.go`
- [x] T026 Code cleanup and refactoring
- [x] T027 Run quickstart.md validation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Integrates with US1 by acting on the created requests, but API can be built independently.
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Integrates with US1 and US2.

### Within Each User Story

- Models before services
- Services before endpoints
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- Foundational DB models can run in parallel
- Frontend API clients and UI components can run in parallel with backend endpoints
- Different user stories can be worked on in parallel by different team members

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

## Phase 7: Convergence

- [x] T028 Add `CurrentState` JSONB to `ChangeRequest` model and repository, and populate it when creating a change request per spec.md: Key Entities (missing)
- [x] T029 Implement true visual diffing (e.g., using `react-diff-viewer`) in `ChangeRequestDiff.tsx` consuming `CurrentState` and `ProposedChanges` per FR-004 (partial)

---

## Phase 8: Convergence

**Purpose**: Addressing unbuilt work identified during the comprehensive codebase audit.

- [x] T030: Wire the existing `ChangeRequestsPage.tsx` page into the application routing `App.tsx` and sidebar navigation. links per US2 (missing)
