# Tasks: User Management API Implementation

**Input**: Design documents from `/specs/034-user-management-api/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/api.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Phase 1: Setup

**Purpose**: Project initialization and basic structure. As this is an existing project, setup is minimal.

- [ ] T001 Initialize database migrations for new entities in `backend/migrations/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [ ] T002 Implement Go models (`User`, `ProjectAccess`, `Invitation`, `SystemConfig`) in `backend/internal/models/`
- [ ] T003 Implement DB migrations for the new tables in `backend/migrations/`

**Checkpoint**: Foundation ready - user story implementation can now begin.

---

## Phase 3: User Story 1 - View Team Members (Priority: P1) 🎯 MVP

**Goal**: Administrators can see a list of all team members, their global roles, and project access.

**Independent Test**: Load the Team Settings page and verify that the displayed users match the database records instead of mock data.

### Implementation for User Story 1

- [ ] T004 [P] [US1] Implement `UserService.GetUsers()` in `backend/internal/services/user_service.go`
- [ ] T005 [US1] Implement `GET /api/v1/users` endpoint in `backend/internal/handlers/user_handler.go`
- [ ] T006 [P] [US1] Create React `useUsers` hook in `frontend/src/hooks/useUsers.ts` to call the new API
- [ ] T007 [US1] Integrate `useUsers` hook into `frontend/src/pages/UsersManagement.tsx`, removing the mock listing data.

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently.

---

## Phase 4: User Story 2 - Invite New User (Priority: P1)

**Goal**: Administrators can invite new colleagues by email and assign them initial roles/projects.

**Independent Test**: Submit an invitation and verify that an email is dispatched (via MailHog) and the user appears in the list as "Pending".

### Implementation for User Story 2

- [ ] T008 [P] [US2] Implement `CryptoService.GenerateToken()` in `backend/internal/services/crypto_service.go`
- [ ] T009 [P] [US2] Implement `EmailService.SendInvitation()` in `backend/internal/services/email_service.go`
- [ ] T010 [US2] Implement `UserService.CreateInvitation()` in `backend/internal/services/user_service.go` (depends on T008)
- [ ] T011 [US2] Implement `POST /api/v1/users/invite` in `backend/internal/handlers/user_handler.go` (depends on T009, T010)
- [ ] T012 [US2] Update `frontend/src/pages/UsersManagement.tsx` to call API on invite form submit (replace mock invite).

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently.

---

## Phase 5: User Story 3 - Edit User Roles and Access (Priority: P2)

**Goal**: Administrators can modify a user's role or assigned projects.

**Independent Test**: Modify a user's role and verify the backend correctly applies and reflects the new permissions.

### Implementation for User Story 3

- [ ] T013 [P] [US3] Implement `UserService.UpdateUserAccess()` in `backend/internal/services/user_service.go`
- [ ] T014 [US3] Implement `PUT /api/v1/users/:id/access` in `backend/internal/handlers/user_handler.go`
- [ ] T015 [US3] Update `frontend/src/pages/UsersManagement.tsx` to call API on edit save (replace mock edit).

**Checkpoint**: All user management capabilities are fully functional.

---

## Phase 6: User Story 4 - Configure System Email Server (Priority: P2)

**Goal**: System Administrator can configure the outbound SMTP server settings within the platform.

**Independent Test**: Save SMTP credentials and successfully trigger a "Test Connection" email from the UI.

### Implementation for User Story 4

- [ ] T016 [P] [US4] Implement `CryptoService.EncryptAES()` and `DecryptAES()` in `backend/internal/services/crypto_service.go`
- [ ] T017 [US4] Implement `GET /api/v1/config/smtp` and `PUT /api/v1/config/smtp` in `backend/internal/handlers/config_handler.go` (depends on T016)
- [ ] T018 [US4] Implement `POST /api/v1/config/smtp/test` in `backend/internal/handlers/config_handler.go`
- [ ] T019 [P] [US4] Create React `useConfig` hook in `frontend/src/hooks/useConfig.ts`
- [ ] T020 [US4] Implement `frontend/src/pages/SystemSettings.tsx` with a form to manage SMTP config.

**Checkpoint**: All user stories should now be independently functional.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories.

- [ ] T021 Code cleanup, removing any unused mock code from the frontend
- [ ] T022 Run `quickstart.md` validation end-to-end

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately.
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories.
- **User Stories (Phase 3+)**: All depend on Foundational phase completion. User stories can proceed in parallel or sequentially.
- **Polish (Final Phase)**: Depends on all user stories being complete.

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2).
- **User Story 2 (P1)**: Can start after Foundational. Needs US4 functionality for real emails, but MailHog can be hardcoded for testing until US4 is built.
- **User Story 3 (P2)**: Can start after Foundational.
- **User Story 4 (P2)**: Can start after Foundational.

### Parallel Opportunities

- Foundational DB migrations (T001, T003) and Models (T002) can run sequentially.
- React hooks (T006, T019) can be built in parallel with backend endpoints.
- Backend services (UserService, EmailService, CryptoService) can be built in parallel.
