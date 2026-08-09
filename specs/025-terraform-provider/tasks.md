---
description: "Task list for Terraform Provider implementation"
---

# Tasks: Terraform Provider

**Input**: Design documents from `/specs/025-terraform-provider/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create Terraform provider directory structure in `providers/terraform`
- [x] T002 Initialize Go module in `providers/terraform/go.mod` with `terraform-plugin-framework` dependencies
- [x] T003 Create `providers/terraform/GNUmakefile` for build and testing automation
- [x] T004 [P] Initialize the main provider entrypoint in `providers/terraform/main.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 Implement Management API client configuration and authentication in `providers/terraform/internal/client/client.go`
- [x] T006 Define the core provider schema (api_url, api_key, bypass_change_requests) in `providers/terraform/internal/provider/provider.go`
- [x] T007 Implement acceptance test setup and provider factory in `providers/terraform/internal/provider/provider_test.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Declarative Infrastructure as Code (Priority: P1) 🎯 MVP

**Goal**: Allow DevOps engineers to define and manage projects, environments, and feature flags.

**Independent Test**: Can run `terraform apply` with a basic configuration to create projects, environments, and flags successfully.

### Implementation for User Story 1

- [x] T008 [P] [US1] Implement Project client API operations in `providers/terraform/internal/client/project.go`
- [x] T009 [P] [US1] Implement Environment client API operations in `providers/terraform/internal/client/environment.go`
- [x] T010 [P] [US1] Implement Feature Flag client API operations in `providers/terraform/internal/client/feature_flag.go`
- [x] T011 [P] [US1] Implement Flag State client API operations in `providers/terraform/internal/client/flag_state.go`
- [x] T012 [US1] Implement the `flagmanagment_project` resource in `providers/terraform/internal/provider/project_resource.go`
- [x] T013 [US1] Implement the `flagmanagment_environment` resource in `providers/terraform/internal/provider/environment_resource.go`
- [x] T014 [US1] Implement the `flagmanagment_feature_flag` resource in `providers/terraform/internal/provider/feature_flag_resource.go`
- [x] T015 [US1] Implement the `flagmanagment_flag_state` resource in `providers/terraform/internal/provider/flag_state_resource.go`
- [x] T016 [US1] Implement Change Request bypass logic during environment and flag state updates

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Managing Governance and RBAC via Terraform (Priority: P2)

**Goal**: Provision API keys, service accounts, and RBAC permissions using Terraform.

**Independent Test**: Can be tested by provisioning a service account and API key via Terraform, then using that API key to query the FlagManagment REST API.

### Implementation for User Story 2

- [x] T017 [P] [US2] Implement Service Account client API operations in `providers/terraform/internal/client/service_account.go`
- [x] T018 [US2] Implement the `flagmanagment_service_account` resource in `providers/terraform/internal/provider/service_account_resource.go`
- [x] T019 [US2] Ensure sensitive output masking for Service Account tokens in Terraform state

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Telemetry and Rollback Triggers Integration (Priority: P3)

**Goal**: Define telemetry triggers and automated rollback actions on specific flags via Terraform.

**Independent Test**: Provisioning a flag with attached telemetry triggers using Terraform successfully persists to the backend.

### Implementation for User Story 3

- [x] T020 [P] [US3] Extend the Flag State client in `providers/terraform/internal/client/flag_state.go` to support telemetry payload fields
- [x] T021 [US3] Extend the `flagmanagment_flag_state` resource in `providers/terraform/internal/provider/flag_state_resource.go` to support `telemetry_trigger` blocks
- [x] T022 [US3] Add validation for telemetry trigger conditions within the resource schema

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T023 [P] Update `README.md` in `providers/terraform/` with usage and installation instructions
- [x] T024 Add end-to-end acceptance tests utilizing `terraform-plugin-testing` across all resources in `providers/terraform/internal/provider/`
- [x] T025 Execute `quickstart.md` validation workflow locally

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### Parallel Opportunities

- Foundation API client (`T005`) and Provider schema (`T006`) can be written mostly in parallel.
- US1 Client API operations (`T008` to `T011`) can be developed simultaneously.
- Provider resources in US1 (`T012` to `T015`) can be built in parallel once their respective client implementations exist.

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently using a real Terraform apply with the HCL configuration in the contracts document.
