# Feature Specification: API-Driven Environment Cloning (Ephemeral Environments)

**Feature Branch**: `[018-environment-cloning]`

**Created**: 2026-08-08

**Status**: Draft

**Input**: User description: "what's next? taking into consideration the initial requirements @[/wsl+ubuntu/home/tarikelmallah/Projects/FlagManagment/docs/g-requirements.md]@[/wsl+ubuntu/home/tarikelmallah/Projects/FlagManagment/docs/p-requirements.md]"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - CI/CD Ephemeral Environment Creation (Priority: P1)

As a DevOps Engineer integrating FlagManagment into our CI/CD pipelines, I want to programmatically clone an existing environment (e.g., "Production" or "Staging") into a new, temporary environment so that automated integration tests can run against isolated, production-like flag states.

**Why this priority**: Fast, isolated testing is critical for modern engineering teams. Without API-driven cloning, teams are forced to test in shared QA environments, leading to test flakiness and collisions. This fulfills PRD Section 10.2.

**Independent Test**: Can be tested by invoking the clone API endpoint with a source environment ID, asserting that a new environment is created with a distinct SDK token, and verifying that all feature flag states from the source are perfectly replicated in the new environment.

**Acceptance Scenarios**:

1. **Given** an existing "Staging" environment with 50 flags configured, **When** I issue a clone request via the API with the name `PR-123-Test`, **Then** a new environment is created, a new SDK authentication token is generated and returned, and all 50 flag states are exactly mirrored into the new environment.
2. **Given** a new ephemeral environment, **When** a CI/CD test toggles a flag within it, **Then** the change only applies to the ephemeral environment and does not impact the source "Staging" environment.

---

### User Story 2 - Ephemeral Environment Teardown (Priority: P1)

As a DevOps Engineer, I want to programmatically delete ephemeral environments once my CI/CD pipeline completes so that I do not clutter the FlagManagment dashboard with hundreds of dead test environments.

**Why this priority**: Required to prevent database bloat and UI clutter when generating environments dynamically on every pull request.

**Independent Test**: Can be tested by deleting an environment via the API and verifying that its SDK token is revoked and its flag states are removed, without affecting other environments.

**Acceptance Scenarios**:

1. **Given** an ephemeral environment `PR-123-Test`, **When** I issue a delete request via the API, **Then** the environment is permanently removed, its API key is invalidated, and all associated environment flag states are deleted.
2. **Given** a protected environment (e.g., "Production"), **When** I issue a delete request via the API, **Then** the system rejects the deletion with a 403 Forbidden to prevent catastrophic data loss.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST expose a REST API endpoint `POST /api/v1/projects/{projectId}/environments/{sourceEnvId}/clone` that creates a new environment based on the source.
- **FR-002**: The cloning process MUST copy all targeting rules, remote configuration payloads, boolean states, and rollout percentages from the source environment into the new environment.
- **FR-003**: The cloned environment MUST generate and return a unique, cryptographically secure SDK authentication token.
- **FR-004**: The system MUST expose a REST API endpoint `DELETE /api/v1/projects/{projectId}/environments/{envId}` to permanently delete an environment.
- **FR-005**: Environments marked as "Protected" MUST NOT be deletable unless their protection status is first removed by an authorized user.
- **FR-006**: Both the cloning and deletion of environments MUST be recorded in the immutable audit log, capturing the actor ID, timestamp, and target environment details.

### Key Entities

- **Environment**: The core entity being duplicated. Includes metadata and a secure SDK token.
- **EnvironmentFlagState**: The state matrices that must be deep-copied from the source environment to the destination environment.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An environment containing 1,000 feature flags can be cloned completely in under 2 seconds.
- **SC-002**: Automated scripts can successfully provision, test against, and destroy an ephemeral environment with a 100% success rate.
- **SC-003**: Deleting an ephemeral environment immediately invalidates its SDK token, terminating any connected server-side SDK streams within 5 seconds.

## Assumptions

- **Concurrency**: Environment cloning is relatively infrequent compared to flag evaluation, so the cloning operation can rely on standard PostgreSQL transactions without requiring specialized queuing or asynchronous worker processes.
- **Permissions**: The API token or user invoking the clone/delete endpoints possesses the appropriate Project Owner or System Administrator RBAC permissions.
