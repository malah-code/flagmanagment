# Implementation Tasks: Enterprise SSO & Authentication (SAML/OIDC)

## Phase 1: Setup

- [ ] T001 Define required IdP environment variables in `.env.example`
- [ ] T002 Add Go dependencies (`go get coreos/go-oidc/v3 github.com/crewjam/saml`)

## Phase 2: Foundational

- [ ] T003 [P] Update `User` model to include `AuthProvider` and `ExternalID` in `backend/internal/models/user.go`
- [ ] T004 [P] Update `ServiceAccountKey` model in `backend/internal/models/service_account.go`
- [ ] T005 Create database migrations for `users` and `service_account_keys` updates

## Phase 3: [US1] Secure Identity Provider Integration

**Goal**: Support OIDC and SAML logins for the admin dashboard.
**Independent Test**: Navigate to `/auth/login/oidc`, authenticate via IdP, and verify a valid session is created.

- [ ] T006 [US1] Implement OIDC HTTP handlers (`/auth/login/oidc`, `/auth/callback/oidc`) in `backend/internal/api/auth.go`
- [ ] T007 [P] [US1] Implement SAML HTTP handlers (`/auth/login/saml`, `/auth/callback/saml`) in `backend/internal/api/auth.go`

## Phase 4: [US2] Just-In-Time (JIT) User Provisioning

**Goal**: Automatically provision users in the database upon successful IdP authentication.
**Independent Test**: Login with a new IdP account; verify the row exists in `users` with the appropriate `auth_provider`.

- [ ] T008 [US2] Extract Subject/NameID from IdP response and implement JIT insertion logic in `backend/internal/api/auth.go` (if user does not exist)
- [ ] T009 [US2] Assign default global role (e.g., Viewer) to newly provisioned users

## Phase 5: [US3] Service Account / API Key Authentication

**Goal**: Machine-to-machine authentication via Service Accounts and opaque/JWT API keys.
**Independent Test**: Call `/api/v1/service-accounts` to generate a key, then use it as a Bearer token.

- [ ] T010 [US3] Implement Service Account creation handler in `backend/internal/api/service_accounts.go`
- [ ] T011 [US3] Implement API Key generation (cryptographically secure opaque token), hash it via SHA-256, and store in `service_account_keys`
- [ ] T012 [P] [US3] Update `backend/internal/api/middleware_jwt.go` to accept and validate Service Account Bearer tokens (via checking the hash in the DB)

## Phase 6: Polish

- [ ] T013 Update `backend/internal/api/middleware_rbac.go` to ensure Service Account roles are properly respected for API authorization
- [ ] T014 Write integration tests for Service Account token generation and validation

## Dependencies
- US1 (OIDC/SAML integration) depends on T003 and T005.
- US2 (JIT Provisioning) depends on US1.
- US3 (Service Accounts) depends on T004 and T005, and can be implemented in parallel with US1/US2.
