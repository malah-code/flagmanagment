# Quickstart: Testing Enterprise SSO & Service Accounts

## Prerequisites
- A local FlagManagment backend running.
- (For OIDC testing) An Auth0 or Google Cloud OIDC application configured with callback `http://localhost:8080/auth/callback/oidc`.

## Validating Service Accounts
1. **Create a Service Account:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/service-accounts \
     -H "Content-Type: application/json" \
     -d '{"name": "Terraform SA", "role": "admin"}'
   ```
2. **Extract Token:** Save the `token` from the JSON response.
3. **Authenticate as Service Account:**
   ```bash
   curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/v1/projects
   ```
   (Should return 200 OK, not 401 Unauthorized)

## Validating OIDC Login
1. Start the server with `OIDC_ISSUER_URL`, `OIDC_CLIENT_ID`, and `OIDC_CLIENT_SECRET` environment variables.
2. Navigate browser to `http://localhost:8080/auth/login/oidc`.
3. You should be redirected to the IdP, authenticate, and land back on the dashboard with a valid session cookie.
