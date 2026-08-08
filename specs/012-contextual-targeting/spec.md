# Feature Specification: Contextual Targeting Engine

**Feature Branch**: `012-contextual-targeting`

**Created**: 2026-07-25

**Status**: Draft

**Input**: User description: "Contextual Targeting Engine (Equals, Contains, Regex matching)"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Rule Creation for Flags (Priority: P1)

As a feature flag administrator, I want to define targeting rules (Equals, Contains, Regex) for a feature flag within an environment so that I can control precisely which users receive specific flag variations based on their contextual attributes (e.g., email domain, tenant ID, region).

**Why this priority**: Core functionality of contextual targeting; without this, flags can only be globally enabled or disabled.

**Independent Test**: Can be tested by defining a rule in the dashboard and verifying it saves correctly into the environment's flag state.

**Acceptance Scenarios**:

1. **Given** I am editing a flag in a specific environment, **When** I add a targeting rule requiring "email" to "contain" "@company.com", **Then** the rule is persisted in the targeting rules configuration.
2. **Given** a targeting rule configuration, **When** I save it, **Then** the JSON structure reflects the operators and values correctly.

---

### User Story 2 - SDK Evaluation of Contextual Rules (Priority: P1)

As an application developer using the SDK, I want the SDK to evaluate the contextual targeting rules against the user context I provide locally, so that users seamlessly receive the correct flag variation without network latency.

**Why this priority**: Required for the rules to actually affect application behavior.

**Independent Test**: Can be tested by providing a user context to the SDK and asserting the correct boolean/variant result based on the defined rules.

**Acceptance Scenarios**:

1. **Given** a flag with a rule "tenant_id EQUALS '123'", **When** the SDK is queried with a context where `tenant_id = '123'`, **Then** the flag evaluates to TRUE (or the specified variation).
2. **Given** a flag with a rule "email REGEX '.*@beta\\.com'", **When** the SDK is queried with `email = 'user@beta.com'`, **Then** the flag evaluates to TRUE.
3. **Given** a flag with a rule, **When** the SDK context does not match the rule, **Then** the flag evaluates to the default fallback variation (e.g., FALSE).

---

### User Story 3 - Complex Multi-Condition Rules (Priority: P2)

As a feature flag administrator, I want to combine multiple conditions (AND/OR logic) within a single targeting rule, so that I can create sophisticated rollout strategies (e.g., "Internal Users" AND "US Region").

**Why this priority**: High-value for enterprise customers who need complex targeting segments.

**Independent Test**: Can be tested by saving a multi-condition rule and verifying the SDK evaluates it using logical AND/OR correctly.

**Acceptance Scenarios**:

1. **Given** a rule with two conditions joined by AND, **When** the context satisfies only one condition, **Then** the rule evaluates to FALSE.
2. **Given** a rule with two conditions joined by AND, **When** the context satisfies both conditions, **Then** the rule evaluates to TRUE.

### Edge Cases

- What happens when a context attribute used in a rule is missing from the provided SDK context? (Should evaluate to false for that rule).
- How does the system handle invalid Regex patterns? (Should fail validation at creation time in the dashboard).
- What if the context attribute type doesn't match the expected type (e.g., comparing a number to a string regex)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow users to define targeting rules on a per-environment flag state.
- **FR-002**: System MUST support the following operators: `EQUALS`, `CONTAINS`, and `REGEX`.
- **FR-003**: System MUST validate `REGEX` patterns upon rule creation/update to prevent evaluation errors in SDKs.
- **FR-004**: System MUST evaluate targeting rules in the Server-Side SDKs locally using the provided user context.
- **FR-005**: System MUST default to a fallback variation if no targeting rules match the provided context.
- **FR-006**: System MUST support combining multiple conditions using logical `AND`.

### Key Entities *(include if feature involves data)*

- **TargetingRule**: A logical structure defining conditions (Attribute, Operator, Value) and the resulting variation if the conditions are met. (Stored within `environment_flag_states` JSONB).
- **EvaluationContext**: A dynamic dictionary of key-value pairs representing the user/application context (e.g., `{"email": "...", "tenant": "..."}`).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Targeting rules can be evaluated by the SDK in under 1 millisecond.
- **SC-002**: 100% of invalid Regex rules are caught during creation and rejected by the API.
- **SC-003**: Complex targeting rule JSON payloads do not exceed reasonable limits (e.g., < 10KB per flag) to ensure fast gRPC streaming updates.

## Assumptions

- Users will provide a valid `EvaluationContext` map when querying the SDK.
- The UI will provide a builder interface for these rules (though UI specifics can be handled in a frontend-focused task).
- All targeting rules are stored inside the `targeting_rules` JSONB column of the `environment_flag_states` table.
- "OR" logic is achieved by creating multiple independent rules that evaluate sequentially (first match wins), while "AND" is used within a single rule.
