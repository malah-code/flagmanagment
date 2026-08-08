# Tasks: PII Hashing & Data Privacy

**Input**: Design documents from `/specs/021-pii-hashing-privacy/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

*(No setup tasks required. Structure exists.)*

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T001 Create UP database migration for `salt` column in `backend/migrations/000008_add_environment_salt.up.sql`
- [x] T002 [P] Create DOWN database migration in `backend/migrations/000008_add_environment_salt.down.sql`
- [x] T003 Update `Environment` model in `backend/internal/models/environment.go` to include `Salt` and generate a 32-byte hex string in `BeforeCreate`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Secure Identity Storage (Priority: P1) 🎯 MVP

**Goal**: Hash PII with SHA-256 and an environment-level salt before storing analytics.

**Independent Test**: Evaluate a flag using an email, verify the stored analytic record contains a hashed email, not plaintext.

### Implementation for User Story 1

- [x] T004 [US1] Update `HashPII` function in `backend/internal/sdk/evaluator.go` to accept a `salt string` parameter and use it for SHA-256 hashing.
- [x] T005 [P] [US1] Update unit tests in `backend/internal/sdk/evaluator_test.go` to verify `HashPII` salting and SHA-256 logic.
- [x] T006 [US1] Update API handlers in `backend/internal/api/` (e.g. `backend/internal/api/evaluate.go`) to pass the environment salt to `HashPII` during flag evaluation analytics tracking.

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Consistent Identity Bucketing (Priority: P1)

**Goal**: Ensure the SDK uses MurmurHash3 with the environment salt for deterministic and secure bucketing.

**Independent Test**: Evaluate rollout rules for the same identity across multiple evaluations in the same environment and confirm they hit the same bucket.

### Implementation for User Story 2

- [x] T007 [US2] Update `EvaluateRolloutSplit` in `backend/internal/sdk/evaluator.go` to accept a `salt string` parameter and incorporate it into the MurmurHash3 calculation.
- [x] T008 [P] [US2] Update unit tests in `backend/internal/sdk/evaluator_test.go` to verify `EvaluateRolloutSplit` salting logic.
- [x] T009 [US2] Update `EvaluateFlag` in `backend/internal/sdk/evaluator.go` to pass `salt` down to `EvaluateRolloutSplit`.

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Configurable Data Retention (Priority: P2)

**Goal**: Implement a scheduled job to automatically clean up evaluation analytics older than 30 days.

**Independent Test**: Run the cleanup job manually and verify old records are deleted from the database.

### Implementation for User Story 3

- [x] T010 [US3] Add a `CleanupOldAnalytics(ctx context.Context, days int)` method to the Audit/Analytics repository and service in `backend/internal/services/audit.go` (or analytics equivalent).
- [x] T011 [US3] Implement a background ticker in `backend/cmd/server/main.go` that invokes the cleanup service method periodically (e.g., every 24 hours).

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T012 Run quickstart.md validation locally.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: N/A
- **Foundational (Phase 2)**: BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - US1 and US2 can run in parallel since they modify different functions in `evaluator.go` (though they may conflict on file saving, they don't block each other logically).
  - US3 is independent and can run in parallel.
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - No dependencies on other stories

### Parallel Opportunities

- T002 can be written while T001 is being written.
- Unit test updates (T005, T008) can be run concurrently with the core logic updates (T004, T007).
- Different user stories can be worked on in parallel by different team members.

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
2. Complete Phase 3: User Story 1
3. **STOP and VALIDATE**: Test User Story 1 independently

### Incremental Delivery

1. Complete Foundational → Foundation ready
2. Add User Story 1 → Test independently
3. Add User Story 2 → Test independently
4. Add User Story 3 → Test independently
