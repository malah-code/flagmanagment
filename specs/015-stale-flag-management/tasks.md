# Tasks: Stale Flag Detection & Lifecycle Management

**Input**: Design documents from `/specs/015-stale-flag-management/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Database schema initialization for the feature.

- [X] T001 Create migration for `flag_lifecycle_state` ENUM, `environment_flag_states` columns (`lifecycle_state`, `last_evaluated_at`, `last_state_change_at`), and `stale_flag_policies` table in `backend/migrations/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure and models that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 [P] Add `FlagLifecycleState` enum definition in `backend/internal/models/stale.go`
- [X] T003 [P] Add `StaleFlagPolicy` model struct in `backend/internal/models/stale.go`
- [X] T004 [P] Update `EnvironmentFlagState` struct to include `LifecycleState` and timestamps in `backend/internal/models/flag.go`
- [X] T005 [P] Update `FlagStateRepo` to support filtering by `lifecycle_state` in `backend/internal/repository/flag_state_repo.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Automatic Stale Flag Detection & Tracking (Priority: P1) 🎯 MVP

**Goal**: Automatically track evaluation activity and state stability, detecting flags that have remained unchanged or inactive for 30+ days.

**Independent Test**: Can be tested by simulating flag evaluation events and timestamp aging via DB updates, verifying that the background job correctly transitions flag status to `STALE`.

### Implementation for User Story 1

- [X] T006 [P] [US1] Create repository method `UpdateLastEvaluatedAtBatch` for batch updating evaluation timestamps in `backend/internal/repository/flag_state_repo.go`
- [X] T007 [P] [US1] Create repository method `FindActiveFlagsForStalenessScan` to fetch flags for scanning in `backend/internal/repository/flag_state_repo.go`
- [X] T008 [US1] Implement `MetricAggregationService` for in-memory batching of SDK evaluation timestamps (flushing every ~10s) in `backend/internal/services/metric_aggregation.go`
- [X] T009 [US1] Integrate `MetricAggregationService` into the SDK flag evaluation endpoints to asynchronously record hits in `backend/internal/api/sdk.go`
- [X] T010 [US1] Implement `StaleScannerService` background worker to evaluate active flags against the default 30-day threshold and update their lifecycle state to `STALE` in `backend/internal/services/stale_scanner.go`
- [X] T011 [US1] Register and schedule `StaleScannerService` in `backend/cmd/server/main.go`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. Flags should be automatically marked as stale.

---

## Phase 4: User Story 2 - Stale Flag Dashboard & Lifecycle Actions (Priority: P2)

**Goal**: Provide a dashboard view to review stale flags and take lifecycle actions (Archive, Deprecate, Restore).

**Independent Test**: Can be tested by navigating to the flag list, filtering by `STALE` status, executing an `ARCHIVE` action, and verifying the flag is excluded from the SDK streaming API.

### Implementation for User Story 2

- [X] T012 [P] [US2] Implement API endpoints for `POST /lifecycle/{action}` (`archive`, `deprecate`, `restore`) with RBAC checks and audit logging in `backend/internal/api/lifecycle.go`
- [X] T013 [P] [US2] Update SDK streaming payload generator to exclude flags with `ARCHIVED` lifecycle state in `backend/internal/services/sdk_payload.go`
- [X] T014 [P] [US2] Update React types for `LifecycleState` enum in `frontend/src/types/flag.ts`
- [X] T015 [US2] Update API client to include lifecycle action endpoints in `frontend/src/services/api.ts`
- [X] T016 [US2] Add Stale filter and status badges to the Flag List UI in `frontend/src/components/flags/FlagList.tsx`
- [X] T017 [US2] Implement Lifecycle Action Menu (Archive, Deprecate, Restore) on the Flag Detail/List views in `frontend/src/components/flags/FlagActions.tsx`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently. Stale flags can be viewed and archived via the UI.

---

## Phase 5: User Story 3 - Configurable Staleness Policies (Priority: P3)

**Goal**: Allow Project Owners to configure custom staleness thresholds per project or environment.

**Independent Test**: Can be tested by updating a project's staleness threshold from 30 days to 14 days and verifying that flags unchanged for 15 days are categorized as stale.

### Implementation for User Story 3

- [X] T018 [P] [US3] Create repository methods for CRUD operations on `stale_flag_policies` in `backend/internal/repository/stale_policy_repo.go`
- [X] T019 [US3] Add API endpoints for managing `stale_flag_policies` in `backend/internal/api/stale_policy.go`
- [X] T020 [US3] Update `StaleScannerService` to fetch and use custom environment/project policies instead of the hardcoded default in `backend/internal/services/stale_scanner.go`
- [X] T021 [P] [US3] Add Stale Policy React types and API client methods in `frontend/src/services/api.ts`
- [X] T022 [US3] Add Stale Policy configuration UI to Project Settings and Environment Settings in `frontend/src/pages/ProjectSettings.tsx` and `frontend/src/pages/EnvironmentSettings.tsx`

**Checkpoint**: All user stories should now be independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T023 [P] Add backend unit tests for `MetricAggregationService` and `StaleScannerService` in `backend/internal/services/`
- [X] T024 Run `quickstart.md` validation scenarios manually to verify End-to-End correctness.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2)
- **User Story 2 (P2)**: Can start after Foundational (Phase 2). API endpoints and UI depend on Foundational models.
- **User Story 3 (P3)**: Can start after Foundational (Phase 2). Modifies US1 scanner logic to use policy DB.

### Parallel Opportunities

- All Setup/Foundational tasks marked `[P]` can run in parallel.
- Once Foundational phase completes, User Stories 1, 2, and 3 can be started in parallel since they touch different parts of the system (Background worker, UI/API endpoints, Policy CRUD/UI).
- UI development for US2 and US3 can proceed in parallel with their respective backend API tasks.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently (Simulate flag aging and ensure background worker detects it).
5. Deploy MVP if needed.

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently
3. Add User Story 2 → Test UI/Actions independently
4. Add User Story 3 → Test policy configurations
5. Each story adds value without breaking previous stories
