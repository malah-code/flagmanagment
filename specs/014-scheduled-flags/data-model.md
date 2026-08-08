# Data Model: Scheduled Flag Changes

**Feature**: `014-scheduled-flags`
**Phase**: 1 — Design & Contracts

---

## New Entity: ScheduledChange

### Table: `scheduled_changes`

| Column           | Type             | Constraints                                      | Description                                       |
|------------------|------------------|--------------------------------------------------|---------------------------------------------------|
| `id`             | `UUID`           | PK, DEFAULT gen_random_uuid()                    | Primary key                                       |
| `project_id`     | `UUID`           | NOT NULL, FK → projects(id) ON DELETE CASCADE    | Owning project                                    |
| `environment_id` | `UUID`           | NOT NULL, FK → environments(id) ON DELETE CASCADE| Target environment                                |
| `target_type`    | `VARCHAR(20)`    | NOT NULL, CHECK IN ('FLAG','CHANGE_REQUEST')     | Discriminator for the target entity               |
| `target_id`      | `UUID`           | NOT NULL                                         | FK to feature_flags.id or change_requests.id      |
| `action`         | `VARCHAR(20)`    | NOT NULL, CHECK IN ('ENABLE','DISABLE','APPLY')  | Action to perform on execution                    |
| `scheduled_for`  | `TIMESTAMPTZ`    | NOT NULL                                         | UTC timestamp when action should execute          |
| `status`         | `VARCHAR(20)`    | NOT NULL, DEFAULT 'PENDING', CHECK IN ('PENDING','EXECUTED','CANCELLED') | Lifecycle state |
| `created_by`     | `UUID`           | NOT NULL, FK → users(id)                         | User who created the schedule                     |
| `executed_at`    | `TIMESTAMPTZ`    | NULL                                             | UTC timestamp when action was executed            |
| `cancelled_at`   | `TIMESTAMPTZ`    | NULL                                             | UTC timestamp when cancelled                      |
| `created_at`     | `TIMESTAMPTZ`    | NOT NULL, DEFAULT NOW()                          | Record creation time                              |
| `updated_at`     | `TIMESTAMPTZ`    | NOT NULL, DEFAULT NOW()                          | Last update time                                  |

### Indexes

```sql
-- Primary conflict-prevention index: only one PENDING schedule per flag
CREATE UNIQUE INDEX uq_scheduled_changes_pending_flag
  ON scheduled_changes (target_id)
  WHERE status = 'PENDING' AND target_type = 'FLAG';

-- Efficient scheduler polling: find due PENDING records quickly
CREATE INDEX idx_scheduled_changes_due
  ON scheduled_changes (scheduled_for, status)
  WHERE status = 'PENDING';

-- Lookup by environment for listing
CREATE INDEX idx_scheduled_changes_env
  ON scheduled_changes (environment_id, status);
```

### Validation Rules

- `scheduled_for` MUST be strictly in the future at creation time (validated at API layer, not DB layer).
- `action` MUST be `'ENABLE'` or `'DISABLE'` when `target_type = 'FLAG'`.
- `action` MUST be `'APPLY'` when `target_type = 'CHANGE_REQUEST'`.
- Only one `PENDING` schedule may exist per `target_id` where `target_type = 'FLAG'` (enforced by partial unique index).
- `status` transitions are one-way: `PENDING → EXECUTED` or `PENDING → CANCELLED`. Terminal states cannot be changed.

### State Transitions

```
PENDING ──(scheduled_for reached)──► EXECUTED
PENDING ──(user cancels)───────────► CANCELLED
```

---

## Go Model

```go
// internal/models/scheduled_change.go
package models

import (
    "time"
    "github.com/google/uuid"
)

type ScheduledChangeTargetType string
type ScheduledChangeAction     string
type ScheduledChangeStatus     string

const (
    TargetTypeFlag          ScheduledChangeTargetType = "FLAG"
    TargetTypeChangeRequest ScheduledChangeTargetType = "CHANGE_REQUEST"

    ActionEnable  ScheduledChangeAction = "ENABLE"
    ActionDisable ScheduledChangeAction = "DISABLE"
    ActionApply   ScheduledChangeAction = "APPLY"

    ScheduleStatusPending   ScheduledChangeStatus = "PENDING"
    ScheduleStatusExecuted  ScheduledChangeStatus = "EXECUTED"
    ScheduleStatusCancelled ScheduledChangeStatus = "CANCELLED"
)

type ScheduledChange struct {
    ID            uuid.UUID                 `json:"id"             db:"id"`
    ProjectID     uuid.UUID                 `json:"project_id"     db:"project_id"`
    EnvironmentID uuid.UUID                 `json:"environment_id" db:"environment_id"`
    TargetType    ScheduledChangeTargetType `json:"target_type"    db:"target_type"`
    TargetID      uuid.UUID                 `json:"target_id"      db:"target_id"`
    Action        ScheduledChangeAction     `json:"action"         db:"action"`
    ScheduledFor  time.Time                 `json:"scheduled_for"  db:"scheduled_for"`
    Status        ScheduledChangeStatus     `json:"status"         db:"status"`
    CreatedBy     uuid.UUID                 `json:"created_by"     db:"created_by"`
    ExecutedAt    *time.Time                `json:"executed_at,omitempty" db:"executed_at"`
    CancelledAt   *time.Time                `json:"cancelled_at,omitempty" db:"cancelled_at"`
    CreatedAt     time.Time                 `json:"created_at"     db:"created_at"`
    UpdatedAt     time.Time                 `json:"updated_at"     db:"updated_at"`
}
```

---

## Relationships to Existing Models

```
projects ──────────────┐
                       │ project_id
environments ──────────┤
                       │ environment_id
users ─────────────────┤ created_by
                       ▼
               scheduled_changes
                    │ target_id (polymorphic)
          ┌─────────┴─────────┐
          │                   │
   feature_flags      change_requests
   (when target_type   (when target_type
    = 'FLAG')           = 'CHANGE_REQUEST')
```

---

## Migration Number

Next migration: `000012_create_scheduled_changes`

```sql
-- 000012_create_scheduled_changes.up.sql
CREATE TABLE scheduled_changes (
    id             UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    project_id     UUID        NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    environment_id UUID        NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    target_type    VARCHAR(20) NOT NULL CHECK (target_type IN ('FLAG', 'CHANGE_REQUEST')),
    target_id      UUID        NOT NULL,
    action         VARCHAR(20) NOT NULL CHECK (action IN ('ENABLE', 'DISABLE', 'APPLY')),
    scheduled_for  TIMESTAMPTZ NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'EXECUTED', 'CANCELLED')),
    created_by     UUID        NOT NULL REFERENCES users(id),
    executed_at    TIMESTAMPTZ,
    cancelled_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_scheduled_changes_pending_flag
  ON scheduled_changes (target_id)
  WHERE status = 'PENDING' AND target_type = 'FLAG';

CREATE INDEX idx_scheduled_changes_due
  ON scheduled_changes (scheduled_for, status)
  WHERE status = 'PENDING';

CREATE INDEX idx_scheduled_changes_env
  ON scheduled_changes (environment_id, status);

-- 000012_create_scheduled_changes.down.sql
DROP TABLE IF EXISTS scheduled_changes;
```
