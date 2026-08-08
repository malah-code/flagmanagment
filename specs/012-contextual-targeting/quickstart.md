# Quickstart Validation Guide: Contextual Targeting Engine

## Prerequisites
- A running FlagManagment backend and database.
- Go 1.21+ installed to run SDK evaluations.
- A test project, environment, and feature flag.

## Step 1: Update Flag State with Targeting Rules

Use the API to set a targeting rule for the feature flag.

```bash
curl -X PUT http://localhost:8080/api/v1/projects/{PROJECT_ID}/environments/{ENV_ID}/flags/{FLAG_ID} \
  -H "Authorization: Bearer {API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": false,
    "targeting_rules": {
      "rules": [
        {
          "id": "rule-beta-users",
          "conditions": [
            {
              "attribute": "email",
              "operator": "REGEX",
              "value": ".*@beta\\.com$"
            }
          ],
          "variation": true
        }
      ]
    }
  }'
```

**Expected Outcome:** HTTP 200 OK. The flag is fundamentally `disabled`, but the targeting rule says that users with an `@beta.com` email should receive `true`.

## Step 2: Validate Target Rule Rejection for Bad Regex

Try to update the flag with a malformed regex pattern.

```bash
curl -X PUT http://localhost:8080/api/v1/projects/{PROJECT_ID}/environments/{ENV_ID}/flags/{FLAG_ID} \
  -H "Authorization: Bearer {API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": false,
    "targeting_rules": {
      "rules": [
        {
          "id": "rule-bad-regex",
          "conditions": [
            {
              "attribute": "email",
              "operator": "REGEX",
              "value": "[a-z"
            }
          ],
          "variation": true
        }
      ]
    }
  }'
```

**Expected Outcome:** HTTP 400 Bad Request. The server rejects the invalid regex `[a-z`.

## Step 3: Test SDK Evaluation via Unit Test

Run the backend SDK evaluation tests to verify the rules engine evaluates the context correctly.

```bash
cd backend
go test -v ./internal/sdk/... -run TestTargetingRules
```

**Expected Outcome:** The tests pass, verifying that passing `{"email": "user@beta.com"}` returns `true`, and `{"email": "user@prod.com"}` returns `false` (the fallback value).
