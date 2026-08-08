# Quickstart: API-Driven Environment Cloning

This guide explains how to validate the new Ephemeral Environment Cloning and Deletion endpoints using `curl` against a running local instance.

## Prerequisites

1. Ensure the FlagManagment backend is running locally (`make run-backend`).
2. Have a valid API token or JWT for an authenticated admin user.
3. Have an existing Project ID and a Source Environment ID (e.g., "Staging").

## Step 1: Clone an Environment

Execute this API call to clone the Staging environment into a new ephemeral environment for PR-123.

```bash
curl -X POST http://localhost:8080/api/v1/projects/{projectId}/environments/{sourceEnvId}/clone \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "PR-123-Test"}'
```

**Expected Outcome:**
You should receive an HTTP `201 Created` response containing the new environment ID and its newly generated SDK token.

```json
{
  "id": "new-env-uuid-here",
  "name": "PR-123-Test",
  "api_key": "fm-sdk-token-12345..."
}
```

## Step 2: Validate Flag States

Use the newly received SDK token to hit the SDK local evaluation bootstrapping endpoint.

```bash
curl -X GET http://localhost:8080/sdk/v1/flags \
  -H "Authorization: Bearer fm-sdk-token-12345..."
```

**Expected Outcome:**
The response should contain the exact same matrix of flag states that were present in the source environment.

## Step 3: Teardown (Delete) the Ephemeral Environment

When testing is complete, destroy the environment via the API.

```bash
curl -X DELETE http://localhost:8080/api/v1/projects/{projectId}/environments/{new-env-uuid-here} \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Expected Outcome:**
HTTP `204 No Content`. The environment is completely removed from the database, along with its associated flags. Attempting to use the SDK token from Step 1 will now result in an HTTP 401 Unauthorized.
