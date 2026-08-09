# Authentication API Contracts

## OIDC Login Flow
`GET /auth/login/oidc`
- Initiates OIDC flow. Redirects user to IdP.

`GET /auth/callback/oidc`
- Callback URL for OIDC IdP.
- Exchanges code for token, authenticates user, sets HTTP-only session cookie.

## SAML Login Flow
`GET /auth/login/saml`
- Initiates SAML login flow. Redirects to IdP.

`POST /auth/callback/saml`
- Accepts SAML Response via POST.
- Validates signature, authenticates user, sets HTTP-only session cookie.

## Service Accounts
`POST /api/v1/service-accounts`
- Request: `{"name": "CI/CD Pipeline", "role": "admin"}`
- Response: `{"id": "uuid", "name": "CI/CD Pipeline", "token": "fm_sa_eyJhbGci..."}` (Token is only returned once during creation).
