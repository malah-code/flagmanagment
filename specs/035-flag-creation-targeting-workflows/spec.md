# Feature Specification: Feature Flag Creation & Targeting Workflows

**Feature Branch**: `035-flag-creation-targeting-workflows`

**Created**: 2026-08-15

**Status**: Draft

**Input**: User description: "proceed to testing the Feature Flag Creation & Targeting workflows with Puppeteer"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create & Configure Multi-Type Flags (Priority: P1)

Administrators and Project Editors can create feature flags of different types (`BOOLEAN`, `MULTIVARIATE`, and `JSON`) with customizable default variations, descriptions, tags, and lifecycle states directly within their workspace.

**Why this priority**: Flag creation is the fundamental prerequisite for feature flag lifecycle management, deployment gating, and SDK evaluations.

**Independent Test**: Create a Boolean, Multivariate, and JSON flag via the "New Feature Flag" modal, and verify they appear immediately in the Flag Definitions and Environment Targeting tables.

**Acceptance Scenarios**:

1. **Given** an authenticated user on the project flags overview page, **When** they click "New Feature Flag", choose type `BOOLEAN`, fill in a unique key `enable-dark-mode` and name `Enable Dark Mode`, and submit, **Then** the flag is created with Default ON/OFF variations and appears in the list.
2. **Given** a user creating a `MULTIVARIATE` flag with key `button-color-experiment`, **When** they define variations (e.g., `blue`, `green`, `purple`) with custom values, **Then** the flag persists all variations and allows weight distribution across them.
3. **Given** a user creating a `JSON` remote config flag, **When** they provide a valid JSON default payload, **Then** the flag validates the schema and stores the remote config payload.

---

### User Story 2 - Contextual Attribute Targeting Rules (Priority: P1)

Feature owners can configure user segmentation and contextual targeting rules within specific environments so that flags evaluate conditionally based on user attributes (such as country, subscription tier, beta tester status, or email domain).

**Why this priority**: Targeted rollouts and canary releases are essential for safe production deployments and beta testing.

**Independent Test**: Open the Targeting modal for an active flag, add a rule `country EQUALS 'US'`, set rollout variation to ON, and verify the rule is persisted and returned in the environment flag evaluation state.

**Acceptance Scenarios**:

1. **Given** an active feature flag in a target environment, **When** an editor opens the Targeting configuration and adds a rule targeting `tier EQUALS 'enterprise'`, **Then** the rule is saved to the environment flag state.
2. **Given** multiple targeting rules, **When** the editor reorders rules, **Then** the system evaluates rules in strict priority order (first match wins).
3. **Given** percentage rollout rules, **When** the user assigns 20% to variation A and 80% to variation B, **Then** the system validates that total weights sum to 100%.

---

### User Story 3 - Emergency Kill Switches & Instant Deactivation (Priority: P2)

Operators can trigger an emergency kill switch on any flag in an environment to instantly disable it and revert evaluations to a safe fallback without requiring code redeployments or full change requests.

**Why this priority**: Fast mitigation of critical production regressions is vital for site reliability and incident management.

**Independent Test**: Click "Kill Switch" on an active flag, confirm the action, and verify the flag status immediately updates to Inactive with an incident reason logged.

**Acceptance Scenarios**:

1. **Given** an active flag causing production errors, **When** an administrator clicks "Kill Switch" and provides a reason `High latency in checkout service`, **Then** the flag is immediately forced OFF in that environment and an audit event is recorded.

---

### User Story 4 - Scheduled Flag Transitions (Priority: P3)

Teams can schedule future state changes for a flag (such as turning ON a promotional banner at a specific date and time) that will be automatically enacted by the backend scheduler.

**Why this priority**: Enables timed campaign launches and off-hours release automation without manual intervention.

**Independent Test**: Schedule a flag state change for a future timestamp, and verify the scheduled change is listed in the pending schedule queue.

**Acceptance Scenarios**:

1. **Given** an inactive flag in an environment, **When** a user clicks "Schedule", selects a future timestamp, and sets the target state to `ACTIVE`, **Then** a scheduled change entry is created and shown with its scheduled execution time.

---

### Edge Cases

- What happens when a user creates a flag with a key that already exists in the project? The system rejects the creation with a descriptive "Flag key already exists" error.
- What happens when invalid JSON is submitted for a JSON-type flag? The UI and backend validate JSON syntax and reject malformed JSON with specific parse error line indicators.
- What happens when targeting rule attributes contain whitespace or special characters? Attributes are trimmed and properly escaped.
- What happens when a kill switch is triggered on an already disabled flag? The action succeeds gracefully and logs the event state.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow creating feature flags of type `BOOLEAN`, `MULTIVARIATE`, and `JSON` with unique keys within a project.
- **FR-002**: System MUST automatically initialize environment flag states across all existing environments when a new flag is created.
- **FR-003**: System MUST allow configuring contextual targeting rules based on attribute operators (`EQUALS`, `NOT_EQUALS`, `CONTAINS`, `GREATER_THAN`, `LESS_THAN`, `IN_LIST`).
- **FR-004**: System MUST support percentage-based rollouts with deterministic hashing for consistent user experience.
- **FR-005**: System MUST provide an emergency Kill Switch mechanism that immediately disables flag evaluations in the selected environment.
- **FR-006**: System MUST record all flag creations, updates, targeting changes, and kill switch activations in the immutable audit log.
- **FR-007**: System MUST support scheduling flag state transitions for automated future execution.

### Key Entities

- **FeatureFlag**: Core flag definition belonging to a project, having key, name, description, type, and variations.
- **Variation**: Represents a distinct evaluation output (value, name, description) for boolean or multivariate flags.
- **EnvironmentFlagState**: Environment-specific state containing enabled status, targeting rules, rollout percentages, and override values.
- **TargetingRule**: Conditional rule containing attribute name, operator, comparison values, and target variation.
- **ScheduledChange**: Queued state transition with target environment, scheduled execution time, and target state payload.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can create and configure a new feature flag from the dashboard in under 30 seconds.
- **SC-002**: Flag state toggles and targeting changes update the dashboard UI instantly (< 500ms).
- **SC-003**: 100% of flag mutations generate a corresponding structured entry in the workspace audit log.
- **SC-004**: All flag evaluation states persist accurately across browser reloads and server restarts.

## Assumptions

- Standard local development environment running Docker with PostgreSQL and Redis.
- Users executing administrative operations possess appropriate project or global roles (`ADMIN` or `EDITOR`).
- OpenFeature-compatible evaluation format is preserved for all flag rules and variations.
