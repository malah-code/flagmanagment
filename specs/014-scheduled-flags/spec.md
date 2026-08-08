# Feature Specification: Scheduled Flag Changes

**Feature Branch**: `014-scheduled-flags`

**Created**: 2026-07-25

**Status**: Draft

**Input**: User description: "what is the next unimplemented feature to work with" (Interpreted as: "The ability to schedule a feature flag to turn on or off at a specific date and time in the future, or schedule a change request to apply automatically at a given time.")

## Clarifications

### Session 2026-07-25
- Q: How should the system handle scheduling conflicts when multiple schedules exist for the same flag? → A: Reject new conflicting schedules

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Schedule Flag State Change (Priority: P1)

As a Release Manager, I want to schedule a flag to turn on at a specific future date and time so that marketing launches can happen automatically without my manual intervention.

**Why this priority**: It is the core capability of the feature, providing immediate business value for coordinated product launches.

**Independent Test**: Can be fully tested by creating a schedule for a flag, waiting for the time to pass, and verifying the flag state updates automatically.

**Acceptance Scenarios**:

1. **Given** a flag is currently OFF, **When** I schedule it to turn ON tomorrow at 9:00 AM, **Then** the flag remains OFF today, and automatically turns ON tomorrow at exactly 9:00 AM.
2. **Given** I have a scheduled flag change, **When** I view the flag details, **Then** I see the upcoming scheduled change clearly indicated.
3. **Given** a scheduled flag change in the future, **When** I decide to cancel the schedule, **Then** the schedule is removed and the flag state remains unchanged.

---

### User Story 2 - Schedule Change Request Approval (Priority: P2)

As a Release Manager, I want to schedule an approved Change Request to be automatically applied at a specific future date and time.

**Why this priority**: Enhances the existing change request workflow with scheduling capabilities, useful for complex configuration rollouts during off-hours.

**Independent Test**: Can be tested by creating a change request, approving it, and attaching a scheduled application time instead of applying it immediately.

**Acceptance Scenarios**:

1. **Given** an approved Change Request, **When** I schedule its application for Saturday at 2:00 AM, **Then** the configuration updates are automatically applied at that exact time.

---

### Edge Cases

- What happens when a scheduled change conflicts with a manual change made before the scheduled time? (e.g., I schedule it to turn ON tomorrow, but someone turns it ON manually today).
- How does the system handle scheduled timezones? (e.g., scheduling for 9:00 AM PST vs EST).
- What happens if the backend scheduler process restarts right when a schedule was supposed to trigger? (Needs to be robust and catch up).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow users with `RELEASE_MANAGER` or `ADMIN` roles to create a scheduled flag state change (ON/OFF).
- **FR-002**: System MUST store the scheduled time in UTC and allow the UI to convert to the user's local timezone.
- **FR-003**: System MUST provide a robust background scheduler that processes pending scheduled changes within 1 minute of their target time.
- **FR-004**: System MUST create an Audit Log entry when a scheduled change is executed, indicating it was triggered automatically by the scheduler.
- **FR-005**: System MUST allow users to cancel or modify a scheduled change before it executes.
- **FR-006**: System MUST reject requests to create a scheduled change for a flag if another pending schedule already exists for that same flag.

### Key Entities *(include if feature involves data)*

- **ScheduledChange**: Represents a pending action on a flag or change request at a specific timestamp. Attributes: TargetType (Flag/ChangeRequest), TargetID, Action (TurnOn/TurnOff/Apply), ScheduledFor (Timestamp), Status (Pending/Executed/Cancelled).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Scheduled changes are executed within 60 seconds of their target timestamp.
- **SC-002**: Users can successfully schedule, view, and cancel a flag state change via the Dashboard UI.
- **SC-003**: Background scheduling mechanism scales to support at least 1,000 concurrent scheduled events triggering at the exact same minute without degrading core API performance.

## Assumptions

- We will implement a polling-based background worker or cron-like mechanism in the Go backend rather than relying on external scheduling services like Cloud Tasks, to keep the architecture self-contained.
- Timezone handling will be done by sending UTC timestamps from the frontend to the backend.
- Modifying a flag's state manually while a schedule is pending does not automatically cancel the schedule (it will still execute and potentially overwrite the manual change).
