# Tasks: SDK Event Forwarding for Analytics

**Input**: Design documents from `/specs/017-sdk-event-forwarding/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create `backend/internal/sdk/hooks` directory

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

*(No foundational DB or infrastructure changes required as this is purely an SDK feature inside the existing Go backend codebase)*

---

## Phase 3: User Story 1 - Standardized SDK Evaluation Hooks (Priority: P1) 🎯 MVP

**Goal**: Allow developers to attach evaluation hooks to capture flag variation events asynchronously without performance penalty.

**Independent Test**: Register a mock hook, evaluate a flag, and assert the hook is invoked with correct context and details asynchronously.

### Implementation for User Story 1

- [ ] T002 [US1] Implement OpenFeature Hook invocation in `backend/internal/sdk/evaluator.go` inside `EvaluateSingleFlag` (execute `After` hooks asynchronously via goroutines).
- [ ] T003 [US1] Ensure `Error` hooks are invoked in `evaluator.go` if flag evaluation fails.
- [ ] T004 [P] [US1] Create unit tests in `backend/internal/sdk/evaluator_test.go` to verify hooks are invoked asynchronously and panics are safely recovered.

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently.

---

## Phase 4: User Story 2 - Integration with External Analytics Providers (Priority: P2)

**Goal**: Provide an integration pattern (reference hook) to seamlessly forward evaluation events to external analytics tools.

**Independent Test**: Use the reference analytics hook and verify that it formats and "forwards" the event correctly (e.g., via logging for the quickstart).

### Implementation for User Story 2

- [ ] T005 [P] [US2] Implement reference `AnalyticsHook` (conforming to OpenFeature `Hook` interface) in `backend/internal/sdk/hooks/analytics_hook.go`.
- [ ] T006 [P] [US2] Create unit tests in `backend/internal/sdk/hooks/analytics_hook_test.go` to verify the `After` method executes properly without blocking.

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T007 [P] Run `make test-backend` to ensure all tests pass.
- [ ] T008 [P] Verify `quickstart.md` validation scenario manually (if necessary).

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: N/A
- **User Stories (Phase 3+)**: Depend on Setup completion.
- **Polish (Final Phase)**: Depends on all user stories being complete.

### User Story Dependencies

- **User Story 1 (P1)**: No dependencies.
- **User Story 2 (P2)**: Depends on User Story 1 (requires hook support in evaluator).

### Parallel Opportunities

- Unit tests in `evaluator_test.go` (T004) can run in parallel with the hook implementation.
- `AnalyticsHook` implementation (T005, T006) can run once `evaluator.go` logic is established.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 3: User Story 1
3. **STOP and VALIDATE**: Test User Story 1 independently

### Incremental Delivery

1. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
2. Add User Story 2 → Test independently → Deploy/Demo
