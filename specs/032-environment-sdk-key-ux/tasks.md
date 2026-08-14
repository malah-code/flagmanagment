# Tasks: Environment SDK Key UX & Integration Guide

**Input**: Design documents from `/specs/032-environment-sdk-key-ux/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, quickstart.md

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

*(No shared infrastructure setup needed)*

---

## Phase 2: Foundational (Blocking Prerequisites)

*(No blocking backend foundational tasks required)*

---

## Phase 3: User Story 1 - Always-Visible Client SDK Key (Priority: P1) 🎯 MVP

**Goal**: Display the Client SDK Key in Environment Settings cards with a 1-click copy button.

**Independent Test**: Can be verified by navigating to Environment Settings and clicking the copy button next to an environment card.

### Implementation for User Story 1

- [ ] T001 [US1] Update `Environment` interface in `frontend/src/types/index.ts` to ensure `apiKey` / client key is typed.
- [ ] T002 [US1] Update `frontend/src/components/environments/EnvironmentsList.tsx` to render the Client SDK Key on environment cards with a 1-click copy button and feedback toast.

**Checkpoint**: User Story 1 is functional and testable independently.

---

## Phase 4: User Story 2 - Interactive SDK Integration Guide (Priority: P1)

**Goal**: Create an interactive Integration Guide modal showing ready-to-use code snippets for React, Node.js, Python, and Go pre-filled with the environment key.

**Independent Test**: Can be verified by clicking the `< /> Integration` button on any environment card and switching code tabs.

### Implementation for User Story 2

- [x] T003 [US2] Create the `SDKIntegrationModal.tsx` component in `frontend/src/components/environments/SDKIntegrationModal.tsx` with tabs for React, Node.js, Python, and Go.
- [x] T004 [US2] Add the `< /> Integration` trigger button to environment cards in `frontend/src/components/environments/EnvironmentsList.tsx` to open `SDKIntegrationModal`.
- [x] T005 [US2] Add 1-click "Copy Code Snippet" functionality inside `SDKIntegrationModal.tsx`.

**Checkpoint**: User Stories 1 and 2 work independently.

---

## Phase 5: User Story 3 - Distinct Public vs. Private Key Management (Priority: P2)

**Goal**: Clearly distinguish between Public Client SDK Keys and Private Admin Keys using visual badges.

**Independent Test**: Can be verified by checking the Environment Settings UI for "Client-Side / Public Key" badges.

### Implementation for User Story 3

- [ ] T006 [US3] Add "Client-Side / Public Key" badges and tooltip helpers to `EnvironmentsList.tsx` and `SDKIntegrationModal.tsx`.

**Checkpoint**: All user stories are complete.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T007 [P] Run `quickstart.md` validation scenarios.
- [ ] T008 [P] Verify responsive behavior of `SDKIntegrationModal` on mobile screens.

---

## Dependencies & Execution Order

- US1, US2, and US3 can proceed in sequence or in parallel as they refine the environment settings components.
- Polish phase depends on user stories.
