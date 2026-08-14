# Tasks: Competitor UX Parity

**Input**: Design documents from `/specs/031-flagcompatitor-ux-parity/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

*(No setup required as this is a frontend-only feature in an existing project)*

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

*(No foundational backend or infrastructure tasks needed, the frontend hooks already exist)*

---

## Phase 3: User Story 1 - One-Click Flag Toggling (Priority: P1) 🎯 MVP

**Goal**: Implement a prominent toggle switch directly on the flag listing rows when viewed within an environment context.

**Independent Test**: Can be verified by rendering the flag list and observing a large toggle switch for boolean states that responds instantly.

### Implementation for User Story 1

- [x] T001 [US1] Create or update a robust Toggle/Switch UI component in `frontend/src/components/ui/Switch.tsx` (or similar) using Tailwind.
- [x] T002 [US1] Integrate the toggle switch into the flag list row in `frontend/src/components/flags/FlagsList.tsx`.
- [x] T003 [US1] Integrate the toggle switch into the environment-specific flag list row in `frontend/src/components/flagStates/FlagStatesList.tsx`.
- [x] T004 [US1] Ensure the toggle properly invokes `useUpdateFlagState` and handles optimistic updates/loading states.

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Dedicated Left Sidebar for Environments (Priority: P1)

**Goal**: Show all available environments in a persistent, aesthetic left sidebar for quick context switching.

**Independent Test**: Can be verified by observing the project dashboard layout containing a left-aligned navigation pane listing environments.

### Implementation for User Story 2

- [x] T005 [US2] Create the sidebar navigation component in `frontend/src/components/layout/EnvironmentSidebar.tsx`.
- [x] T006 [US2] Update the main project dashboard layout in `frontend/src/pages/ProjectDashboard.tsx` to include the `EnvironmentSidebar`.
- [x] T007 [US2] Implement active state styling in the sidebar based on the currently selected environment URL or state.
- [x] T008 [US2] Remove the redundant environment dropdown/tabs from the main view now that the sidebar handles environment selection.

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Optional Initial Value on Flag Creation (Priority: P2)

**Goal**: Optionally specify an initial value when creating a flag, so users do not have to navigate to targeting rules immediately after creation.

**Independent Test**: Can be verified by opening the "Create Flag" dialog and observing an input field for the initial value.

### Implementation for User Story 3

- [x] T009 [US3] Add an "Initial Value" (boolean/string/number) input field to the form in `frontend/src/components/flags/CreateFlagDialog.tsx`.
- [x] T010 [US3] Update the API payload in `useFlags.ts` (or relevant hook) to include the `initial_value` when creating the flag.

**Checkpoint**: All P1 and P2 user stories should now be independently functional

---

## Phase 6: User Story 4 - Tag Management and Filtering (Priority: P2)

**Goal**: Assign tags to flags and filter the flag list by selecting these tags.

**Independent Test**: Can be verified by adding a tag to a flag and subsequently filtering the list to show only flags matching that tag.

### Implementation for User Story 4

- [x] T011 [US4] Create a tag multi-select/filter component in `frontend/src/components/flags/TagFilter.tsx`.
- [x] T012 [US4] Integrate the `TagFilter` component into `frontend/src/components/flags/FlagsList.tsx` above the table.
- [x] T013 [US4] Implement local filtering logic in `FlagsList.tsx` to only display flags that match the selected tags.
- [x] T014 [US4] (If not already present) Add a UI to apply tags to flags in the flag creation/edit dialogs or detail views.

**Checkpoint**: All user stories should now be independently functional

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T015 [P] Run `quickstart.md` validation scenarios.
- [x] T016 [P] Verify responsive behavior of the new sidebar layout on mobile viewports.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: N/A
- **Foundational (Phase 2)**: N/A
- **User Stories (Phase 3+)**: Can proceed in parallel as they touch different parts of the UI (toggles vs. layout vs. creation modal).
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: No dependencies.
- **User Story 2 (P1)**: No dependencies.
- **User Story 3 (P2)**: No dependencies.
- **User Story 4 (P2)**: Depends on the existence of the flag list to apply the filter.

### Parallel Opportunities

- US1, US2, US3, and US4 can all be implemented in parallel by different developers since they affect mostly isolated React components.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 3: User Story 1 (Flag Toggles).
2. **STOP and VALIDATE**: Test User Story 1 independently.

### Incremental Delivery

1. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
2. Add User Story 2 (Sidebar) → Test independently → Deploy/Demo
3. Add User Story 3 (Initial Value) → Test independently → Deploy/Demo
4. Add User Story 4 (Tag Filter) → Test independently → Deploy/Demo
