# Feature Specification: User Management API Implementation

**Feature Branch**: `[034-user-management-api]`

**Created**: 2026-08-14

**Status**: Draft

**Input**: User description: "we need to fix the mock, and make it real implementation"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Team Members (Priority: P1)

As an Administrator, I need to see a list of all team members, their global roles, and their project-specific access so that I can audit who has access to the system.

**Why this priority**: Visibility into system access is the foundation of RBAC and governance.

**Independent Test**: Can be fully tested by loading the Team Settings page and verifying that the displayed users match the database records.

**Acceptance Scenarios**:

1. **Given** I am logged in as an Administrator, **When** I navigate to the Team Settings page, **Then** I see a list of all users, their roles, and their assigned projects.
2. **Given** a user has a global role, **When** they are displayed in the list, **Then** their assigned projects show as "All Projects".

---

### User Story 2 - Invite New User (Priority: P1)

As an Administrator, I need to invite new colleagues by email and assign them initial roles/projects so that they can securely join the platform.

**Why this priority**: Without a way to invite users, the platform cannot be adopted by a team.

**Independent Test**: Can be fully tested by submitting an invitation and verifying that an email is dispatched and the user appears in the list as "Pending".

**Acceptance Scenarios**:

1. **Given** I am an Administrator, **When** I submit the invite form with an email and role, **Then** the system sends an invitation email.
2. **Given** I have sent an invite, **When** I view the user list, **Then** the invited user appears with a "Pending Invitation" status.

---

### User Story 3 - Edit User Roles and Access (Priority: P2)

As an Administrator, I need to modify a user's role or assigned projects so that I can adjust access as team responsibilities change.

**Why this priority**: Access needs change over time; an immutable access model is not viable for enterprises.

**Independent Test**: Can be fully tested by modifying a user's role and verifying the backend correctly applies and enforces the new permissions.

**Acceptance Scenarios**:

1. **Given** I am an Administrator, **When** I save edits to a user's role from "Project Editor" to "Global Administrator", **Then** the user immediately gains system-wide access.

---

### User Story 4 - Configure System Email Server (Priority: P2)

As a System Administrator, I need to configure the outbound SMTP server settings (host, port, credentials) within the platform so that the system can actually dispatch invitation and notification emails to users.

**Why this priority**: Without configuring a mail server, the invitation capability built in Story 2 cannot function in a self-hosted/enterprise environment.

**Independent Test**: Can be fully tested by saving SMTP credentials and successfully triggering a "Test Connection" email from the UI.

**Acceptance Scenarios**:

1. **Given** I am an Administrator, **When** I navigate to the System Configuration page and input SMTP details, **Then** the system securely persists these settings.
2. **Given** configured SMTP settings, **When** I click "Test Connection", **Then** the system dispatches a test email to my address and confirms success or displays the relevant connection error.

### Edge Cases

- What happens when an Administrator tries to invite an email address that already belongs to an existing user?
- How does the system handle an Administrator trying to demote their own role or remove their own global access?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide an API endpoint to retrieve a paginated list of all users, including their roles and project assignments.
- **FR-002**: System MUST provide an API endpoint to invite a new user via email address, assigning them an initial role and optional projects.
- **FR-003**: System MUST dispatch an email containing a secure registration link when an invitation is created.
- **FR-004**: System MUST provide an API endpoint to update an existing user's role and project assignments.
- **FR-005**: System MUST provide secure API endpoints to get and update system configuration parameters (specifically SMTP settings: host, port, username, secure password storage, sender email).
- **FR-006**: System MUST provide a "Test SMTP Connection" API endpoint to dispatch a diagnostic email using the current configuration.
- **FR-007**: Frontend MUST integrate with these API endpoints, completely replacing the hardcoded mock data in the `UsersManagement` view and adding a new `System Settings` view.
- **FR-008**: System MUST track and expose the "Last Active" status or "Pending" state for each user.

### Key Entities

- **User**: Represents an authenticated identity in the system.
- **Role Assignment**: The mapping between a User, a Role (Global Admin, Editor, Auditor), and optionally a specific Project.
- **Invitation**: A temporary, secure token-based record used to onboard a new user.
- **System Configuration**: A global singleton entity holding operational parameters (e.g., SMTP settings) for the platform instance.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Administrators can view the accurate list of users loaded from the backend API in under 500ms.
- **SC-002**: An invitation email is successfully dispatched to the target email address within 5 seconds of form submission, utilizing the dynamically configured SMTP settings.
- **SC-003**: Role modifications are persisted to the database and reflect immediately upon the next page load or API fetch.
- **SC-004**: System Configuration changes (e.g. SMTP passwords) are securely encrypted at rest.
- **SC-005**: Zero hardcoded mock users remain in the frontend production bundle.

## Assumptions

- The backend database schema already supports users and roles, and only the API layer and frontend integration need to be built.
- The `admin@example.com` default seeded user will serve as the initial Administrator for testing.
