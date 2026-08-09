# Feature Specification: Enterprise SSO & Authentication (SAML/OIDC)

**Feature Branch**: `022-enterprise-sso`

**Created**: 2026-08-08

**Status**: Draft

**Input**: User description: "Enterprise SSO & Authentication (SAML/OIDC)"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Secure Identity Provider Integration (Priority: P1)

As a System Administrator, I want to configure the platform to authenticate users via our corporate Identity Provider (IdP) using SAML 2.0 or OIDC, so that we can enforce central access controls, MFA, and lifecycle management.

**Why this priority**: Required for enterprise adoption and SOC2 compliance.

**Independent Test**: Can be tested by configuring a mock IdP (e.g. Auth0 or Keycloak), initiating a login flow, and successfully authenticating and returning to the dashboard with an active session.

**Acceptance Scenarios**:

1. **Given** a correctly configured OIDC integration, **When** a user clicks "Log in with SSO", **Then** they are redirected to the IdP, authenticate successfully, and are redirected back to the dashboard with a valid session.
2. **Given** an IdP integration, **When** a user's account is disabled in the IdP, **Then** their next login attempt via SSO fails.

---

### User Story 2 - Just-In-Time (JIT) User Provisioning (Priority: P2)

As a System Administrator, I want new users who authenticate via SSO to be automatically provisioned in FlagManagment with a default role, so that I don't have to manually invite every developer.

**Why this priority**: Reduces administrative overhead for large teams.

**Independent Test**: Can be tested by logging in with a new user account via the IdP and verifying the user record is created in the database with the default "Read-Only" role.

**Acceptance Scenarios**:

1. **Given** JIT provisioning is enabled, **When** a new user successfully authenticates via SSO, **Then** a user record is created and assigned the default global role.

---

### User Story 3 - Service Account / API Key Authentication (Priority: P1)

As a Developer, I want to authenticate my CI/CD pipelines and external automation tools using robust API keys or service accounts instead of user-based SSO, so that automation can interact with the API securely.

**Why this priority**: SSO is for humans; we need machine-to-machine authentication to enable Terraform and CI/CD pipelines.

**Independent Test**: Can be tested by generating a Service Account API key and successfully invoking a protected REST endpoint.

**Acceptance Scenarios**:

1. **Given** a valid Service Account API key, **When** a request is made to the REST API, **Then** the system authenticates the request and enforces the Service Account's RBAC permissions.

### Edge Cases

- What happens if the IdP goes down? (Fallback to local admin auth for emergency access).
- What happens if a user's role in the IdP changes? (Do we sync roles? Wait, let's keep it simple: we assign a default role on creation, and admins manage RBAC within FlagManagment).
- How do we handle mixed authentication (some users local, some SSO)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST support OIDC (OpenID Connect) for authenticating dashboard users.
- **FR-002**: The system MUST support SAML 2.0 for authenticating dashboard users.
- **FR-003**: The system MUST support Just-In-Time (JIT) provisioning for SSO users, creating local records with a configurable default role.
- **FR-004**: The system MUST provide an "Emergency Admin" local login bypass in case of IdP misconfiguration or downtime.
- **FR-005**: The system MUST support the creation of non-human "Service Accounts" with API keys for machine-to-machine automation.
- **FR-006**: Service Accounts MUST be assignable to standard RBAC roles (Global, Project, Environment).

### Key Entities

- **SSOConfiguration**: Stores the IdP issuer URL, Client ID, Client Secret (encrypted), or SAML XML metadata.
- **User**: Needs a flag `AuthProvider` (Local, OIDC, SAML) and an `ExternalID` mapping to the IdP subject.
- **ServiceAccount**: A non-human actor representing automation, capable of holding roles and generating API keys.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can successfully authenticate via Auth0/Okta using OIDC or SAML in under 3 seconds.
- **SC-002**: Automated scripts can authenticate using Service Account tokens and execute API calls.
- **SC-003**: 100% of authentication events (success and failure) are logged in the Immutable Audit Log.

## Assumptions

- For the MVP, role mapping (syncing IdP groups to FlagManagment roles) is excluded. We will use JIT provisioning with a default role, and admins will manage specific permissions locally.
- SAML/OIDC configuration will be done via environment variables or a secure configuration file for the initial release to avoid building a complex UI for IdP configuration immediately.
