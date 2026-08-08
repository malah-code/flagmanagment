# Quickstart: PII Hashing & Data Privacy

This guide demonstrates how to validate that PII hashing and salting is working correctly in the backend.

## Prerequisites
- FlagManagment backend running locally (`make run` or `go run cmd/server/main.go`).
- Database migrated to the latest version including the `Salt` column on the `environments` table.
- curl or a similar API client.

## Validation 1: Verify Environment Salt Generation

When a new environment is created, it should automatically be assigned a cryptographically secure 32-byte salt.

```bash
# 1. Create a project
PROJECT_ID=$(curl -s -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer dev-admin-token" \
  -d '{"name": "PII Test Project"}' | jq -r '.id')

# 2. Create an environment
ENV_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/projects/$PROJECT_ID/environments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer dev-admin-token" \
  -d '{"name": "Production"}')

# 3. Extract the salt from the database directly to verify it was generated
ENV_ID=$(echo $ENV_RESPONSE | jq -r '.id')
docker exec -it flagmanagment-db psql -U postgres -d flagmanagment -c "SELECT salt FROM environments WHERE id = '$ENV_ID';"
```
**Expected Outcome**: The database should return a 64-character hexadecimal string representing the 32-byte salt.

## Validation 2: Verify PII is Hashed in Evaluation Context

When evaluating a flag, the `HashPII` function should apply SHA-256 + Salt to sensitive fields.

```bash
# 1. Evaluate a flag using an email address
curl -s -X POST http://localhost:8080/api/v1/evaluate \
  -H "Content-Type: application/json" \
  -H "X-Environment-Key: <your_env_key>" \
  -d '{
    "flagKey": "test-flag",
    "context": {
      "entityId": "123",
      "attributes": {
        "email": "user@example.com"
      }
    }
  }'

# 2. Check the server logs or analytics sink (database/redis)
```
**Expected Outcome**: The evaluation result will return successfully, and observing the persistent storage for analytics will show `email` as a hashed string instead of `user@example.com`.

## Validation 3: Verify SDK Bucketing uses MurmurHash3

Bucketing must be deterministic.
Evaluating the exact same user ID and flag key repeatedly with the exact same `Environment` salt must yield the exact same bucket assignment. 
*See unit tests in `evaluator_test.go` for programmatic validation of MurmurHash3 bucketing.*
