# Feature Specification: SDK Evaluation API

**Feature Branch**: `[005-sdk-evaluation-api]`

**Created**: 2026-07-21

**Status**: Draft

**Input**: User description: "Feature 005: High-performance SDK Client API flag evaluation endpoints using Redis/In-memory caching"

## Clarifications
### Session 2026-07-22
- Q: Streaming Protocol: The spec assumes Server-Sent Events (SSE), but the project Constitution mandates gRPC/Protobuf for SDK streaming. Should we update the spec to reflect gRPC to comply with the Constitution? → A: Use gRPC/Protobuf. It is faster and more efficient for server-side SDKs, aligning with the Constitution's performance goals.
- Q: Redis Failure Fallback: The spec lists "What happens when Redis cache is unavailable?" as an edge case. Should the SDK endpoints fallback to querying PostgreSQL, or fail fast to protect the database? → A: Fail fast (return 503 Service Unavailable).
- Q: Thundering Herd Mitigation: How should we handle thousands of SDKs reconnecting simultaneously after a disconnect? → A: Require SDKs to implement jitter/exponential backoff.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - SDK Bootstrapping (Priority: P1)

As a server-side SDK, I want to download a complete snapshot of all feature flags and their targeting rules for my environment upon initialization so that I can evaluate flags locally in-memory with zero network latency.

**Why this priority**: Local in-memory evaluation is mandated by the Constitution for performance. The initial snapshot payload is required for any evaluation to occur.

**Independent Test**: Can be fully tested by sending an API request with a valid environment SDK token and receiving a JSON payload containing all active flags and their targeting rules.

**Acceptance Scenarios**:

1. **Given** a valid environment SDK token, **When** the SDK requests the `/api/v1/sdk/rules` endpoint, **Then** the API returns a 200 OK with the full flag ruleset.
2. **Given** an invalid or revoked environment SDK token, **When** the SDK requests the endpoint, **Then** the API returns a 401 Unauthorized.

---

### User Story 2 - Real-time Delta Updates (Priority: P2)

As a server-side SDK, I want to establish a persistent streaming connection to receive lightweight delta updates when flags change, so that my local in-memory cache stays accurate without polling.

**Why this priority**: While polling could work initially, persistent streaming is required to keep evaluations fast and accurate without overwhelming the server, aligning with the project's performance constraints.

**Independent Test**: Can be tested by establishing an SSE or gRPC stream and verifying that a flag change in the admin dashboard triggers an immediate push event over the stream.

**Acceptance Scenarios**:

1. **Given** an established streaming connection, **When** a flag's state or rule is modified in the system, **Then** the API pushes a delta payload with the specific change to the connected SDKs for that environment.
2. **Given** a dropped connection, **When** the SDK reconnects with the last known version hash, **Then** the API provides any missed updates or forces a full ruleset refresh.

---

### User Story 3 - Server-side Evaluation (Thin Clients) (Priority: P3)

As a thin client (e.g., mobile app or single-page application), I want to call an endpoint to evaluate specific flags for my user context so that I don't have to download the entire ruleset and expose all targeting logic to the public internet.

**Why this priority**: While server-side SDKs (thick clients) must evaluate locally, frontend/mobile SDKs (thin clients) shouldn't download full rulesets for security and payload size reasons.

**Independent Test**: Can be tested by sending an evaluation request with a user context and verifying the correct flag value is returned.

**Acceptance Scenarios**:

1. **Given** a valid thin-client SDK token and a user context, **When** a request is made to `/api/v1/sdk/evaluate`, **Then** the API evaluates the rules server-side and returns the boolean/multivariate result.

### Edge Cases

- **Redis Cache Unavailable**: The endpoints will fail fast and return a 503 Service Unavailable error to protect PostgreSQL from traffic spikes (thundering herd).
- **Thundering Herd Mitigation**: To handle thousands of SDKs reconnecting simultaneously (e.g., during a deployment), the SDKs MUST implement jitter and exponential backoff. The server will not implement queueing and will drop excess connection requests if limits are reached.
- What happens if the payload size for the ruleset exceeds typical HTTP response size constraints?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a high-performance HTTP endpoint (`/api/v1/sdk/rules`) to serve the complete ruleset snapshot for an environment.
- **FR-002**: System MUST authenticate all SDK API requests using the environment-specific SDK token.
- **FR-003**: System MUST provide a streaming mechanism (e.g., Server-Sent Events or gRPC) for pushing real-time updates to connected SDKs.
- **FR-004**: System MUST leverage Redis (or in-memory caching) to serve the ruleset endpoint and evaluation endpoints to meet sub-millisecond response targets.
- **FR-005**: System MUST invalidate or update the relevant Redis cache entries immediately when a flag state or targeting rule is changed via the management API.
- **FR-006**: System MUST provide a server-side evaluation endpoint (`/api/v1/sdk/evaluate`) for thin clients that accepts a user context and returns the calculated flag variation.
- **FR-007**: System MUST salt and hash any PII provided in the evaluation context before processing or logging, per Constitution VII.

### Key Entities *(include if feature involves data)*

- **Ruleset Snapshot**: A JSON representation of all flags, their default states, and targeting rules for a specific environment.
- **SDK Token**: A cryptographically secure token tied to a specific environment, used to authenticate the SDK client.
- **Evaluation Context**: A dictionary of attributes (e.g., user ID, email, country, custom properties) used to evaluate targeting rules.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The `/api/v1/sdk/rules` snapshot endpoint responds in under 10ms at the 95th percentile under normal load.
- **SC-002**: Flag changes are propagated to connected SDKs via the streaming connection in under 500ms at the 99th percentile.
- **SC-003**: The server-side `/api/v1/sdk/evaluate` endpoint responds in under 5ms at the 99th percentile when rules are cached.
- **SC-004**: The system supports at least 10,000 concurrent streaming connections per server instance.

## Assumptions

- Redis 7+ is deployed and accessible to the backend application, as mandated by the Constitution.
- gRPC will be the primary mechanism for streaming updates to server-side SDK clients, replacing initial assumptions of SSE, to ensure maximum performance and type safety.
- SDK clients are responsible for implementing the local caching and MurmurHash3 bucketing logic as defined in the Constitution; this feature only covers the server-side API endpoints providing the data.
