# Quickstart Validation: Server-Side Keys Management

## 1. Create a Named Server-Side Key
Verify that admins can issue a new named server key via the API.

```bash
# Obtain a project and environment ID
PROJECT_ID="your-project-id"
ENV_ID="your-environment-id"

# Create a new server key
curl -X POST http://localhost:8080/api/v1/projects/$PROJECT_ID/environments/$ENV_ID/server-keys \
  -H "Content-Type: application/json" \
  -d '{"name": "test-service-key"}'
```
**Expected Outcome**: 
- Response `201 Created` with a JSON payload containing `id`, `name`, and `key` (e.g., `env_xxxxxxxxxx`).

## 2. List Server-Side Keys
Verify the key appears in the environment's key list.

```bash
curl -X GET http://localhost:8080/api/v1/projects/$PROJECT_ID/environments/$ENV_ID/server-keys \
  -H "Content-Type: application/json"
```
**Expected Outcome**: 
- Response `200 OK`. The array `data` contains the newly created key, but the `key` field itself is OMITTED (only metadata is returned).

## 3. Verify Local Evaluation Auth with the New Key
Ensure the SDK local evaluation endpoint accepts the newly created key.

```bash
# Use the key obtained in step 1
NEW_SERVER_KEY="env_xxxxxxxxxx"

curl -X GET http://localhost:8080/api/v1/sdk/evaluations \
  -H "Authorization: Bearer $NEW_SERVER_KEY"
```
**Expected Outcome**:
- Response `200 OK` returning the evaluation payload. (If a legacy key or invalid key is used, ensure correct behavior).

## 4. Delete the Server-Side Key
Verify the key can be revoked and deleted.

```bash
KEY_ID="the-id-from-step-1"

curl -X DELETE http://localhost:8080/api/v1/projects/$PROJECT_ID/environments/$ENV_ID/server-keys/$KEY_ID
```
**Expected Outcome**:
- Response `204 No Content`. Subsequent attempts to use this key with `/api/v1/sdk/evaluations` MUST return `401 Unauthorized`.
