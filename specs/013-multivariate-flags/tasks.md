# Tasks: Multivariate Flags & Percentage Rollouts (A/B/n Testing)

**Input**: Design documents from `/specs/013-multivariate-flags/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `backend/` and `frontend/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Add `github.com/spaolacci/murmur3` dependency to the Go backend via `go get github.com/spaolacci/murmur3` in `backend/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T002 Update Go `models` constants in `backend/internal/models/types.go` to explicitly support `MULTIVARIATE` string constants.

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Define Multivariate Flags and Variations (Priority: P1) 🎯 MVP

**Goal**: Allow creation of feature flags with multiple variations instead of just true/false.

**Independent Test**: Can be tested by creating a flag of type MULTIVARIATE, defining custom variation payloads, and verifying they are saved in the database.

### Implementation for User Story 1

- [x] T003 [P] [US1] Add `Variations` (JSONB array of variation objects) to `FeatureFlag` model in `backend/internal/models/feature_flag.go`
- [x] T004 [P] [US1] Update API request/response structs for flag creation and updates to include variations in `backend/internal/api/handlers/flag_handler.go`
- [x] T005 [P] [US1] Add frontend TypeScript interfaces for `Variation` inside `frontend/src/types/index.ts`
- [x] T006 [US1] Update frontend flag creation form to dynamically add/remove variations when `MULTIVARIATE` type is selected in `frontend/src/components/flags/CreateFlagDialog.tsx`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Configure Percentage-Based Rollouts (Priority: P1)

**Goal**: Allocate specific percentages of traffic to each variation (e.g., 10% to Red, 40% to Blue, 50% to Green).

**Independent Test**: Can be tested by setting rollout percentages in a specific environment and verifying the total equals 100%.

### Implementation for User Story 2

- [x] T007 [P] [US2] Update `EnvironmentFlagState` model to include `DefaultVariation` and `RolloutRules` JSONB fields in `backend/internal/models/environment_flag_state.go`
- [x] T008 [US2] Add backend validation to ensure `RolloutRules` percentages sum to exactly 10,000 basis points (100%) during state updates in `backend/internal/models/environment_flag_state.go` and `backend/internal/api/flags.go`
- [x] T009 [P] [US2] Update frontend TypeScript interfaces for `EnvironmentFlagState` in `frontend/src/types/index.ts`
- [x] T010 [US2] Create or update frontend UI for rollout sliders/inputs in `frontend/src/components/flagStates/RolloutRuleBuilder.tsx`, strictly enforcing the 100% total validation on the client.

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Deterministic SDK Bucketing (Priority: P1)

**Goal**: Use identity hashes to bucket users deterministically based on percentages so the same user always sees the same variation.

**Independent Test**: Can be tested by passing the same identity to the SDK 100 times and getting the exact same result in < 1ms.

### Tests for User Story 3

- [x] T011 [P] [US3] Write unit tests in `backend/internal/sdk/evaluator_test.go` to assert deterministic output and roughly even distribution over 10,000 hashes.

### Implementation for User Story 3

- [x] T012 [US3] Implement `murmur3.Sum32([]byte(flagKey + identityKey)) % 10000` logic in `backend/internal/sdk/evaluator.go` to select the bucket and match against `RolloutRules`.
- [x] T013 [US3] Ensure the SDK Evaluation endpoint returns the correct variation ID and payload structure based on the hash evaluation in `backend/internal/api/sdk.go`.

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T014 Run validation commands from `quickstart.md`
- [x] T015 Ensure Edge SDK caching propagates rollout rule changes correctly in `backend/internal/services/change_request_service.go` or relevant publisher logic.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Can start immediately
- **Foundational (Phase 2)**: Depends on Setup
- **User Stories (Phase 3+)**: US1, US2, and US3 are highly sequenced here. US2 depends on US1's variations. US3 depends on US2's rollout rules.
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### Parallel Opportunities

- T003, T004, and T005 can be executed in parallel.
- T007, T009, and T011 can be executed in parallel.
