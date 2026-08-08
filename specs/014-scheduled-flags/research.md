# Research: Scheduled Flag Changes

**Feature**: `014-scheduled-flags`
**Phase**: 0 — Outline & Research
**Date**: 2026-08-08

---

## 1. Background Scheduler Pattern (Go)

**Decision**: Use a time-ticker-based polling goroutine (in-process background worker) started at server startup, scanning the `scheduled_changes` table every 30 seconds.

**Rationale**:
- Aligned with the spec assumption: "a polling-based background worker or cron-like mechanism in the Go backend rather than relying on external scheduling services."
- The project has no existing job-queue infrastructure; adding one would be over-engineered.
- A 30-second poll tick comfortably satisfies SC-001 (≤60s execution window).
- `time.Ticker` is idiomatic Go with no external dependencies.

**Alternatives Considered**:
- **External cron service (Google Cloud Scheduler)**: Rejected — violates spec assumption of self-contained architecture.
- **Asynq / Redis-backed job queue**: Rejected — adds Redis queue complexity unnecessarily.
- **robfig/cron library**: Would work but adds a dependency for a straightforward polling loop.

---

## 2. Conflict Detection (FR-006)

**Decision**: Enforce a DB-level partial unique index on `(flag_id)` WHERE `status = 'PENDING'` in PostgreSQL, plus an application-level pre-check inside a transaction with FOR UPDATE row lock.

**Rationale**:
- PostgreSQL partial indexes allow only one pending schedule per flag while permitting historical executed/cancelled records.
- Application-level check provides clean error messaging before hitting DB constraint.

**Alternatives Considered**:
- **Application-only check (no DB constraint)**: Rejected — race condition possible under concurrent requests.
- **Full unique index on `(flag_id)`**: Rejected — blocks re-scheduling after execution/cancellation.

---

## 3. Timezone Handling (FR-002)

**Decision**: All `scheduled_for` values stored as UTC `TIMESTAMPTZ` in PostgreSQL. Frontend sends UTC ISO-8601 strings. UI converts to/from user's local timezone via browser's native `Intl.DateTimeFormat` API.

**Rationale**:
- Consistent with the spec assumption: "Timezone handling will be done by sending UTC timestamps from the frontend to the backend."
- PostgreSQL `TIMESTAMPTZ` stores in UTC and handles DST transparently.
- No timezone conversion logic required server-side.

---

## 4. Scheduler Restart Resilience (Edge Case: Missed Triggers)

**Decision**: On every scheduler tick and on startup, query for all `PENDING` schedules where `scheduled_for <= NOW()` and process them as a startup catchup sweep.

**Rationale**:
- Ensures zero triggers are permanently missed after a restart.
- Worst-case delay after restart is one poll interval (30s), well within SC-001's 60s window.

---

## 5. Concurrent Execution Safety (SC-003)

**Decision**: Bounded goroutine worker pool inside the scheduler. Each tick fetches due schedules and dispatches to a buffered channel consumed by N=20 fixed workers.

**Rationale**:
- Prevents spawning 1,000 goroutines simultaneously, which could spike memory.
- Bounded workers prevent DB connection pool exhaustion.
- Wraps each execution in a transaction via the existing `WithTx` pattern.

---

## 6. Audit Log Integration (FR-004)

**Decision**: Reuse `AuditService.LogAction` with `ActorID` set to a well-known system sentinel UUID (`00000000-0000-0000-0000-000000000001`) and `Action = "SCHEDULED_EXECUTION"`.

**Rationale**:
- Existing `AuditLog.ActorID` is non-nullable (uuid.UUID); a sentinel UUID is the cleanest solution.
- Consistent with how `AuditService` is used elsewhere — no schema changes needed to `audit_logs`.

---

## 7. Unified ScheduledChange Entity (US-1 + US-2)

**Decision**: Single `scheduled_changes` table with `target_type IN ('FLAG', 'CHANGE_REQUEST')` handles both flag toggles and CR applications. Scheduler dispatches to the appropriate service based on `target_type`.

**Rationale**:
- A single table and scheduler loop handles both user stories uniformly.
- Avoids schema fragmentation; aligns with the `ScheduledChange` entity defined in the spec.

---

## 8. API Design (Constitution I — API-First)

**Decision**: REST endpoints following existing project conventions (`/api/v1/...`, chi router, JSON responses):

```
POST   /api/v1/environments/{envId}/scheduled-changes   — Create schedule
GET    /api/v1/environments/{envId}/scheduled-changes   — List schedules (with status filter)
GET    /api/v1/scheduled-changes/{id}                   — Get by ID
PATCH  /api/v1/scheduled-changes/{id}                   — Modify schedule time
DELETE /api/v1/scheduled-changes/{id}                   — Cancel schedule
```

Full OpenAPI contract defined in `contracts/scheduled-changes.openapi.yaml`.

---

## 9. RBAC (FR-001)

**Decision**: Create/cancel/modify requires `RELEASE_MANAGER` or `ADMIN` role. Viewing requires `VIEWER` minimum. Mirrors the RBAC pattern used for Change Requests.

---

## 10. Frontend Integration

**Decision**: Add `ScheduledChangeBadge` and `ScheduleDialog` components to the existing flag list/detail view in `ProjectDetail.tsx`. A new `scheduledChangesApi.ts` service wraps the REST endpoints. The dialog uses `<input type="datetime-local">` with JavaScript converting to UTC before sending.

**Rationale**:
- Scheduling is surfaced as a modal panel on the existing flag card — no new page required.
- Keeps the frontend surface minimal for this feature.
