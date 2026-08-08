# Feature Specification: Stale Flag Detection & Lifecycle Management

**Feature Branch**: `015-stale-flag-management`

**Created**: 2026-08-08

**Status**: Draft

**Input**: User description: "Stale Flag Management per Section 12.1 and Appendix A.3 of PRD requirements (docs/p-requirements.md and docs/g-requirements.md). Automatically detect, highlight, and manage stale feature flags that have remained unchanged or fully rolled out for 30+ days to reduce technical debt."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Automatic Stale Flag Detection & Tracking (Priority: P1)

As a Release Manager or Developer, I want the system to automatically track evaluation activity and state stability for all feature flags across environments, so that flags that have remained unchanged at 100% rollout or inactive for 30+ days are automatically detected and flagged.

**Why this priority**: Stale flags accumulate code debt and create cognitive overhead for teams. Automated tracking is the foundation of technical debt reduction and prevents obsolete code paths from lingering in production.

**Independent Test**: Can be tested by simulating flag evaluation events and timestamp aging, verifying that the system correctly transitions flag status to `STALE` when criteria are met.

**Acceptance Scenarios**:

1. **Given** a feature flag has been at 100% rollout with no state modifications for 30 consecutive days, **When** the stale flag evaluation job runs, **Then** the flag's lifecycle status for that environment is automatically updated to `STALE`.
2. **Given** a feature flag has received zero evaluation requests across all SDKs for 30 consecutive days, **When** the system scans flag activity, **Then** the flag is flagged as `INACTIVE_STALE` with the last evaluation timestamp recorded.
3. **Given** a flag marked as `STALE` undergoes a rule or targeting modification by a user, **When** the change is saved, **Then** the flag's lifecycle status automatically reverts to `ACTIVE` and its staleness timer resets.

---

### User Story 2 - Stale Flag Dashboard & Lifecycle Actions (Priority: P2)

As a Project Owner or Release Manager, I want a dedicated dashboard view and filtering options to review all stale flags across projects and environments, and take lifecycle actions such as archiving, deprecating, or confirming flags.

**Why this priority**: Detection alone is insufficient without actionable UX tools to review stale flags and safely archive them once code paths are cleaned up.

**Independent Test**: Can be tested by navigating to the flag list, filtering by `STALE` status, and executing an `ARCHIVE` or `DEPRECATE` action, verifying that the flag is archived and hidden from active evaluation rulesets.

**Acceptance Scenarios**:

1. **Given** a user is viewing the feature flag list, **When** they apply the `STALE` filter, **Then** only flags categorized as stale are displayed along with their staleness reason (e.g., "100% rollout for 45 days", "No evaluation in 60 days").
2. **Given** a Release Manager is inspecting a stale flag, **When** they select the "Archive Flag" action and confirm, **Then** the flag's status transitions to `ARCHIVED`, it is excluded from active SDK rulesets, and an audit log entry is recorded.
3. **Given** an archived feature flag, **When** an SDK requests evaluation for that flag, **Then** the SDK returns the default fallback value without error.

---

### User Story 3 - Configurable Staleness Policies (Priority: P3)

As a Project Owner, I want to configure custom staleness thresholds (e.g., 14 days for Dev, 60 days for Production) per project or environment, so that team-specific SLA rules dictate when flags are flagged as stale.

**Why this priority**: Different environments and projects have varying release cadence expectations; fixed 30-day thresholds may not fit fast-paced dev environments or slow enterprise release cycles.

**Independent Test**: Can be tested by updating a project's staleness threshold from 30 days to 14 days and verifying that flags unchanged for 15 days are immediately categorized as stale under the new policy.

**Acceptance Scenarios**:

1. **Given** a Project Owner is on the Project Settings page, **When** they update the staleness threshold to 14 days for Development and save, **Then** flags in Development unchanged for >14 days are marked as stale on the next scan.

---

### Edge Cases

- What happens when a stale flag is referenced by an active sequential flag dependency?
  - The system MUST prevent archiving the parent flag until all child dependent flags are resolved or detached.
- How does the system handle high-volume evaluation timestamp updates without database write bottlenecks?
  - Timestamp updates MUST be aggregated and batch-flushed to persistent storage rather than executing a synchronous database write per evaluation request.
- What happens if an archived flag is restored?
  - Restoring an archived flag transitions its status back to `ACTIVE` and re-includes it in environment SDK rulesets.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST record and maintain a `last_evaluated_at` timestamp for every feature flag per environment.
- **FR-002**: System MUST automatically evaluate flag staleness against configured threshold policies (default: 30 days of inactivity or 100% unchanged rollout).
- **FR-003**: System MUST support lifecycle states for feature flags: `ACTIVE`, `STALE`, `DEPRECATED`, and `ARCHIVED`.
- **FR-004**: System MUST allow filtering and searching feature flags by lifecycle state across global, project, and environment views.
- **FR-005**: System MUST allow authorized users (Release Managers and Project Owners) to perform lifecycle transitions (`MARK_STALE`, `DEPRECATE`, `ARCHIVE`, `RESTORE`).
- **FR-006**: System MUST prevent archiving a feature flag if active dependent flags still rely on it as a parent dependency.
- **FR-007**: System MUST record an immutable audit log entry for every manual or automated flag lifecycle state change.
- **FR-008**: System MUST allow Project Owners to customize staleness thresholds (in days) per project and environment.

### Key Entities

- **FlagLifecycleState**: Enum representing the current lifecycle stage (`ACTIVE`, `STALE`, `DEPRECATED`, `ARCHIVED`).
- **StaleFlagPolicy**: Configuration entity defining staleness criteria (threshold in days, inclusion of 100% rollout vs zero evaluation) per project/environment.
- **FlagEvaluationMetric**: Aggregate entity tracking evaluation counts and `last_evaluated_at` timestamps for flags across environments.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Stale flags are detected and highlighted in the UI within 1 hour of surpassing the configured staleness threshold.
- **SC-002**: Engineering teams can identify and filter all stale flags in a project in under 3 seconds.
- **SC-003**: 100% of flag lifecycle state transitions (stale detection, archiving, restoration) are captured in the immutable audit log.
- **SC-004**: Archiving a stale flag immediately removes it from active SDK rulesets within the normal delta synchronization window (under 1 second).

## Assumptions

- Flag evaluation metrics are aggregated in-memory or via cache (e.g., Redis) before batch persistence to prevent database write saturation.
- Default staleness criteria is 30 days of no state modifications while at 100% rollout, or 30 days without evaluation requests.
- Archived flags remain stored for historical audit compliance but are excluded from active evaluation streams.
