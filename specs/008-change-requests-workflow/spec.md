# Feature Specification: Change Requests Workflow

**Feature Branch**: `[008-change-requests-workflow]`

**Created**: 2026-07-25

**Status**: Draft

**Input**: User description: "Implement Change Request workflow for protected environments (Production), including state machine (Pending -> Approved -> Applied / Rejected), visual JSON diff generator, and approval enforcement by Release Managers."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create Change Request in Protected Environment (Priority: P1)

As a team member, I want to propose a flag configuration change in a protected environment, so that it can be safely reviewed before affecting production users.

**Why this priority**: Core mechanism to prevent accidental or unauthorized changes in critical environments.

**Independent Test**: Can be tested by attempting to toggle a flag in an environment marked as "Protected" and verifying a Change Request is generated instead of an immediate state change.

**Acceptance Scenarios**:

1. **Given** an environment marked as protected, **When** I attempt to change a feature flag's state or targeting rules, **Then** my changes are saved as a Pending Change Request and the live flag remains unchanged.
2. **Given** a non-protected environment, **When** I attempt to change a feature flag, **Then** the change applies immediately without generating a Change Request.

---

### User Story 2 - Review and Approve Change Request (Priority: P1)

As a Release Manager, I want to review a visual diff of proposed changes and approve them, so that safe configurations are atomically applied to the protected environment.

**Why this priority**: Required to complete the lifecycle and actually deploy changes.

**Independent Test**: Can be tested by logging in as a Release Manager, viewing the pending request diff, and clicking Approve to verify the changes go live.

**Acceptance Scenarios**:

1. **Given** a Pending Change Request, **When** a Release Manager reviews the visual JSON diff and clicks Approve, **Then** the state transitions to Approved and the changes are atomically applied to the live environment.
2. **Given** a Pending Change Request created by myself, **When** I attempt to approve it, **Then** the system prevents self-approval and requires a different Release Manager.

---

### User Story 3 - Reject Change Request (Priority: P2)

As a Release Manager, I want to reject an unsafe or incorrect Change Request, so that bad configurations do not reach production.

**Why this priority**: Critical for governance to discard invalid proposals cleanly.

**Independent Test**: Can be tested by rejecting a Pending request and verifying it is marked as Rejected without altering the live environment.

**Acceptance Scenarios**:

1. **Given** a Pending Change Request, **When** a Release Manager rejects it with a reason, **Then** the state transitions to Rejected and the proposed changes are discarded.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow any Environment to be toggled as "Protected".
- **FR-002**: System MUST intercept state-mutating requests (toggle, rule updates) on flags in Protected environments and convert them into Change Requests.
- **FR-003**: System MUST track Change Requests through a strict state machine: `Pending` -> `Approved` -> `Applied`, or `Pending` -> `Rejected`.
- **FR-004**: System MUST generate and provide a visual diff comparing the current live configuration against the proposed JSON configuration.
- **FR-005**: System MUST enforce that only users with the `Release Manager` (or `Admin`) role can Approve or Reject a Change Request.
- **FR-006**: System MUST prevent users from approving their own Change Requests (Self-approval restriction).
- **FR-007**: System MUST atomically apply the proposed changes to the live environment immediately upon successful Approval.
- **FR-008**: System MUST record all Change Request creations, approvals, rejections, and applications to the Immutable Audit Log.

### Key Entities 

- **Environment**: Enhanced with an `is_protected` boolean flag.
- **ChangeRequest**: Represents a proposed mutation. Stores the target environment ID, flag ID, author user ID, proposed state (JSON), current state (JSON), and status.
- **Approval**: Represents the decision (Approve/Reject), the reviewing user ID, and an optional comment/reason.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of mutations directed at protected environments successfully route into the Change Request workflow rather than applying instantly.
- **SC-002**: Self-approvals are prevented by the backend API 100% of the time, returning a 403 Forbidden.
- **SC-003**: The visual diff UI accurately highlights additions, deletions, and modifications between the old and new JSON targeting rules.
- **SC-004**: Once Approved, the transition to Applied and the live environment mutation happen in a single, atomic database transaction to prevent partial state updates.

## Assumptions

- The frontend UI will utilize a standard JSON diffing library (e.g., `react-diff-viewer`) to render the visual differences.
- If the `Release Manager` role does not explicitly exist yet in the DB, it will be added as part of this feature's implementation, alongside `Admin`, `Editor`, and `Viewer`.
- Approval and Application happen simultaneously in this MVP (i.e., Scheduled Deployments are out of scope for now; approving applies the change instantly).
