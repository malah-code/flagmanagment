# Data Model: Immutable Audit Logs & SIEM Webhooks

## Entities

### `audit_logs`

The append-only ledger tracking all administrative actions.

| Field | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | Primary Key | Unique identifier for the log entry. |
| `actor_id` | UUID | Foreign Key (`users.id`) | The user who performed the action. |
| `target_project_id` | UUID | Foreign Key (`projects.id`), Nullable | The project affected (if applicable). |
| `target_environment_id` | UUID | Foreign Key (`environments.id`), Nullable | The environment affected (if applicable). |
| `target_flag_id` | UUID | Foreign Key (`feature_flags.id`), Nullable | The feature flag affected (if applicable). |
| `action_type` | VARCHAR | Not Null | E.g., `FLAG_CREATED`, `RULE_UPDATED`, `ROLE_ASSIGNED`. |
| `previous_state` | JSONB | Nullable | The state of the entity before the action (sanitized). |
| `new_state` | JSONB | Nullable | The state of the entity after the action (sanitized). |
| `actor_ip` | VARCHAR | Nullable | The IP address from which the action was initiated. |
| `created_at` | TIMESTAMP | Default `NOW()` | When the action occurred. |

**Indexes**:
- `idx_audit_logs_project_id` on `target_project_id`
- `idx_audit_logs_actor_id` on `actor_id`
- `idx_audit_logs_created_at` on `created_at` (DESC)

### `webhook_integrations`

Configuration for external SIEM streaming.

| Field | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | Primary Key | Unique identifier for the webhook. |
| `project_id` | UUID | Foreign Key (`projects.id`) | The project this webhook belongs to. |
| `url` | VARCHAR | Not Null | The destination URL for the HTTP POST request. |
| `secret_key` | VARCHAR | Nullable | Optional signing key for HMAC validation on the receiver end. |
| `events` | JSONB | Not Null | Array of event types to trigger on (e.g., `["audit.*"]`). |
| `is_active` | BOOLEAN | Default `TRUE` | Whether this webhook is currently active. |
| `created_at` | TIMESTAMP | Default `NOW()` | |
| `updated_at` | TIMESTAMP | Default `NOW()` | |

## Data Flow & Validations

1. **Mutation Hook**: The backend HTTP handlers wrap state-mutating actions (create, update, delete) inside a database transaction.
2. **Sanitization**: Before committing the transaction, the `previous_state` and `new_state` payloads are passed through a `Scrub(payload map[string]any)` function that redacts fields matching a known sensitive list (`password`, `token`, `api_key`).
3. **Commit**: The original entity mutation and the `audit_logs` insert are committed together atomically.
4. **Webhook Trigger**: Immediately after a successful commit, the `audit_logs` entry is sent to a Go channel for asynchronous webhook dispatching to any active `webhook_integrations` matching the project.
