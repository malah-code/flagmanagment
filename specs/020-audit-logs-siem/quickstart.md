# Quickstart Validation: Audit Logs & SIEM Webhooks

This guide documents how to validate that the Immutable Audit Log and SIEM Webhook integration are functioning end-to-end.

## Prerequisites

- FlagManagment backend running locally (`make run`).
- `curl` or Postman for making API requests.
- A local webhook testing tool like [Webhook.site](https://webhook.site/) or `ngrok`.

## Scenario 1: Validating Audit Log Creation and Sanitization

**1. Create an Environment (Mutating Action)**
```bash
# Create an environment (this creates an API key internally)
curl -X POST http://localhost:33363/api/v1/projects/<project-id>/environments \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Audit Test Env"}'
```

**2. Query the Audit Log**
```bash
curl -X GET "http://localhost:33363/api/v1/audit-logs?project_id=<project-id>&limit=1" \
  -H "Authorization: Bearer <admin-token>"
```

**Validation**:
- The API returns a `200 OK` containing the audit log entry.
- The `action_type` is `ENVIRONMENT_CREATED`.
- In the `new_state` JSON payload, the `api_key` field MUST explicitly read `"[REDACTED]"`. The plaintext API key must NOT be present.

## Scenario 2: Validating SIEM Webhook Dispatch

**1. Register a Webhook Integration**
Obtain a unique URL from [Webhook.site](https://webhook.site/) or spin up a local listener.
```bash
curl -X POST http://localhost:33363/api/v1/projects/<project-id>/webhooks \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://webhook.site/your-unique-id", "events": ["audit.*"]}'
```

**2. Perform a Flag Mutation**
```bash
curl -X POST http://localhost:33363/api/v1/projects/<project-id>/flags \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"key": "test-webhook-flag", "type": "boolean"}'
```

**Validation**:
- Check your Webhook.site dashboard.
- You should see an incoming `POST` request within 2 seconds.
- The payload should match the `Webhook Delivery Payload` schema from `contracts/api.md`, detailing the `FLAG_CREATED` action.

## Scenario 3: Validating CSV Export

**1. Request CSV Export**
```bash
curl -X GET "http://localhost:33363/api/v1/audit-logs/export?project_id=<project-id>" \
  -H "Authorization: Bearer <admin-token>" \
  -o audit_export.csv
```

**Validation**:
- The file `audit_export.csv` is downloaded successfully.
- Opening the file in a text editor reveals standard comma-separated values containing the log entries, including the flag creation event from Scenario 2.
