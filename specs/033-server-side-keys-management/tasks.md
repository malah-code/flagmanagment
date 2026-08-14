# Tasks: Server-Side Environment Keys & Tabbed Keys Management

**Input**: Design documents from `/specs/033-server-side-keys-management/`

**Prerequisites**: plan.md, spec.md, data-model.md, contracts/api.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create database migration file for `environment_server_keys` table in `backend/migrations/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

- [ ] T002 Implement `EnvironmentServerKey` model in `backend/internal/models/environment.go`

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Dedicated "Keys" Settings Tab (Priority: P1) 🎯 MVP

**Goal**: Access a dedicated "Keys" tab under Environment Settings to manage client-side credentials in one clear view.

**Independent Test**: Navigate to Environment Settings and verify the "Keys" tab renders the Client-Side Environment Key section.

### Implementation for User Story 1

- [ ] T003 [P] [US1] Create `frontend/src/components/environments/EnvironmentSettingsTabs.tsx` to handle tabbed layout (`General`, `Keys`, `SDK Settings`).
- [ ] T004 [P] [US1] Create `frontend/src/components/environments/ClientSideKeyCard.tsx` extracting existing client-side key UI.
- [ ] T005 [US1] Update environment detail views (e.g. `frontend/src/components/environments/EnvironmentsList.tsx`) to integrate the new `EnvironmentSettingsTabs` and `ClientSideKeyCard`.

**Checkpoint**: User Story 1 should be fully functional with the tabbed UI and Client key visible.

---

## Phase 4: User Story 2 - Server-Side Environment Keys Management (Priority: P1)

**Goal**: Allow admins to create, list, copy, and delete named Server-Side Environment Keys.

**Independent Test**: Create a server-side key via the UI, copy its value, verify it is masked by default, and delete it.

### Implementation for User Story 2

- [ ] T006 [P] [US2] Define Create/List server key request and response structs in `backend/internal/dto/requests.go` and `backend/internal/dto/responses.go`
- [ ] T007 [P] [US2] Implement server keys database operations in `backend/internal/repository/environment_repo.go` (or a dedicated repo file)
- [ ] T008 [US2] Implement business logic (Create, List, Delete) for Server Keys in `backend/internal/services/environment_service.go`
- [ ] T009 [US2] Add HTTP handlers for Server Keys endpoints in `backend/internal/api/handlers/environment_handler.go`
- [ ] T010 [US2] Register Server Keys routes in `backend/internal/api/router.go`
- [ ] T011 [P] [US2] Add ServerKey types in `frontend/src/types/index.ts`
- [ ] T012 [P] [US2] Add ServerKey API calls (create, list, delete) in `frontend/src/services/apiClient.ts`
- [ ] T013 [P] [US2] Create `frontend/src/components/environments/CreateServerKeyDialog.tsx`
- [ ] T014 [US2] Create `frontend/src/components/environments/ServerSideKeysPanel.tsx` listing keys with "Show/Hide", "Copy", and "Revoke" functionality
- [ ] T015 [US2] Integrate `ServerSideKeysPanel` and `CreateServerKeyDialog` into the "Keys" tab of `EnvironmentSettingsTabs.tsx`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently. Server Keys can be fully managed.

---

## Phase 5: User Story 3 - Search and Filter Server-Side Keys (Priority: P2)

**Goal**: Search and filter server-side keys by name in the UI.

**Independent Test**: Enter a query in the search bar and verify only matching server keys are displayed in the table.

### Implementation for User Story 3

- [x] T016 [US3] Add a search bar input and client-side filtering logic to `frontend/src/components/environments/ServerSideKeysPanel.tsx`

**Checkpoint**: All user stories should now be independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T017 Validate backend evaluation endpoint (`backend/internal/middleware/sdk_auth.go`) logic correctly falls back to validating against `environment_server_keys` if the legacy `api_key_hash` on `environments` table does not match.
- [x] T018 Run `quickstart.md` manual tests to ensure full end-to-end functionality.
- [x] T019 Update any necessary frontend documentation or helper text for the Keys tab.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion
- **User Stories (Phase 3+)**: Depend on Foundational phase completion
- **Polish (Final Phase)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Independent frontend work
- **User Story 2 (P1)**: Depends on US1 for UI placement, but backend is independent
- **User Story 3 (P2)**: Depends on US2 frontend

### Parallel Opportunities

- **Backend / Frontend**: T006, T007, T011, T012 can all be run in parallel since they don't depend on each other directly.
- **US1 & US2**: Frontend developers can build `EnvironmentSettingsTabs` (US1) while backend developers implement Server Keys endpoints (US2).

## Implementation Strategy

### Incremental Delivery

1. Setup Database + Backend Model (T001-T002)
2. Add Tabbed UI (T003-T005) -> Deploy/Demo 
3. Build Server Keys Backend (T006-T010)
4. Build Server Keys UI & Integration (T011-T015) -> Deploy/Demo
5. Add Search Filtering (T016) -> Deploy/Demo
6. Final verification & Polish (T017-T019) -> Completion
