# Quickstart Validation

## Prerequisites
- A running instance of the FlagManagment backend and frontend.
- An existing Environment and its API Key.
- An existing Feature Flag.

## 1. Configure the Kill Switch
Using the UI (or API), create a Kill Switch Rule for your feature flag:
- **Flag**: `payment-gateway-v2`
- **Alert Identifier**: `payment_errors_high`

## 2. Verify Initial State
Ensure the flag `payment-gateway-v2` is `enabled = true`.

## 3. Simulate an APM Webhook
Send a mock webhook simulating a Datadog/NewRelic alert:

```bash
curl -X POST http://localhost:8080/api/v1/webhooks/apm \
  -H "Authorization: Bearer <YOUR_ENV_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{
    "alert_identifier": "payment_errors_high",
    "status": "firing",
    "description": "5xx errors spiking"
  }'
```

## 4. Expected Outcomes
- The curl command returns `202 Accepted`.
- Refreshing the UI shows the `payment-gateway-v2` flag is now `enabled = false`.
- The Audit Log shows `FLAG_KILLED_AUTOMATICALLY` attributed to the webhook.
