# Quickstart: Granular RBAC and User Management Validation

## Prerequisites
- Backend API running (`make run`)
- Access to the database (e.g., via `psql`)

## Scenario 1: Authentication & Authorization

1. **Create a User** (Backend CLI or DB injection for testing):
   ```sql
   INSERT INTO users (id, email, password_hash) VALUES ('00000000-0000-0000-0000-000000000001', 'admin@example.com', '<bcrypt_hash>');
   ```
2. **Login and Obtain JWT**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
        -H "Content-Type: application/json" \
        -d '{"email":"admin@example.com", "password":"password123"}'
   # Ensure you get a 200 OK and a JWT token in the response
   ```
3. **Verify RBAC Enforcement**:
   Attempt to mutate a flag in Project A using a token from a user with a `VIEWER` role in Project A:
   ```bash
   curl -X POST http://localhost:8080/api/v1/projects/{project_id}/flags/{flag_id}/toggle \
        -H "Authorization: Bearer <VIEWER_JWT>" \
        -d '{"enabled": true}'
   # Ensure you get a 403 Forbidden response.
   ```

## Scenario 2: Immutable Audit Logging

1. **Trigger an Action**:
   Use an `ADMIN` token to successfully toggle a flag (as in Scenario 1, but with an Admin token).
   ```bash
   curl -X POST http://localhost:8080/api/v1/projects/{project_id}/flags/{flag_id}/toggle \
        -H "Authorization: Bearer <ADMIN_JWT>" \
        -d '{"enabled": true}'
   ```
2. **Verify Audit Log Creation**:
   ```bash
   curl -X GET http://localhost:8080/api/v1/projects/{project_id}/audit-logs \
        -H "Authorization: Bearer <ADMIN_JWT>"
   # Ensure the response contains the FLAG_TOGGLED action and both old_state and new_state.
   ```
3. **Verify PII Redaction**:
   Create a new environment (which generates an API key) and verify the audit log:
   ```bash
   # Check the audit log endpoint again for ENVIRONMENT_CREATED
   # Ensure the new_state JSON payload has `api_key: "***"` instead of the real key.
   ```
