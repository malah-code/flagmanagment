# Feature Specification: Edge Proxy / Relay Node

**Feature Branch**: `[019-edge-proxy-relay]`

**Created**: 2026-08-08

**Status**: Draft

**Input**: User description: "what's next? taking into consideration the initial requirements from g-requirements.md and p-requirements.md"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Enterprise Private Network Isolation (Priority: P1)

As a Security Architect at a financial institution, I want to deploy a lightweight Edge Proxy inside my private corporate subnet so that none of my internal microservices ever need outbound internet access to the FlagManagment cloud backend.

**Why this priority**: High-security environments (banking, government, regulated industries) have strict outbound firewall policies. Without a relay node, every single microservice would need an exception opened to reach the FlagManagment backend. This blocks enterprise adoption entirely.

**Independent Test**: Deploy the Edge Proxy in an isolated network namespace. Configure a Go SDK to point at the Edge Proxy (not the backend). Toggle a flag in the FlagManagment dashboard. Verify the SDK receives the update without any SDK process having a direct route to the backend.

**Acceptance Scenarios**:

1. **Given** the Edge Proxy is running with a valid environment token and connected to the FlagManagment backend, **When** a microservice SDK connects to the Edge Proxy and evaluates a flag, **Then** the evaluation returns the correct result with sub-millisecond latency.
2. **Given** a flag is toggled in the FlagManagment dashboard, **When** the backend pushes a delta update, **Then** the Edge Proxy propagates that delta to all connected SDK clients within 500ms.
3. **Given** the Edge Proxy loses its connection to the FlagManagment backend, **When** a microservice SDK evaluates a flag, **Then** the evaluation is served from the Edge Proxy's last known good in-memory state (no service disruption).

---

### User Story 2 - Edge Proxy Operational Observability (Priority: P2)

As a DevOps/SRE Engineer, I want to monitor the health and connectivity status of the Edge Proxy so that I can detect and alert on connection loss between the proxy and the FlagManagment backend.

**Why this priority**: Enterprises need confidence that their local evaluation hub is operating correctly before they can trust it for production workloads.

**Independent Test**: Can be tested independently by hitting the proxy's health endpoint and asserting the `connected`, `upstream_latency_ms`, and `connected_since` fields are populated correctly.

**Acceptance Scenarios**:

1. **Given** the Edge Proxy is running and connected to the backend, **When** a health check request is made, **Then** the response reports `status: healthy` and `upstream_connected: true`.
2. **Given** the Edge Proxy has lost its upstream connection, **When** a health check request is made, **Then** the response reports `status: degraded` and `upstream_connected: false` so monitoring tools can fire alerts.
3. **Given** a configurable `HEALTH_CHECK_PORT` environment variable is set, **When** the Edge Proxy starts, **Then** the health endpoint is available on the configured port.

---

### User Story 3 - Multi-SDK Client Support (Priority: P2)

As a Backend Engineering Lead, I want the Edge Proxy to support multiple simultaneous internal microservice SDK connections so that hundreds of internal services can evaluate flags locally without individual upstream connections.

**Why this priority**: The entire value proposition of the relay node is fan-out — one upstream connection serving many internal consumers.

**Independent Test**: Connect 100 SDK clients concurrently to a single Edge Proxy instance and verify all clients receive flag updates within 1 second of a backend delta push.

**Acceptance Scenarios**:

1. **Given** 500 SDK clients are connected to the Edge Proxy, **When** a flag delta is received from the backend, **Then** all 500 connected clients receive the update without any client being dropped.
2. **Given** a client disconnects from the Edge Proxy, **When** that client reconnects, **Then** it receives the full current ruleset snapshot from the Edge Proxy's in-memory state without needing to contact the backend.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The Edge Proxy MUST be deployable as a standalone, stateless Docker container with configuration driven entirely by environment variables.
- **FR-002**: The Edge Proxy MUST maintain exactly one persistent upstream streaming connection to the FlagManagment backend using the environment SDK token.
- **FR-003**: The Edge Proxy MUST serve an in-memory flag ruleset snapshot to any connecting SDK client upon authentication.
- **FR-004**: When the backend pushes a delta update, the Edge Proxy MUST propagate that delta to all connected SDK clients within 500ms under normal load.
- **FR-005**: If the upstream connection to the FlagManagment backend is lost, the Edge Proxy MUST continue serving evaluation requests from its last known good in-memory state.
- **FR-006**: The Edge Proxy MUST implement exponential backoff reconnection to the backend on connection loss.
- **FR-007**: The Edge Proxy MUST expose a health check HTTP endpoint reporting upstream connectivity status and connected client count.
- **FR-008**: The Edge Proxy MUST support at least 500 concurrent downstream SDK client connections per instance.
- **FR-009**: All communication between the Edge Proxy and the FlagManagment backend MUST be secured with TLS.
- **FR-010**: Communication between internal SDK clients and the Edge Proxy MUST also support TLS (configurable).

### Key Entities

- **EdgeProxy**: The service instance. Holds a single upstream backend connection and N downstream SDK connections.
- **InMemoryRuleset**: The cached flag snapshot maintained in the proxy, updated on delta receipt.
- **DownstreamConnection**: A connected internal SDK client session.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A single Edge Proxy instance sustains 500 concurrent SDK client connections with zero dropped clients during a 10-minute stress test.
- **SC-002**: Flag delta propagation from backend receipt to all connected clients completes within 500ms under a load of 500 clients.
- **SC-003**: During a simulated backend outage, 100% of connected SDK clients continue to receive correct flag evaluations from the proxy's in-memory state with zero evaluation errors.
- **SC-004**: The Docker image for the Edge Proxy is less than 30MB compressed to support Raspberry Pi 4 deployment.

## Assumptions

- **Single Environment Token**: Each Edge Proxy instance is configured with exactly one environment token and serves only that environment's flag state. Multiple environments require multiple proxy instances.
- **Stateless Between Restarts**: The proxy does not persist its ruleset to disk. Upon restart, it reconnects to the backend and re-bootstraps from a fresh snapshot.
- **SDK Compatibility**: Internal SDK clients that previously connected directly to the FlagManagment backend can connect to the Edge Proxy without any SDK code changes — only a configuration endpoint change is required.
