---
description: "Task list for implementing Data Model & State Management"
---

# Tasks: Data Model & State Management

**Input**: Design documents from `specs/002-data-model-state/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and base type definitions.

- [x] T001 Create `backend/internal/models/` and `backend/internal/repository/` directories
- [x] T002 Create `backend/migrations/` directory for SQL migration files
- [x] T003 Add `golang-migrate/migrate/v4` and `github.com/google/uuid` dependencies to `backend/go.mod`
- [x] T004 [P] Create base shared types, enums, and JSONB wrapper structures in `backend/internal/models/types.go`
- [x] T005 [P] Define `Store` interface and shared Error variables (`ErrNotFound`, `ErrAlreadyExists`, etc.) in `backend/internal/repository/repository.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented.

- [x] T006 Add `migrate-up` and `migrate-down` Makefile targets to root or backend Makefile to run `golang-migrate` CLI.

**Checkpoint**: Foundation ready - user story implementation can now begin.

---

## Phase 3: User Story 1 - Database Schema Provisioning & Migration Pipeline (P1) 🎯 MVP

**Goal**: Define version-controlled migrations for the entire database schema to allow automated provisioning and teardown across environments.

**Independent Test**: Scenario 1 and 2 from `quickstart.md` (Running migrate up and migrate down, then verifying tables).

### Implementation for User Story 1

- [x] T007 [P] [US1] Create migration `000001_create_projects.up.sql` and `down.sql` in `backend/migrations/`
- [x] T008 [P] [US1] Create migration `000002_create_environments.up.sql` and `down.sql` in `backend/migrations/`
- [x] T009 [P] [US1] Create migration `000003_create_feature_flags.up.sql` and `down.sql` in `backend/migrations/`
- [x] T010 [P] [US1] Create migration `000004_create_environment_flag_states.up.sql` and `down.sql` in `backend/migrations/`
- [x] T011 [P] [US1] Create migration `000005_create_change_requests.up.sql` and `down.sql` in `backend/migrations/`
- [x] T012 [P] [US1] Create migration `000006_create_audit_logs.up.sql` and `down.sql` in `backend/migrations/`
- [x] T013 [P] [US1] Create migration `000007_create_roles.up.sql` and `down.sql` in `backend/migrations/`
- [x] T014 [P] [US1] Create migration `000008_create_user_roles.up.sql` and `down.sql` in `backend/migrations/`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently via database migration CLI commands.

---

## Phase 4: User Story 2 - Feature Flag State & Targeting Rule Persistence (P2)

**Goal**: Implement Go data models and repository layers for Projects, Environments, Feature Flags, and Environment Flag States to persist configuration.

**Independent Test**: Scenario 3, 4, 5 from `quickstart.md` (testing lookups, constraints, and Go repository tests).

### Implementation for User Story 2

- [x] T015 [P] [US2] Create Project struct model in `backend/internal/models/project.go`
- [x] T016 [P] [US2] Create Environment struct model in `backend/internal/models/environment.go`
- [x] T017 [P] [US2] Create FeatureFlag struct model in `backend/internal/models/feature_flag.go`
- [x] T018 [P] [US2] Create EnvironmentFlagState struct model in `backend/internal/models/environment_flag_state.go`
- [x] T019 [US2] Implement ProjectRepository (`Create`, `GetByID`, `GetByKey`, `List`, `Update`, `Delete`) in `backend/internal/repository/project_repo.go`
- [x] T020 [US2] Implement EnvironmentRepository (`Create`, `GetByID`, `GetByAPIKeyHash`, `ListByProject`, `Update`, `Delete`) in `backend/internal/repository/environment_repo.go`
- [x] T021 [US2] Implement FlagRepository (`Create`, `GetByID`, `GetByKey`, `ListByProject`, `Update`, `Delete`, `UpdateLastEvaluatedAt`) in `backend/internal/repository/flag_repo.go`
- [x] T022 [US2] Implement FlagStateRepository (`Create`, `GetByID`, `GetByEnvAndFlag`, `ListByEnvironment`, `Update`, `Delete`) in `backend/internal/repository/flag_state_repo.go`
- [ ] T023 [P] [US2] Write unit tests for ProjectRepository in `backend/internal/repository/project_repo_test.go`
- [ ] T024 [P] [US2] Write unit tests for EnvironmentRepository in `backend/internal/repository/environment_repo_test.go`
- [ ] T025 [P] [US2] Write unit tests for FlagRepository in `backend/internal/repository/flag_repo_test.go`
- [ ] T026 [P] [US2] Write unit tests for FlagStateRepository in `backend/internal/repository/flag_state_repo_test.go`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently. Projects, Environments, and Flags can be created and queried using Go repositories.

