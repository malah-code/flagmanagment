# Tasks: Telemetry Ingestion and Kill-Switches

**Input**: Design documents from `/specs/009-telemetry-kill-switches/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Setup project structure and feature documentation in `specs/009-telemetry-kill-switches/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

- [x] T002 Create `kill_switches` database migration in `backend/migrations/000010_create_kill_switches.up.sql`
- [x] T003 Create `KillSwitchRule` model in `backend/internal/models/kill_switch.go`
- [x] T004 Create `KillSwitchRepository` in `backend/internal/repository/kill_switch_repo.go`
- [x] T005 Register `KillSwitchRepo` in central store interface `backend/internal/repository/store.go`

---

## Phase 3: User Story 1 - Ingest APM Alerts (Priority: P1) 🎯 MVP

**Goal**: Ingest APM webhook alerts securely and validate environment bearer tokens.

**Independent Test**: Send a mock webhook payload with a valid authorization token and verify response HTTP 202.

- [x] T006 [US1] Create `WebhookService` skeleton in `backend/internal/services/webhook_service.go`
- [x] T007 [US1] Create `WebhookHandler` and ingestion endpoint in `backend/internal/api/webhooks.go`
- [x] T008 [US1] Register `POST /api/v1/webhooks/apm` endpoint with `AuthMiddleware` in `backend/cmd/server/main.go`
- [x] T009 [US1] Add unit tests for webhook authentication and ingestion in `backend/internal/api/webhooks_test.go`

---

## Phase 4: User Story 2 - Trigger Automated Kill Switch (Priority: P1)

**Goal**: Automatically disable linked feature flags when a matching APM alert is received and record an audit log.

**Independent Test**: Configure a flag with a kill switch, send a matching alert webhook, and verify flag becomes disabled with audit log recorded.

- [x] T010 [US2] Implement matching logic in `WebhookService` to query rules by `alert_identifier` in `backend/internal/services/webhook_service.go`
- [x] T011 [US2] Implement flag state update and Redis cache invalidation in `WebhookService`
- [x] T012 [US2] Log automated kill-switch action to `AuditService` in `backend/internal/services/webhook_service.go`
- [x] T013 [US2] Add unit tests for `WebhookService` kill-switch matching logic in `backend/internal/services/webhook_service_test.go`

---

## Phase 5: User Story 3 - View Alert and Kill-Switch History (Priority: P2)

**Goal**: Allow Release Managers to view and configure kill switch rules for feature flags in the UI.

**Independent Test**: Open a feature flag's targeting page in the UI, view, add, or remove kill switch rules.

- [x] T014 [US3] Create `KillSwitchHandler` CRUD API in `backend/internal/api/kill_switches.go` and register routes in `backend/cmd/server/main.go`
- [x] T015 [US3] Create `killSwitchApi.ts` frontend service in `frontend/src/services/killSwitchApi.ts`
- [x] T016 [US3] Create `KillSwitchForm.tsx` component in `frontend/src/components/KillSwitchForm.tsx`
- [x] T017 [US3] Integrate `KillSwitchForm` into `frontend/src/components/flagStates/FlagStatesList.tsx`

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification, cleanup, and documentation

- [x] T018 Code cleanup and review
- [x] T019 Verify backend unit tests and frontend build

## Phase 7: Convergence

- [x] T020 Pass full APM webhook payload to `ProcessAPMAlert` and include it in the Audit Log entry so the exact alert details are visible in history per US3/AC1 (partial)
