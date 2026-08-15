# Tasks: Immutable Audit Logs & SIEM Webhooks

**Input**: Design documents from `/specs/020-audit-logs-siem/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: The examples below include test tasks. Tests are OPTIONAL - only include them if explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create package directories `backend/internal/audit` and `backend/internal/webhook`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T002 Generate and write database migrations for `audit_logs` and `webhook_integrations` tables in `backend/migrations/`
- [x] T003 [P] Create `AuditLog` entity in `backend/internal/audit/models.go`
- [x] T004 [P] Create `WebhookIntegration` entity in `backend/internal/webhook/models.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Compliance Auditing and Traceability (Priority: P1) 🎯 MVP

**Goal**: Record an immutable, append-only ledger of all administrative actions in the database.

**Independent Test**: Perform a series of CRUD operations on a feature flag via the API. Query the audit logs endpoint and assert that a distinct log entry exists for each action.

### Implementation for User Story 1

- [x] T005 [US1] Implement `AuditRepository` in `backend/internal/audit/repository.go` for inserting and querying audit logs
- [x] T006 [US1] Implement `AuditService` with transactional hook logic in `backend/internal/audit/service.go` (depends on T005)
- [x] T007 [P] [US1] Implement `GET /api/v1/audit-logs` endpoint in `backend/internal/audit/handler.go`
- [x] T008 [P] [US1] Implement CSV streaming export `GET /api/v1/audit-logs/export` endpoint in `backend/internal/audit/handler.go`
- [x] T009 [US1] Integrate `AuditService` into existing feature flag handlers (e.g., `backend/internal/api/handlers.go`) to record logs on mutations

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Security Log Sanitization (Priority: P1)

**Goal**: Actively sanitize log payloads so that plaintext API keys and PII are never stored.

**Independent Test**: Generate a new environment API key. Query the audit log and assert the key is `[REDACTED]`.

### Implementation for User Story 2

- [x] T010 [P] [US2] Implement generic JSON `Scrub()` function in `backend/internal/audit/sanitizer.go`
- [x] T011 [US2] Integrate `Scrub()` into the `AuditService` (in `backend/internal/audit/service.go`) to sanitize `previous_state` and `new_state` before database insertion

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - SIEM Webhook Streaming (Priority: P2)

**Goal**: Stream audit log events in real-time via webhooks to external SIEM tools.

**Independent Test**: Register a mock webhook endpoint in the project settings. Perform a flag change and assert the mock endpoint receives the POST request.

### Implementation for User Story 3

- [x] T012 [P] [US3] Implement `WebhookRepository` in `backend/internal/webhook/repository.go` for managing configurations
- [x] T013 [P] [US3] Implement POST `/api/v1/projects/:id/webhooks` endpoint in `backend/internal/webhook/handler.go`
- [x] T014 [US3] Implement in-memory `Dispatcher` with `time.AfterFunc` retries in `backend/internal/webhook/dispatcher.go`
- [x] T015 [US3] Wire the `Dispatcher` to listen to an `AuditService` channel for newly created logs and fire HTTP POST requests

**Checkpoint**: All user stories should now be independently functional

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T016 [P] Write unit tests for the `Scrub()` function in `backend/internal/audit/sanitizer_test.go`
- [x] T017 [P] Write unit tests for `Dispatcher` retry logic in `backend/internal/webhook/dispatcher_test.go`
- [x] T018 Register the new HTTP routes (`audit-logs`, `webhooks`) in `backend/cmd/server/main.go`
- [x] T019 Run quickstart.md validation end-to-end

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed sequentially in priority order (US1 -> US2 -> US3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2)
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Augments US1
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Augments US1

### Parallel Opportunities

- All Foundational tasks marked [P] can run in parallel (T003, T004)
- Endpoint implementations (T007, T008) can be built in parallel to core logic
- Unit testing in Phase N can run in parallel

---

## Implementation Strategy

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 (Audit Log CRUD) → Test independently → MVP!
3. Add User Story 2 (Sanitization) → Test independently
4. Add User Story 3 (Webhooks) → Test independently
5. Each story adds value without breaking previous stories

---

## Phase 7: Convergence

**Purpose**: Addressing unbuilt work identified during the comprehensive codebase audit.

- [x] T020: Ensure `AuditLogs.tsx` is wired into the application routing and reachable via the sidebar.
