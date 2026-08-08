# Feature Specification: One-Click Flag Environment Promotions

**Feature Branch**: `011-flag-promotions`

**Created**: 2026-07-25

**Status**: Draft

**Input**: User description: "Implement one-click feature flag configuration promotions between environments."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Promote Flag Configuration to Target Environment (Priority: P1)

As a Release Manager, I want to copy a feature flag's state, targeting rules, and remote config from a source environment (e.g., QA) to a target environment (e.g., Staging or Production) with one click, so that configurations are promoted safely without manual errors.

**Why this priority**: Core productivity and governance feature to prevent configuration drift across environments.

**Independent Test**: Can be tested by selecting a flag in QA, clicking "Promote to Staging", and verifying that Staging receives the exact ruleset from QA.

**Acceptance Scenarios**:

1. **Given** a flag configured in a source environment, **When** a user promotes the flag to an un-protected target environment, **Then** the target environment's flag state is updated immediately to match the source environment configuration.
2. **Given** the target environment is marked as "Protected", **When** a user promotes the flag, **Then** the system automatically generates a Change Request for approval instead of applying the state immediately.

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST expose a Promotion API endpoint `POST /api/v1/projects/{projectId}/flags/{flagId}/promote` accepting source and target environment IDs.
- **FR-002**: System MUST copy flag state (`enabled`), targeting rules (`JSONB`), and remote config (`JSONB`) from source to target.
- **FR-003**: System MUST respect protected target environments by generating a pending `ChangeRequest` if the target environment is protected.
- **FR-004**: System MUST record an Audit Log entry for the promotion action.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Promotion execution completes in under 100ms.
- **SC-002**: Zero configuration drift (100% attribute copy precision) between source and target environments.

## Assumptions

- Users must have `EDITOR` role on the target environment to initiate a promotion.
