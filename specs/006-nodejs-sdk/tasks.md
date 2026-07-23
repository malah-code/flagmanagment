# Tasks: Node.js / TypeScript SDK

**Input**: Design documents from `/specs/006-nodejs-sdk/`

**Prerequisites**: plan.md, spec.md, data-model.md, contracts/sdk-interface.md, research.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure for the Node.js SDK

- [x] T001 Create project directory `sdk/node/` and structure per implementation plan
- [x] T002 Initialize `package.json` with TypeScript, `jest`, `@openfeature/server-sdk`, `murmurhash3js`, and `@grpc/grpc-js` dependencies
- [x] T003 [P] Configure `tsconfig.json` for CommonJS and ESM output
- [x] T004 [P] Configure `jest.config.js` for TypeScript testing

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 [P] Create base types in `sdk/node/src/types.ts` (`FlagRule`, `EvaluationContext`, `EvaluationResult`)
- [x] T006 Implement the `RuleStore` class in `sdk/node/src/store.ts` for in-memory flag caching

**Checkpoint**: Foundation ready - user story implementation can now begin.

---

## Phase 3: User Story 1 - Initialize the SDK Client (Priority: P1) 🎯 MVP

**Goal**: Establish a secure connection to the backend and retrieve the initial ruleset snapshot.

**Independent Test**: Write a script instantiating the client and verifying it pulls the initial ruleset via the REST endpoint.

### Tests for User Story 1 ⚠️

- [x] T007 [P] [US1] Integration test for client initialization in `sdk/node/tests/client.test.ts`

### Implementation for User Story 1

- [x] T008 [P] [US1] Implement `fetchSnapshot` logic in `sdk/node/src/sync.ts` using native fetch or HTTP client
- [x] T009 [US1] Create the main `FlagManagmentClient` class in `sdk/node/src/client.ts` with the `init()` method (depends on T008, T006)
- [x] T010 [US1] Add robust error handling for unauthorized tokens (401) and unavailable servers (503)

**Checkpoint**: At this point, the SDK can successfully connect and download the ruleset into memory.

---

## Phase 4: User Story 2 - Local Flag Evaluation (Priority: P1)

**Goal**: Evaluate feature flags locally in memory without network latency (under 1ms).

**Independent Test**: Evaluate multivariate flags using different user identities and verify consistent bucketing.

### Tests for User Story 2 ⚠️

- [x] T011 [P] [US2] Unit test for MurmurHash3 bucketing logic in `sdk/node/tests/evaluator.test.ts`

### Implementation for User Story 2

- [x] T012 [P] [US2] Implement `EvaluateFlag` logic in `sdk/node/src/evaluator.ts` using `murmurhash3js` (must exactly mirror backend logic)
- [x] T013 [US2] Add `evaluateBoolean`, `evaluateString`, etc., methods to `FlagManagmentClient` in `sdk/node/src/client.ts` (calls T012)
- [x] T014 [US2] Ensure evaluations fallback to the `defaultVariation` securely on missing flags

**Checkpoint**: The SDK can now fully evaluate flags offline using the initial snapshot.

---

## Phase 5: User Story 3 - Real-time Delta Updates (Priority: P2)

**Goal**: Automatically receive rule updates over a streaming connection.

**Independent Test**: Use the backend Admin API to change a flag state, and observe the SDK update instantly without restarting.

### Tests for User Story 3 ⚠️

- [x] T015 [P] [US3] Unit/Integration test for streaming reconnects in `sdk/node/tests/sync.test.ts`

### Implementation for User Story 3

- [x] T016 [US3] Implement SSE/gRPC streaming subscription logic in `sdk/node/src/sync.ts`
- [x] T017 [US3] Update `client.ts` to process incoming deltas and apply updates to `RuleStore`
- [x] T018 [US3] Add exponential backoff and reconnection logic for dropped streams

**Checkpoint**: The SDK is now fully synchronized in real-time.

---

## Phase 6: User Story 4 - OpenFeature Provider Compliance (Priority: P2)

**Goal**: Use the FlagManagment SDK as an OpenFeature Provider to avoid vendor lock-in.

**Independent Test**: Evaluate a flag using `@openfeature/server-sdk` directly.

### Tests for User Story 4 ⚠️

- [x] T019 [P] [US4] Unit test for provider compliance in `sdk/node/tests/provider.test.ts`

### Implementation for User Story 4

- [x] T020 [US4] Implement `FlagManagmentProvider` implementing `@openfeature/server-sdk`'s `Provider` interface in `sdk/node/src/provider.ts`
- [x] T021 [US4] Map OpenFeature `EvaluationContext` to FlagManagment `EvaluationContext`

**Checkpoint**: The SDK natively supports OpenFeature!

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T022 Export all public interfaces in `sdk/node/src/index.ts`
- [x] T023 Run `quickstart.md` validation script manually against the Go backend to ensure 100% functionality
- [x] T024 Write the final `README.md` for the SDK package

## Phase 8: Convergence

- [x] T025 Implement true SSE/fetch streaming client with exponential backoff in `sync.ts` per FR-004/SC-002 (partial)
- [x] T026 Configure `tsconfig.json` or build tools to output both CommonJS and ESM per FR-001 (partial)