---

## Phase 5: User Story 3 - Governance & Audit Trail Data Storage (P3)

**Goal**: Implement Go data models and repository layers for Change Requests, Approvals, and Audit Logs.

**Independent Test**: Scenario 6 from `quickstart.md` (Audit Log Immutability verification and CRUD tests for change requests).

### Implementation for User Story 3

- [x] T027 [P] [US3] Create ChangeRequest and ChangeRequestApproval struct models in `backend/internal/models/change_request.go`
- [x] T028 [P] [US3] Create AuditLog struct model in `backend/internal/models/audit_log.go`
- [x] T029 [US3] Implement ChangeRequestRepository (`Create`, `GetByID`, `ListByEnvironment`, `UpdateStatus`, `AddApproval`, `ListApprovals`) in `backend/internal/repository/change_request_repo.go`
- [x] T030 [US3] Implement AuditRepository (`Create`, `ListByProject`, `ListByEnvironment`) in `backend/internal/repository/audit_repo.go`
- [x] T031 [P] [US3] Write unit tests for ChangeRequestRepository in `backend/internal/repository/change_request_repo_test.go`
- [x] T032 [P] [US3] Write unit tests for AuditRepository in `backend/internal/repository/audit_repo_test.go`

**Checkpoint**: Governance and Audit models and persistence layer are complete.

---

## Phase 6: User Story 4 - RBAC & Access Control Model Persistence (P4)

**Goal**: Implement Go data models and repository layers for Roles and User Roles to support granular access control.

**Independent Test**: Create roles and assign users, verifying constraints via unit tests.

### Implementation for User Story 4

- [x] T033 [P] [US4] Create Role and UserRole struct models in `backend/internal/models/role.go`
- [x] T034 [US4] Implement RoleRepository (`Create`, `GetByID`, `GetByName`, `List`, `AssignUserRole`, `RemoveUserRole`, `GetUserRoles`) in `backend/internal/repository/role_repo.go`
- [x] T035 [P] [US4] Write unit tests for RoleRepository in `backend/internal/repository/role_repo_test.go`

**Checkpoint**: All user stories are independently functional. RBAC persistence layer is completed.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Verification and quality checks across all implemented components.

- [ ] T036 Run all validation scenarios defined in `quickstart.md` to ensure DB schema constraints, migrations, and performance are up to par.
- [ ] T037 Run `go vet`, `golangci-lint` and `go test -race ./...` on the `backend` project to ensure zero regressions and high coverage.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - US1 (Phase 3) generates SQL migrations.
  - US2, US3, US4 can be started in parallel or sequentially. Note that Go model/repo implementation in US2/US3/US4 doesn't strictly depend on SQL files existing if using mocks, but tests will require the SQL migrations from US1.
- **Polish (Final Phase)**: Depends on all desired user stories being complete.

### Parallel Opportunities

- All migration SQL scripts within US1 (T007-T014) can be written in parallel.
- All Go model structures within US2, US3, US4 can be written in parallel.
- After repository implementations are done, all unit tests can be written in parallel.

---

## Phase 8: Convergence

- [x] T039 Add JSONB Schema validation for `targeting_rules` payload at database level or app level per FR-014 (missing)
- [x] T040 Update `idx_audit_logs_query` to `idx_audit_logs_project_env_created` in `backend/migrations/000006_create_audit_logs.up.sql` per FR-012 (partial)
- [x] T041 Add migration runner execution on startup in `backend/cmd/server/main.go` per plan (missing)
