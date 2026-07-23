# Tasks: 007-rbac-user-management

**Input**: Design documents from `/specs/007-rbac-user-management/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md, contracts/api.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Define JWT and bcrypt environment variables (or configuration) for the backend
- [x] T002 Add `golang-jwt/jwt` and `golang.org/x/crypto/bcrypt` as Go dependencies
- [x] T003 Create foundational models `backend/internal/models/user.go`, `backend/internal/models/role_assignment.go`, and `backend/internal/models/audit_log.go`
- [x] T004 Create database migration up/down scripts for `users`, `project_role_assignments`, and `audit_logs` tables

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 Implement `internal/auth/password.go` to handle hashing and comparing bcrypt passwords
- [x] T006 Implement `internal/auth/jwt.go` to handle JWT token generation and validation

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - System Authentication (Priority: P1) 🎯 MVP

**Goal**: As a team member, I want to authenticate into the dashboard using a secure login method.

**Independent Test**: Can be fully tested by creating a user in the database, calling the login API, and receiving a JWT.

### Implementation for User Story 1

- [x] T007 [P] [US1] Implement `backend/internal/store/postgres/auth.go` for fetching a User by email
- [x] T008 [US1] Implement `backend/internal/api/auth.go` containing the `POST /api/v1/auth/login` HTTP handler
- [x] T009 [US1] Wire the auth endpoints into the backend router
- [x] T010 [P] [US1] Create frontend API service `frontend/src/services/auth.ts`
- [x] T011 [US1] Build React components `frontend/src/components/auth/LoginForm.tsx` and `frontend/src/pages/Login.tsx`
- [x] T012 [US1] Setup React Router to protect authenticated routes and redirect to `/login` if unauthenticated

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Role-Based Access Control (RBAC) (Priority: P1)

**Goal**: As an administrator, I want to assign roles to restrict who can view, create, or modify feature flags.

**Independent Test**: Can be tested by invoking state-mutating endpoints with a `VIEWER` token and ensuring a 403 Forbidden.

### Implementation for User Story 2

- [x] T013 [P] [US2] Implement `backend/internal/store/postgres/role_assignment.go` to lookup roles for a given `user_id` and `project_id`
- [x] T014 [US2] Implement HTTP middleware in `backend/internal/api/middleware_rbac.go` to enforce RBAC based on the `Authorization` JWT and requested project context
- [x] T015 [US2] Apply the RBAC middleware to all `POST`, `PUT`, `DELETE` endpoints within the project router group

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Immutable Audit Logging (Priority: P2)

**Goal**: As a security auditor, I want to view a history of all flag changes, environment creations, and role assignments.

**Independent Test**: Can be tested by mutating a flag, then querying the audit-logs endpoint to see the old/new states.

### Implementation for User Story 3

- [x] T016 [P] [US3] Implement `backend/internal/store/postgres/audit_log.go` to append new logs and fetch paginated logs
- [x] T017 [US3] Add Audit Log service in `backend/internal/services/audit.go` that explicitly sanitizes `api_key` or other PII from JSONB payloads before storing
- [x] T018 [US3] Update all state-mutating API handlers (e.g., flag toggle, environment creation) to asynchronously dispatch an audit log event
- [x] T019 [US3] Implement HTTP handler `GET /api/v1/projects/{project_id}/audit-logs`
- [x] T020 [US3] Build React page `frontend/src/pages/AuditLogs.tsx` and a display component `frontend/src/components/audit/AuditLogTable.tsx`
- [x] T021 [US3] Update Project Detail layout to include an "Audit Logs" tab pointing to the new page

**Checkpoint**: All user stories 1-3 should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T022 [P] Polish styling, loading skeletons, and error toasts for Login and Audit Logs across components
- [x] T023 Ensure all backend endpoints return standardized error formats (e.g., `401 Unauthorized`, `403 Forbidden`)
- [x] T024 Write or update unit tests for the RBAC middleware and JWT generation
- [x] T025 Run `quickstart.md` validation scenarios manually
