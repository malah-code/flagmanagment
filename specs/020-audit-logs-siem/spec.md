# Feature Specification: Immutable Audit Logs & SIEM Webhooks

**Feature Branch**: `[020-audit-logs-siem]`

**Created**: 2026-08-08

**Status**: Draft

**Input**: User description: "Immutable Audit Logs & SIEM Webhooks"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Compliance Auditing and Traceability (Priority: P1)

As a Compliance Officer, I want the system to record an immutable, append-only ledger of all administrative actions (like flag creations, rule modifications, role assignments, and environment changes) so that I can satisfy SOC 2 and internal auditing requirements.

**Why this priority**: Enterprise customers cannot adopt a feature flag platform that mutates production state without a strict, unalterable trail of who changed what, when, and from where.

**Independent Test**: Perform a series of CRUD operations on a feature flag via the API. Query the audit logs endpoint and assert that a distinct log entry exists for each action containing the correct Actor ID, Timestamp, Action Type, Previous State, and New State.

**Acceptance Scenarios**:

1. **Given** an authenticated user creates a new feature flag, **When** the transaction commits, **Then** a corresponding audit log entry is synchronously recorded in the database.
2. **Given** a user modifies a targeting rule, **When** viewing the audit log, **Then** the log entry clearly displays the previous JSON state and the new JSON state for diffing.
3. **Given** an actor attempts to delete or modify an existing audit log entry via the application API, **Then** the request is rejected (append-only constraint).

---

### User Story 2 - Security Log Sanitization (Priority: P1)

As a Security Engineer, I want the audit logging system to actively sanitize log payloads so that plaintext API keys, PII (like email addresses or identity hashes), and sensitive configuration data are never stored in the audit ledger.

**Why this priority**: Storing plaintext secrets or PII in an audit log violates strict compliance requirements and creates a massive security vulnerability.

**Independent Test**: Generate a new environment API key. Query the audit log for the environment creation/update event. Assert that the API key is completely redacted or hashed in the logged JSON state.

**Acceptance Scenarios**:

1. **Given** an admin generates a new environment API key, **When** the audit log is recorded, **Then** the API key field in the `NewState` JSON payload is replaced with `[REDACTED]`.
2. **Given** a user modifies targeting rules that include PII-based identity rules, **When** logged, **Then** sensitive identity fields remain appropriately hashed/redacted.

---

### User Story 3 - SIEM Webhook Streaming (Priority: P2)

As a SecOps Engineer, I want the platform to stream audit log events in real-time via webhooks to my external SIEM tools (e.g., Splunk, Datadog) so that I can set up organization-wide alerts on unauthorized changes or suspicious behavior.

**Why this priority**: Security teams monitor infrastructure centrally through SIEMs; they do not want to log into the FlagManagment dashboard to look for security events.

**Independent Test**: Register a mock webhook endpoint in the project settings. Perform a flag change. Assert that the mock endpoint receives an HTTP POST request containing the sanitized audit log event payload within 5 seconds.

**Acceptance Scenarios**:

1. **Given** a webhook integration is configured for "Audit Events", **When** an administrative action occurs, **Then** a JSON payload of the audit event is HTTP POSTed to the configured webhook URL.
2. **Given** the remote SIEM webhook endpoint is down or returns a 500 error, **Then** the system implements a retry mechanism with exponential backoff to ensure event delivery.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST record an append-only ledger entry for the following actions: Flag creation/modification/deletion, Environment creation/protection state changes, Role assignments/invitations, and Change Request lifecycle transitions.
- **FR-002**: Each audit log entry MUST include: a UUID, Timestamp, Actor User ID, Target IDs (Project/Environment/Flag), Action Type, Previous JSON State, New JSON State, and the Actor's IP Address.
- **FR-003**: The backend MUST execute a sanitization pass on the Previous/New JSON states before saving to the database to scrub known sensitive fields (e.g., `api_key`, `token`, password fields).
- **FR-004**: The API MUST expose paginated, filterable endpoints to query audit logs by Project, Environment, Actor, and Date Range.
- **FR-005**: The API MUST expose an endpoint to export a filtered set of audit logs in CSV format.
- **FR-006**: The system MUST support registering Webhook URLs at the Project level.
- **FR-007**: When a webhook is configured for audit events, the backend MUST asynchronously HTTP POST the sanitized audit log entry to the webhook URL.
- **FR-008**: The webhook dispatcher MUST implement a retry policy for failed deliveries (e.g., up to 3 retries with backoff).

### Key Entities

- **AuditLog**: The core ledger record containing the actor, action, timestamp, and state diffs.
- **WebhookIntegration**: Configuration tying a Project to an external SIEM endpoint URL.
- **WebhookDelivery**: (Optional/Internal) Tracking the success/failure state of an outbound webhook dispatch for retry logic.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of state-mutating administrative API calls result in a corresponding, sanitized audit log entry.
- **SC-002**: Webhook dispatch to external SIEMs triggers within 2 seconds of the actual transaction committing to the database.
- **SC-003**: The audit log API endpoint can return a paginated response of 50 records in under 200ms when querying a table with 1,000,000 log entries.
- **SC-004**: Automated security scanners confirm zero plaintext API keys or known PII leak into the `audit_logs` table during regular operational load testing.

## Assumptions

- **Read Operations**: Pure read operations (e.g., viewing a flag, fetching the ruleset via SDK) are NOT recorded in the Audit Log to prevent massive database bloat. Only administrative mutations are logged.
- **Storage Strategy**: The PostgreSQL `audit_logs` table will grow indefinitely. We assume organizations will manage data retention (e.g., archiving logs older than 1 year) outside the scope of this immediate feature, or through a future automated cleanup job.
- **Webhook Format**: The webhook payload will be a standard JSON representation of the `AuditLog` entity, rather than implementing specific proprietary formats for Splunk/Datadog natively. Customers can parse the JSON on their SIEM side.
