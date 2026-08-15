# Data Model: Feature Flag Creation & Targeting Workflows

## Entities & Schema

### 1. `FeatureFlag` (`feature_flags` table)

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID (PK) | Unique flag identifier |
| `project_id` | UUID (FK) | Reference to owning project |
| `key` | VARCHAR(255) | Unique alphanumeric flag key (e.g. `enable-new-checkout`) |
| `name` | VARCHAR(255) | Human-readable name |
| `description` | TEXT | Description of flag intent |
| `type` | VARCHAR(50) | `BOOLEAN`, `MULTIVARIATE`, `JSON` |
| `variations` | JSONB | Array of variation objects `[{ id, name, value }]` |
| `tags` | JSONB | Array of string tags `['checkout', 'v2', 'beta']` |
| `created_at` | TIMESTAMPTZ | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | Last update timestamp |

### 2. `EnvironmentFlagState` (`environment_flag_states` table)

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID (PK) | Unique state identifier |
| `environment_id`| UUID (FK) | Reference to environment |
| `flag_id` | UUID (FK) | Reference to feature flag |
| `is_enabled` | BOOLEAN | Global flag toggle status in this environment |
| `rules` | JSONB | Ordered array of targeting rules |
| `rollout_weights`| JSONB | Map of variation IDs to percentage integer weights (sum = 100) |
| `default_variation`| VARCHAR(100)| Fallback variation ID when enabled and no rule matches |
| `off_variation` | VARCHAR(100)| Fallback variation ID when disabled |
| `updated_at` | TIMESTAMPTZ | Last state modification timestamp |

### 3. `TargetingRule` (Structured JSONB inside `rules`)

```json
{
  "id": "rule_01",
  "name": "Beta Testers Only",
  "conditions": [
    {
      "attribute": "email",
      "operator": "ENDS_WITH",
      "values": ["@example.com"]
    },
    {
      "attribute": "tier",
      "operator": "IN_LIST",
      "values": ["enterprise", "pro"]
    }
  ],
  "variation": "var_active",
  "rollout_weights": {
    "var_active": 100
  }
}
```

### 4. `ScheduledChange` (`scheduled_changes` table)

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID (PK) | Unique change identifier |
| `environment_id`| UUID (FK) | Target environment |
| `flag_id` | UUID (FK) | Target flag |
| `scheduled_at` | TIMESTAMPTZ | Scheduled execution timestamp |
| `status` | VARCHAR(50) | `PENDING`, `APPLIED`, `FAILED`, `CANCELLED` |
| `state_payload` | JSONB | State payload to apply |
| `created_by` | UUID (FK) | User who queued the schedule |
| `created_at` | TIMESTAMPTZ | Schedule creation timestamp |
