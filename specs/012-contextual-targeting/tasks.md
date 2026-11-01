# Implementation Tasks: Contextual Targeting Engine

## Phase 1: Foundational

- [x] T001 Define `TargetingRule`, `TargetingCondition`, and `Operator` structs in `backend/internal/models/environment_flag_state.go`
- [x] T002 Update `Validate()` in `backend/internal/models/environment_flag_state.go` to validate Regex patterns

## Phase 2: User Story 1 - SDK Evaluation of Contextual Rules (Backend)

- [x] T003 [US1] Create `backend/internal/sdk/evaluator.go` and define `EvaluationContext`
- [x] T004 [US1] Implement `EQUALS`, `CONTAINS`, and `REGEX` matching logic in `evaluator.go`
- [x] T005 [P] [US1] Create unit tests for evaluator in `backend/internal/sdk/evaluator_test.go`
- [x] T006 [US1] Integrate `evaluator.go` into `SDKService` (or wherever flags are evaluated) in `backend/internal/services/flag_state_service.go`

## Phase 3: User Story 2 - Rule Creation for Flags (Frontend)

- [x] T007 [US2] Create `frontend/src/components/flagStates/TargetingRuleBuilder.tsx`
- [x] T008 [US2] Integrate `TargetingRuleBuilder` into the existing flag state editor modal in `frontend/src/components/flagStates/FlagStatesList.tsx` (or related editor component)

## Dependencies

- Phase 1 must be completed before Phase 2.
- Phase 3 (Frontend) can be executed in parallel with Phase 2 (Backend) since the API contract is already defined in `data-model.md`.

## Implementation Strategy

- Implement the Go evaluation engine first to ensure all rules can be parsed and executed within the performance budget.
- Connect the frontend rule builder so users can dynamically generate these JSON rules.
