# Feature Specification: Core API Service

**Feature Branch**: `003-core-api`

**Created**: 2026-07-20

**Status**: Draft

**Input**: User description: "Design the REST API endpoints for managing projects, environments, and feature flags. Include endpoints for the SDK to fetch the flag state for a specific environment using the environment's api key hash."

## User Scenarios & Testing *(mandatory)*

## Clarifications
### Session 2026-07-20
- Q: Should we implement HTTP caching strategies (like ETag) for the SDK endpoint? → A: Yes, implement ETag / `If-None-Match` support for HTTP 304 Not Modified responses.
- Q: Which pagination strategy should we use for list endpoints? → A: Follow Google API standards using cursor-based `pageToken` and `pageSize`.

### User Story 1 - Workspace & Hierarchy Management (Priority: P1)

As a platform administrator or developer, I want to manage projects and environments via REST API so that I can organize my feature flags across different deployment targets.

**Why this priority**: Without projects and environments, feature flags cannot be scoped or evaluated. This is the foundational management API.

**Independent Test**: Can be fully tested by sending HTTP POST/GET requests to create a project, then create environments within that project, and verifying the JSON responses and HTTP status codes.

**Acceptance Scenarios**:

1. **Given** an authorized user, **When** they submit a valid POST request to `/api/v1/projects`, **Then** the project is created and a 201 Created response is returned with the project ID.
2. **Given** an existing project, **When** a user submits a POST request to `/api/v1/projects/{projectId}/environments`, **Then** the environment is created, an API key is generated and hashed, and a 201 response is returned.

---

### User Story 2 - Feature Flag Configuration (Priority: P1)

As a developer or product manager, I want to create, update, and toggle feature flags via the REST API so that I can control feature releases without code deployments.

**Why this priority**: The core value proposition of the platform is the ability to manage feature flags dynamically.

**Independent Test**: Can be fully tested by creating a flag definition and then updating its targeting rules/state for a specific environment via the API.

**Acceptance Scenarios**:

1. **Given** an existing project, **When** a user submits a POST request to create a flag at `/api/v1/projects/{projectId}/flags`, **Then** the flag definition is created globally for the project.
2. **Given** an existing environment and flag, **When** a user submits a PUT request to `/api/v1/environments/{envId}/flags/{flagId}/state` with new targeting rules, **Then** the environment-specific flag state is updated.

---

### User Story 3 - SDK Evaluation Retrieval (Priority: P1)

As an application integrating the SDK, I want to fetch all active feature flags and their states for my specific environment using my API key, so that my application can evaluate flags locally.

**Why this priority**: The SDK needs to fetch rules from the server to evaluate flags. This is the highest volume and most critical read path in the system.

**Independent Test**: Can be fully tested by sending a GET request with an `Authorization: Bearer {environment_api_key}` header and verifying the correct evaluation payload is returned rapidly.

**Acceptance Scenarios**:

1. **Given** a valid environment API key, **When** the SDK calls the evaluation endpoint `/api/v1/evaluate/flags`, **Then** the system returns a 200 OK with a complete JSON payload of all flags and rules for that environment.
2. **Given** an invalid or revoked API key, **When** the SDK calls the evaluation endpoint, **Then** the system rejects the request with a 401 Unauthorized error.

---

### Edge Cases

- What happens when a user requests a flag or project that doesn't exist? The API MUST return a standard 404 Not Found response.
- How does the system handle an SDK evaluation request with malformed headers? The API MUST return a 400 Bad Request or 401 Unauthorized with a clear error payload.
- What happens if the JSON payload for targeting rules is structurally invalid during a state update? The API MUST return a 400 Bad Request detailing the schema violation.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide RESTful endpoints for CRUD operations on Projects (`/api/v1/projects`).
- **FR-002**: System MUST provide RESTful endpoints for CRUD operations on Environments (`/api/v1/projects/{projectId}/environments`).
- **FR-003**: System MUST provide RESTful endpoints for creating and managing Feature Flags globally per project (`/api/v1/projects/{projectId}/flags`).
- **FR-004**: System MUST provide RESTful endpoints for managing Environment-specific Flag States (`/api/v1/environments/{envId}/flags/{flagId}/state`).
- **FR-005**: System MUST provide a high-performance evaluation endpoint for SDKs (`/api/v1/evaluate/flags`) that authenticates using an environment API key hash and supports ETag / `If-None-Match` headers to return HTTP 304 Not Modified when flags haven't changed.
- **FR-006**: System MUST return standardized JSON error responses (RFC 7807 Problem Details or similar) for all 4xx and 5xx HTTP errors.
- **FR-007**: System MUST validate all incoming JSON payloads against defined schemas before attempting database operations.
- **FR-008**: System MUST implement pagination for all list endpoints following standard API patterns (e.g., Google APIs) using `pageToken` and `pageSize`.

### Key Entities

- **Project API Resource**: JSON representation of a project workspace.
- **Environment API Resource**: JSON representation of an environment, omitting sensitive hashes but returning the plaintext API key exactly once upon creation.
- **Feature Flag API Resource**: JSON definition of a flag.
- **SDK Evaluation Payload**: Highly optimized JSON document containing all targeting rules and states required by the SDK for local evaluation.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Management API endpoints respond to 95% of standard CRUD requests in under 200ms.
- **SC-002**: SDK evaluation endpoint `/api/v1/evaluate/flags` responds to 99% of valid requests in under 50ms, excluding network latency.
- **SC-003**: API achieves 100% compliance with OpenAPI 3.0 specification (no undocumented endpoints or fields).
- **SC-004**: SDK authentication successfully rejects 100% of invalid or malformed API keys.

## Assumptions

- Standard REST HTTP verbs (GET, POST, PUT, DELETE, PATCH) will be used appropriately.
- Authentication for Management APIs (projects, environments) will use standard Bearer tokens (JWT or session), which will be mocked or stubbed for this feature until a full auth provider is integrated.
- SDKs will fetch the entire state for an environment on startup (bulk fetch) rather than querying flags one-by-one.
- No rate limiting is required for the MVP API layer, though standard timeouts will be enforced.
