# Data Model: Slack Notification Webhooks

## Entities

### `SlackWebhookConfig`

Stores the incoming Slack webhook URL for an environment.

| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary Key |
| `environment_id` | UUID | Foreign Key -> `environments.id` (Unique) |
| `webhook_url` | Text | Encrypted/Plaintext Slack Incoming Webhook URL |
| `enabled` | Boolean | Whether notifications are enabled for this environment |
| `created_at` | Timestamp | Creation timestamp |
| `updated_at` | Timestamp | Last update timestamp |
