# Data Model: Stale Flag Detection

**Feature**: `015-stale-flag-management`

## Enums

```sql
CREATE TYPE flag_lifecycle_state AS ENUM ('ACTIVE', 'STALE', 'DEPRECATED', 'ARCHIVED');
```

## Tables Modified

### `environment_flag_states` (Altered)
- **`lifecycle_state`**: `flag_lifecycle_state` (DEFAULT 'ACTIVE', NOT NULL)
- **`last_evaluated_at`**: `TIMESTAMPTZ` (Nullable)
- **`last_state_change_at`**: `TIMESTAMPTZ` (DEFAULT NOW(), NOT NULL)

*Index Updates:*
- Add index on `(environment_id, lifecycle_state)` for fast dashboard filtering.
- Add index on `(environment_id, last_evaluated_at, last_state_change_at)` to optimize the background staleness scanner.

## New Tables

### `stale_flag_policies` (New)
Defines the configurable thresholds for staleness per environment. If an environment has no policy, the project policy is used. If the project has no policy, the system default (30 days) applies.

- **`id`**: `UUID` (Primary Key)
- **`project_id`**: `UUID` (Foreign Key, Not Null)
- **`environment_id`**: `UUID` (Foreign Key, Nullable - allows project-wide defaults)
- **`stale_after_days`**: `INTEGER` (DEFAULT 30, NOT NULL)
- **`created_at`**: `TIMESTAMPTZ` (NOT NULL)
- **`updated_at`**: `TIMESTAMPTZ` (NOT NULL)

*Constraints:*
- Unique index on `(project_id, environment_id)` ensuring only one policy per scope.

## State Transitions

- `ACTIVE` -> `STALE`: Triggered by background scanner when `last_evaluated_at` > `stale_after_days` OR (`last_state_change_at` > `stale_after_days` AND rollout = 100%).
- `STALE` -> `ACTIVE`: Triggered automatically if a user modifies the flag's rules or boolean state.
- `STALE` / `ACTIVE` -> `DEPRECATED`: Manual action by Release Manager (flags it as "do not use for new code").
- `STALE` / `DEPRECATED` -> `ARCHIVED`: Manual action by Release Manager. Excludes the flag from the active SDK gRPC stream payload.
- `ARCHIVED` -> `ACTIVE`: Manual restore action by Release Manager.
