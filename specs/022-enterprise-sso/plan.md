# Implementation Plan: Enterprise SSO & Authentication (SAML/OIDC)

## 1. Technical Context

### Architecture & Tech Stack

- **Backend**: Go with `chi` router and custom middleware.
- **Identity Protocol**: OIDC (OpenID Connect) and SAML 2.0. We will use the standard `golang.org/x/oauth2` for OIDC and something like `github.com/crewjam/saml` for SAML 2.0.
- **Database**: PostgreSQL (currently managing `users` and `service_accounts` tables).

### Core Changes

1. **Authentication Handlers**: Need new HTTP handlers in `backend/internal/api/auth.go` for `/auth/login/sso`, `/auth/callback/oidc`, and `/auth/callback/saml`.
2. **SSO Configuration**: We will configure IdP settings via environment variables (e.g. `OIDC_ISSUER_URL`, `OIDC_CLIENT_ID`, `SAML_IDP_METADATA_URL`).
3. **Just-In-Time (JIT) Provisioning**: If a user successfully authenticates via SSO but does not exist in our database, we automatically insert them into `users` with their OIDC subject or SAML NameID as `external_id`, and `auth_provider` set to "OIDC" or "SAML".
4. **Service Accounts**: Enhance existing `/api/v1/service-accounts` endpoints to allow creation and generation of API tokens (JWT or opaque). Update `middleware_jwt.go` to accept service account tokens.

### Open Questions / NEEDS CLARIFICATION

None. The spec clearly defines that roles won't be mapped directly from the IdP yet; they just get a default global role upon creation.
