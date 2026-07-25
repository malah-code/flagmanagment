# Phase 1: Data Model

## Entities

### `kill_switch_rules` (Table)

Tracks the linkage between a feature flag and an expected APM alert identifier.

| Field | Type | Attributes | Description |
|-------|------|------------|-------------|
| `id` | UUID | Primary Key | Unique ID |
| `flag_id` | UUID | Foreign Key | References `feature_flags(id)` |
| `environment_id` | UUID | Foreign Key | References `environments(id)` |
| `alert_identifier` | VARCHAR(255) | Not Null, Index | The string sent by the APM tool |
| `action` | VARCHAR(50) | Default 'DISABLE' | Action to take when triggered |
| `created_by` | UUID | Foreign Key | References `users(id)` |
| `created_at` | TIMESTAMP | Not Null | Creation time |

**Indexes**:
- Unique compound index on `(environment_id, alert_identifier, flag_id)` to prevent duplicates and allow fast lookups during webhook ingestion.

### `audit_logs` (Existing Table Update)
The `AuditService` will be used with a system actor or the APM tool name as the "actor" when an automated kill switch is triggered. The action type will be `FLAG_KILLED_AUTOMATICALLY`.
