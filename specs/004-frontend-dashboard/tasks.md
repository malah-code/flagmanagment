# Tasks: Frontend Dashboard

**Input**: Design documents from `/specs/004-frontend-dashboard/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Install dependencies (TailwindCSS, react-router-dom, @tanstack/react-query, lucide-react, clsx, tailwind-merge) in frontend
- [x] T002 Configure Tailwind CSS and postcss in `frontend/tailwind.config.js` and `frontend/postcss.config.js`
- [x] T003 Initialize Shadcn/UI and add base components (button, input, card, dialog, table, label, switch, toast)
- [x] T004 Setup global styles in `frontend/src/index.css`
- [x] T005 Setup React Router and React Query Provider in `frontend/src/App.tsx` and `frontend/src/main.tsx`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T006 Define shared TypeScript types in `frontend/src/types/index.ts` (Project, Environment, FeatureFlag, FlagState)
- [x] T007 Create base API client with fetch wrapper in `frontend/src/services/apiClient.ts`
- [x] T008 Create generic app layout and navigation header in `frontend/src/components/shared/Layout.tsx`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Project Management (Priority: P1) 🎯 MVP

**Goal**: As a product manager, I want to view, create, and edit projects.

**Independent Test**: Can be fully tested by creating a new project in the UI and verifying it appears in the projects list.

### Implementation for User Story 1

- [x] T009 [P] [US1] Create Project API service in `frontend/src/services/projects.ts`
- [x] T010 [US1] Create React Query hooks for projects in `frontend/src/hooks/useProjects.ts`
- [x] T011 [US1] Create Projects list page in `frontend/src/pages/ProjectsList.tsx`
- [x] T012 [US1] Create New Project dialog form in `frontend/src/components/projects/CreateProjectDialog.tsx`
- [x] T013 [US1] Wire up `/projects` route to `ProjectsList` in `frontend/src/App.tsx`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Environment Management (Priority: P1)

**Goal**: As a developer, I want to view and create environments within a project.

**Independent Test**: Can be tested by navigating to a project, creating a new environment, and verifying the securely generated API key is presented.

### Implementation for User Story 2

- [x] T014 [P] [US2] Create Environment API service in `frontend/src/services/environments.ts`
- [x] T015 [US2] Create Environment React Query hooks in `frontend/src/hooks/useEnvironments.ts`
- [x] T016 [US2] Create Project Dashboard layout with sidebar/tabs in `frontend/src/pages/ProjectDetail.tsx`
- [x] T017 [US2] Create Environments list component in `frontend/src/components/environments/EnvironmentsList.tsx`
- [x] T018 [US2] Create New Environment dialog with one-time API key display in `frontend/src/components/environments/CreateEnvironmentDialog.tsx`
- [x] T019 [US2] Route `/projects/:projectId` to `ProjectDetail` in `frontend/src/App.tsx`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Feature Flag Management (Priority: P1)

**Goal**: As a developer, I want to create and define feature flags for a project.

**Independent Test**: Can be tested by adding a new flag with a specific key and type to a project.

### Implementation for User Story 3

- [x] T020 [P] [US3] Create Feature Flag API service in `frontend/src/services/flags.ts`
- [x] T021 [US3] Create Feature Flag React Query hooks in `frontend/src/hooks/useFlags.ts`
- [x] T022 [US3] Create Feature Flags list component in `frontend/src/components/flags/FlagsList.tsx`
- [x] T023 [US3] Create New Feature Flag dialog in `frontend/src/components/flags/CreateFlagDialog.tsx`
- [x] T024 [US3] Add Flags tab to `ProjectDetail` in `frontend/src/pages/ProjectDetail.tsx`

**Checkpoint**: All user stories 1-3 should now be independently functional

---

## Phase 6: User Story 4 - Flag State and Rules Configuration (Priority: P1)

**Goal**: As a release manager, I want to toggle flag states within a specific environment.

**Independent Test**: Can be tested by selecting an environment, toggling a flag ON/OFF, and verifying the state persists.

### Implementation for User Story 4

- [x] T025 [P] [US4] Create Flag State API service in `frontend/src/services/flagStates.ts`
- [x] T026 [US4] Create Flag State React Query hooks in `frontend/src/hooks/useFlagStates.ts`
- [x] T027 [US4] Create Flag States list component with toggle switches in `frontend/src/components/flagStates/FlagStatesList.tsx`
- [x] T028 [US4] Update `ProjectDetail` layout to allow selecting an active environment to view its flag states

**Checkpoint**: All user stories 1-4 should now be independently functional

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T029 [P] Polish styling, loading skeletons, and error toasts across components
- [x] T030 Add basic rendering tests for key components in `frontend/tests/` to meet Constitution test goals
- [x] T031 Run `quickstart.md` validation scenarios manually
