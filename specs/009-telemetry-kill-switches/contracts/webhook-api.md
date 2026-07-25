# APM Webhook API Contract

## Ingest Alert Endpoint

**URL**: `POST /api/v1/webhooks/apm`

**Authentication**: 
Requires the target Environment's API Key passed as a Bearer token.
`Authorization: Bearer <env_api_key>`

### Request Payload (JSON)

APM tools must configure their webhooks to send a JSON payload containing an `alert_identifier`.

```json
{
  "alert_identifier": "high_error_rate_payment_service",
  "status": "firing", 
  "description": "Error rate exceeded 5% on /payments"
}
```

*Note: While additional fields can be sent by the APM tool, the `alert_identifier` is the only strictly required field the platform uses to match against `kill_switch_rules`.*

### Responses

**202 Accepted**
The webhook was successfully authenticated, parsed, and any matching kill-switch rules were executed.
```json
{
  "status": "processed",
  "flags_killed": 1
}
```

**401 Unauthorized**
The API key was missing or invalid.

**400 Bad Request**
The JSON payload was malformed or missing the `alert_identifier`.
