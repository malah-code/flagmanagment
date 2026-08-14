# Feature Specification: Server-Side Environment Keys & Tabbed Keys Management

**Feature Branch**: `[###-feature-name]`

**Created**: 2026-08-10

**Status**: Draft

**Input**: User description: "specify the enhancement in detail before we start tasks and implementation: Dedicated environment keys management view matching Flagsmith architecture. Includes a dedicated 'Keys' settings tab per environment, separate Client-side Environment Key section with 1-click copy, and a full Server-side Environment Keys management section allowing admins to create named server keys (for local evaluation / backend SDKs), show/hide masked key values, search server keys, and delete/revoke keys."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Dedicated "Keys" Settings Tab (Priority: P1)

As an environment administrator, I want to access a dedicated "Keys" tab under Environment Settings, so that I can manage both client-side and server-side credentials in one clear, structured view.

**Why this priority**: Directly aligns with Flagsmith's tabbed layout (`General`, `SDK Settings`, `Keys`, `Members`, `Webhooks`), organizing environment management seamlessly.

**Independent Test**: Can be verified by navigating to Environment Settings and clicking the "Keys" tab.

**Acceptance Scenarios**:

1. **Given** the user navigates to Environment Settings, **When** they click the "Keys" tab, **Then** the view renders distinct sections for "Client-side Environment Key" and "Server-side Environment Keys".

### User Story 2 - Server-Side Environment Keys Management (Priority: P1)

As a backend developer or admin, I want to create named Server-Side Environment Keys (e.g. `billing-service-key`, `staging-k8s-pod`) for local evaluation SDKs, so that backend services can securely evaluate flags in-memory.

**Why this priority**: Server-side SDKs require privileged local evaluation keys. Allowing administrators to issue named server keys enables granular key rotation and service tracking.

**Independent Test**: Can be verified by clicking "Create Server-side Environment Key", entering a key name, and observing the new server key appear in the Server-side Keys table.

**Acceptance Scenarios**:

1. **Given** the user is on the "Keys" tab, **When** they click "Create Server-side Environment Key", **Then** a modal prompts for a key name (e.g. `payment-backend`).
2. **Given** the key is created, **When** it renders in the table, **Then** its secret value is masked by default (`••••••••••••••••••••`) with a "Show / Hide" toggle button and a "Copy" button.
3. **Given** the user clicks the trash icon on a server key row, **When** confirmed, **Then** the key is immediately revoked and deleted.

### User Story 3 - Search and Filter Server-Side Keys (Priority: P2)

As a system administrator with dozens of microservices, I want to search and filter server-side keys by name, so that I can quickly inspect or revoke keys for a specific service.

**Why this priority**: Maintains usability as the number of server-side integrations grows.

**Independent Test**: Can be verified by entering a name into the search bar inside the "Server-Side Environment Keys" panel and filtering the table rows.

**Acceptance Scenarios**:

1. **Given** multiple server-side keys exist, **When** the user types in the search bar, **Then** the table filters rows to only match key names containing the query.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a dedicated "Keys" tab in Environment Settings.
- **FR-002**: System MUST render an explicit "Client-side Environment Key" card with an input field and a 1-click "Copy" button.
- **FR-003**: System MUST provide a "Server-side Environment Keys" panel allowing admins to create, list, search, reveal/mask, copy, and delete server keys.
- **FR-004**: System MUST mask server key values by default with a "Show/Hide" toggle button for security.
- **FR-005**: System MUST allow naming each server key upon creation to track which microservice or server environment uses it.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% parity with Flagsmith's Keys management UX layout and functionality.
- **SC-002**: Admins can issue a new named server-side key in under 15 seconds.
- **SC-003**: 0% exposure of plain server key secrets unless explicitly toggled by an authorized admin.

## Assumptions

- Server-side keys are used for high-performance in-memory local evaluation in server SDKs (Go, Node, Python, Java).
- Client-side keys remain public and read-only for browser/mobile SDKs.
