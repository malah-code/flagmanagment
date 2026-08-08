# Feature Specification: PII Hashing & Data Privacy

**Feature Branch**: `021-pii-hashing-privacy`

**Created**: 2026-08-08

**Status**: Draft

**Input**: User description: "PII Hashing & Data Privacy"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Secure Identity Storage (Priority: P1)

As a security-conscious organization, I want all user identity attributes used for flag targeting to be securely hashed and salted before being stored in the database, so that if the database is compromised, user privacy is maintained.

**Why this priority**: Protecting PII is a core constitutional requirement and essential for enterprise compliance (SOC 2, GDPR).

**Independent Test**: Can be fully tested by verifying that raw identity values (like emails or user IDs) passed during SDK evaluation or stored in persistent analytics are never written in plaintext to the database.

**Acceptance Scenarios**:

1. **Given** a user is evaluated against a targeting rule using their email address, **When** the evaluation record is stored for analytics, **Then** the stored identifier must be a salted hash, not the plaintext email.
2. **Given** targeting rules are configured in the dashboard, **When** a user enters an explicit user ID to target, **Then** the system must hash the value with the project's salt before persisting it.

---

### User Story 2 - Consistent Identity Bucketing (Priority: P1)

As a developer using the SDK, I want the hashing mechanism to be deterministic across all SDKs and the backend, so that a specific user consistently falls into the same rollout bucket regardless of which SDK evaluates them.

**Why this priority**: Feature flag stickiness is crucial; users should not bounce between variations simply because they hit different services or platforms.

**Independent Test**: Can be fully tested by generating hashes for identical identities across different SDK implementations and verifying they produce the exact same bucket assignment.

**Acceptance Scenarios**:

1. **Given** a user identity "user@example.com", **When** evaluated by the backend and a local SDK, **Then** both must produce the identical hashed value for bucketing.

---

### User Story 3 - Configurable Data Retention (Priority: P2)

As a compliance officer, I want to configure how long hashed identity data and evaluation analytics are retained, so that we comply with our internal data minimization policies.

**Why this priority**: Hashed PII is still considered sensitive in some contexts. Organizations need control over data lifecycles.

**Independent Test**: Can be fully tested by configuring a 30-day retention policy and verifying that records older than 30 days are automatically purged from the system.

**Acceptance Scenarios**:

1. **Given** a data retention policy is set to 30 days, **When** the scheduled cleanup job runs, **Then** all analytics and telemetry data older than 30 days must be permanently deleted.

### Edge Cases

- What happens if the environment salt is rotated? (Existing targeted hashes may invalidate; needs a clear mitigation or warning).
- How does the system handle extremely long or malformed identity strings?
- What happens if an organization requires zero data retention for analytics?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST support deterministic salting and hashing of all user identity attributes (PII) used in targeting rules or analytics.
- **FR-002**: The system MUST automatically generate a unique, cryptographically secure salt for each Environment upon creation.
- **FR-003**: The system MUST never store plaintext identity values used for evaluation analytics in the database.
- **FR-004**: SDKs MUST receive the environment salt securely to perform consistent local deterministic bucketing without transmitting plaintext PII.
- **FR-005**: The system MUST provide a configurable data retention policy for evaluation analytics and telemetry data, with a default of 30 days.

### Key Entities

- **Environment**: Must include a `Salt` attribute used for hashing identities within that environment.
- **EvaluationAnalytics**: Records of flag evaluations. Must only contain hashed identifiers, never plaintext PII.
- **TargetingRule**: When storing specific user IDs to include/exclude, these must be hashed before storage.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of identity attributes in targeting rules and evaluation analytics are stored as salted hashes.
- **SC-002**: Zero plaintext PII is written to persistent storage by the core evaluation and analytics engine.
- **SC-003**: Hashing algorithms must be highly performant, ensuring server-side evaluation still executes in under 1 millisecond.
- **SC-004**: Hashing is deterministic across 100% of supported SDK languages for consistent bucketing.

## Assumptions

- We assume the MurmurHash3 algorithm (as mandated by the Constitution) is used for deterministic cross-language bucketing (prioritizing <1ms performance), while SHA-256 is used for hashing PII for persistent storage (prioritizing cryptographic security for audit logs and analytics).
- We assume environment-level salting provides sufficient isolation and security.
- We assume the retention cleanup job can run asynchronously during off-peak hours to avoid impacting evaluation performance.
