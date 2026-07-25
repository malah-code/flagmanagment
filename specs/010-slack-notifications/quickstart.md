# Quickstart: Slack Notification Webhooks Validation

## Verification Flow

1. **Configure Webhook**:
   - Make a `POST /api/v1/environments/{envId}/slack` with body `{"webhook_url": "https://hooks.slack.com/services/...", "enabled": true}`.
   - Verify HTTP 200 OK.

2. **Trigger Notification**:
   - Toggle a feature flag state in `{envId}`.
   - Verify that an asynchronous HTTP POST request is dispatched to the configured Slack webhook URL with Slack Block Kit JSON format.
