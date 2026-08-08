# Tasks: Sequential Dependencies

**Input**: Design documents from `/specs/016-sequential-dependencies/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Database schema initialization for the feature.

- [ ] T001 Create migration for `parent_flag_id` UUID column in `feature_flags` table in `backend/migrations/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure and models that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T002 [P] Update `FeatureFlag` struct to include `ParentFlagID *uuid.UUID` in `backend/internal/models/feature_flag.go`
- [ ] T003 [P] Add unit test cases for the cycle detection algorithm in `backend/internal/services/cycle_detector_test.go`
- [ ] T004 Implement `CycleDetectorService` to perform DFS cycle detection in `backend/internal/services/cycle_detector.go`
- [ ] T005 Update `FlagRepo` to fetch flag dependency chains in `backend/internal/repository/flag_repo.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Configure a Feature Flag Dependency (Priority: P1) 🎯 MVP

**Goal**: Allow users to configure a feature flag to depend on the state of another (parent) flag safely without creating infinite loops.

**Independent Test**: Can be tested by creating Flag A and Flag B, and configuring Flag B to depend on Flag A via the API.

### Implementation for User Story 1

- [ ] T006 [P] [US1] Update Flag POST/PUT handlers to accept `parent_flag_id` in `backend/internal/api/flags.go`
- [ ] T007 [US1] Update `FlagService.Create` and `FlagService.Update` to invoke `CycleDetectorService` before saving to DB in `backend/internal/services/flag_service.go`
- [ ] T008 [P] [US1] Update frontend `Flag` type definitions with `parent_flag_id` in `frontend/src/types/index.ts`
- [ ] T009 [US1] Add a "Depends On" optional dropdown in `frontend/src/components/flags/CreateFlagDialog.tsx` selecting from other flags

**Checkpoint**: At this point, User Story 1 should be fully functional. Dependencies can be saved without cycles.

---

## Phase 4: User Story 2 - Evaluate Dependent Flags (Priority: P1)

**Goal**: Safely evaluate dependent flags at runtime, returning the fallback state if the parent flag is not ON.

**Independent Test**: Can be tested by evaluating Flag B via SDK when Flag A is forced to OFF.

### Implementation for User Story 2

- [ ] T010 [P] [US2] Update SDK evaluation unit tests to assert short-circuiting logic in `backend/internal/sdk/evaluator_test.go`
- [ ] T011 [US2] Update `LocalEvaluator` to recursively check `parent_flag_id` state before evaluating targeting rules in `backend/internal/sdk/evaluator.go`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently. Dependencies are evaluated correctly at runtime.

---

## Phase 5: User Story 3 - UI Visibility of Dependencies (Priority: P2)

**Goal**: Clearly show which flags depend on a parent flag in the dashboard.

**Independent Test**: Can be tested by navigating to a flag's detail page and viewing its parent dependency.

### Implementation for User Story 3

- [ ] T012 [P] [US3] Add a visual indicator in the Flag List identifying flags with parents in `frontend/src/components/flags/FlagsList.tsx`
- [ ] T013 [US3] Display the `parent_flag_id` reference explicitly in the flag configuration page in `frontend/src/components/flags/FlagDetails.tsx`

**Checkpoint**: All user stories should now be independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T014 Run `quickstart.md` validation scenarios manually to verify End-to-End correctness.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2)
- **User Story 2 (P1)**: Can start after Foundational (Phase 2).
- **User Story 3 (P2)**: Can start after Foundational (Phase 2).

### Parallel Opportunities

- All Setup/Foundational tasks marked `[P]` can run in parallel.
- Once Foundational phase completes, User Stories 1, 2, and 3 can be started in parallel since they touch different parts of the system (API/UI, SDK, UI).

---

## Implementation Strategy

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently
3. Add User Story 2 → Test evaluation independently
4. Add User Story 3 → Test UI independently
5. Each story adds value without breaking previous stories
