---

description: "Task list template for feature implementation"
---

# Tasks: SDK Evaluation API

**Input**: Design documents from `/specs/005-sdk-evaluation-api/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 [P] Ensure Redis is configured in docker-compose.yml and local environment
- [x] T002 Copy `sdk.proto` and `evaluate-api.yaml` to backend contract directories
- [x] T003 Generate Go code from `sdk.proto` using protoc into `backend/pkg/gen/sdk/v1/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Setup gRPC server initialization alongside HTTP in `backend/cmd/server/main.go`
- [x] T005 [P] Create Redis client wrapper and PubSub helper in `backend/internal/cache/redis.go`
- [x] T006 [P] Create foundational data models (RulesetSnapshot, FlagRule, EvaluationContext) in `backend/internal/models/sdk.go`
- [x] T007 Implement SDK Token validation middleware/interceptor in `backend/internal/middleware/sdk_auth.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - SDK Bootstrapping (Priority: P1) 🎯 MVP

**Goal**: As a server-side SDK, I want to download a complete snapshot of all feature flags and their targeting rules for my environment upon initialization so that I can evaluate flags locally in-memory with zero network latency.

**Independent Test**: Can be fully tested by establishing a gRPC connection with a valid environment SDK token and receiving a JSON/protobuf payload containing all active flags.

### Implementation for User Story 1

- [x] T008 [P] [US1] Implement Redis fetching logic for full snapshot in `backend/internal/cache/ruleset.go`
- [x] T009 [US1] Create gRPC service implementation structure in `backend/internal/sdk/service.go`
- [x] T010 [US1] Implement `FetchSnapshot` handler logic in `backend/internal/sdk/service.go`
- [x] T011 [US1] Integrate gRPC service with the main server router in `backend/cmd/server/main.go`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently via `grpcurl`.

---

## Phase 4: User Story 2 - Real-time Delta Updates (Priority: P2)

**Goal**: As a server-side SDK, I want to establish a persistent streaming connection to receive lightweight delta updates when flags change, so that my local in-memory cache stays accurate without polling.

**Independent Test**: Can be tested by establishing a gRPC stream and verifying that a flag change triggers an immediate push event over the stream.

### Implementation for User Story 2

- [x] T012 [P] [US2] Implement Redis Pub/Sub listener for environment updates in `backend/internal/cache/pubsub.go`
- [x] T013 [US2] Create streaming connection manager in `backend/internal/sdk/stream.go`
- [x] T014 [US2] Implement `StreamRulesets` gRPC handler in `backend/internal/sdk/service.go` linking to the stream manager
- [x] T015 [US2] Wire up the admin API flag update handlers to publish events to Redis Pub/Sub in `backend/internal/handlers/flags.go`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently.

---

## Phase 5: User Story 3 - Server-side Evaluation (Thin Clients) (Priority: P3)

**Goal**: As a thin client (e.g., mobile app or single-page application), I want to call an endpoint to evaluate specific flags for my user context so that I don't have to download the entire ruleset and expose all targeting logic to the public internet.

**Independent Test**: Can be tested by sending an evaluation request with a user context and verifying the correct flag value is returned via the REST API.

### Implementation for User Story 3

- [x] T016 [P] [US3] Implement fast in-memory application cache (e.g., using ristretto) to shadow Redis in `backend/internal/cache/memory.go`
- [x] T017 [US3] Sync in-memory cache with Redis Pub/Sub events in `backend/internal/cache/memory.go`
- [x] T018 [US3] Implement local evaluation engine with MurmurHash3 bucketing in `backend/internal/sdk/evaluator.go`
- [x] T019 [US3] Implement REST POST `/api/v1/sdk/evaluate` handler in `backend/internal/handlers/sdk_evaluate.go`
- [x] T020 [US3] Ensure PII hashing inside context before logging in `backend/internal/handlers/sdk_evaluate.go`

**Checkpoint**: All user stories should now be independently functional.

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T021 [P] Write integration tests using `testcontainers-go` for Redis in `backend/internal/cache/redis_test.go`
- [x] T022 Code cleanup, error wrapping, and structured logging audits across the new SDK package.
- [x] T023 Run quickstart.md validation manually.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2)
- **User Story 2 (P2)**: Can start after Foundational (Phase 2)
- **User Story 3 (P3)**: Can start after Foundational (Phase 2)

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Models and Redis clients marked [P] can run in parallel
