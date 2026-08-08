# Tasks: Scheduled Flag Changes

**Input**: Design documents from `specs/014-scheduled-flags/`

**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅, quickstart.md ✅

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Database & Schema)

**Purpose**: Create the `scheduled_changes` database table and migration files. This phase has no code dependencies and can start immediately.

- [x] T001 Create migration file `backend/migrations/000012_create_scheduled_changes.up.sql` with the full CREATE TABLE statement for the `scheduled_changes` table. The table must have 13 columns: `id` (UUID PK DEFAULT gen_random_uuid()), `project_id` (UUID NOT NULL FK→projects(id) ON DELETE CASCADE), `environment_id` (UUID NOT NULL FK→environments(id) ON DELETE CASCADE), `target_type` (VARCHAR(20) NOT NULL CHECK IN ('FLAG','CHANGE_REQUEST')), `target_id` (UUID NOT NULL), `action` (VARCHAR(20) NOT NULL CHECK IN ('ENABLE','DISABLE','APPLY')), `scheduled_for` (TIMESTAMPTZ NOT NULL), `status` (VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK IN ('PENDING','EXECUTED','CANCELLED')), `created_by` (UUID NOT NULL FK→users(id)), `executed_at` (TIMESTAMPTZ nullable), `cancelled_at` (TIMESTAMPTZ nullable), `created_at` (TIMESTAMPTZ NOT NULL DEFAULT NOW()), `updated_at` (TIMESTAMPTZ NOT NULL DEFAULT NOW()). After the CREATE TABLE, add three indexes: (1) a partial unique index `uq_scheduled_changes_pending_flag` on `(target_id) WHERE status = 'PENDING' AND target_type = 'FLAG'` to enforce FR-006 one-pending-schedule-per-flag constraint, (2) a partial index `idx_scheduled_changes_due` on `(scheduled_for, status) WHERE status = 'PENDING'` for efficient scheduler polling, and (3) an index `idx_scheduled_changes_env` on `(environment_id, status)` for listing by environment. Reference: `specs/014-scheduled-flags/data-model.md` has the exact SQL.

- [x] T002 [P] Create migration file `backend/migrations/000012_create_scheduled_changes.down.sql` containing `DROP TABLE IF EXISTS scheduled_changes;` — this single statement also drops all associated indexes. Follow the same pattern used in existing down migrations like `backend/migrations/000010_create_kill_switches.down.sql`.

---

## Phase 2: Foundational (Go Backend — Model, Repository, Store Interface)

**Purpose**: Create the Go model struct, repository interface, repository implementation, and register them in the Store. These are blocking prerequisites for ALL user stories. No user story work can begin until this phase is complete.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T003 Create the Go model file `backend/internal/models/scheduled_change.go`. Define three custom string types: `ScheduledChangeTargetType`, `ScheduledChangeAction`, and `ScheduledChangeStatus`. Define constants: `TargetTypeFlag = "FLAG"`, `TargetTypeChangeRequest = "CHANGE_REQUEST"`, `ActionEnable = "ENABLE"`, `ActionDisable = "DISABLE"`, `ActionApply = "APPLY"`, `ScheduleStatusPending = "PENDING"`, `ScheduleStatusExecuted = "EXECUTED"`, `ScheduleStatusCancelled = "CANCELLED"`. Define the `ScheduledChange` struct with 13 fields matching the DB columns, using `uuid.UUID` for IDs, `time.Time` for timestamps (`ExecutedAt` and `CancelledAt` are `*time.Time`), and the custom types for `TargetType`, `Action`, `Status`. Use both `json` and `db` struct tags following the same pattern as `backend/internal/models/audit_log.go`. Also define a system sentinel UUID constant: `SystemActorUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")` for the scheduler's audit log actor. Reference: the exact struct definition is in `specs/014-scheduled-flags/data-model.md` under "Go Model".

- [x] T004 Add the `ScheduledChangeRepository` interface to `backend/internal/repository/repository.go`. Add it below the existing `ChangeRequestRepository` interface (around line 88). The interface must define these methods: `Create(ctx context.Context, sc *models.ScheduledChange) error`, `GetByID(ctx context.Context, id uuid.UUID) (*models.ScheduledChange, error)`, `ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*models.ScheduledChange, int, error)`, `GetPendingByTargetID(ctx context.Context, targetID uuid.UUID) (*models.ScheduledChange, error)`, `GetDueSchedules(ctx context.Context, now time.Time, limit int) ([]*models.ScheduledChange, error)`, `MarkExecuted(ctx context.Context, id uuid.UUID, executedAt time.Time) error`, `MarkCancelled(ctx context.Context, id uuid.UUID, cancelledAt time.Time) error`, `UpdateScheduledFor(ctx context.Context, id uuid.UUID, newTime time.Time) error`. Also add `ScheduledChangeRepo() ScheduledChangeRepository` to the `Store` interface (around line 18, after `SlackConfigRepo()`). Add `"time"` to the import block if not already present.

- [x] T005 Create the repository implementation file `backend/internal/repository/scheduled_change_repo.go`. Create a `scheduledChangeRepository` struct with a `*sql.DB` field and a `NewScheduledChangeRepository(db *sql.DB) ScheduledChangeRepository` constructor. Implement all 8 methods from the interface:
    - `Create`: INSERT all 13 columns using positional params `$1` through `$13`. Follow the insert pattern in `backend/internal/repository/kill_switch_repo.go`.
    - `GetByID`: SELECT all columns WHERE `id = $1`, return `ErrNotFound` for `sql.ErrNoRows` (follow the `kill_switch_repo.go` GetByID pattern).
    - `ListByEnvironment`: SELECT with optional status filter. If `status` is empty string, don't filter by status. Use `COUNT(*) OVER()` for total count to support pagination. ORDER BY `scheduled_for DESC`. Use `LIMIT $N OFFSET $N` for pagination. Return `([]*models.ScheduledChange, int, error)` where the int is total count. Follow the pagination pattern in `backend/internal/repository/change_request_repo.go`.
    - `GetPendingByTargetID`: `SELECT * FROM scheduled_changes WHERE target_id = $1 AND status = 'PENDING' LIMIT 1`. Return `ErrNotFound` if no rows. This is used for conflict detection in the service layer.
    - `GetDueSchedules`: `SELECT * FROM scheduled_changes WHERE status = 'PENDING' AND scheduled_for <= $1 ORDER BY scheduled_for ASC LIMIT $2`. The `$1` param is `now time.Time`, `$2` is `limit int` (use 100 as the default when called from the scheduler). This is the core scheduler polling query.
    - `MarkExecuted`: `UPDATE scheduled_changes SET status = 'EXECUTED', executed_at = $2, updated_at = $2 WHERE id = $1 AND status = 'PENDING'`. Only update if current status is PENDING (idempotency guard).
    - `MarkCancelled`: `UPDATE scheduled_changes SET status = 'CANCELLED', cancelled_at = $2, updated_at = $2 WHERE id = $1 AND status = 'PENDING'`. Only update if current status is PENDING.
    - `UpdateScheduledFor`: `UPDATE scheduled_changes SET scheduled_for = $2, updated_at = NOW() WHERE id = $1 AND status = 'PENDING'`. Only update if PENDING.
    Scan all 13 columns in the same order for every SELECT query. Use `&sc.ExecutedAt` and `&sc.CancelledAt` (pointer fields) for nullable columns.

- [x] T006 Register the new repository in the Store by modifying `backend/internal/repository/store.go`. Add a new method `func (s *store) ScheduledChangeRepo() ScheduledChangeRepository { return NewScheduledChangeRepository(s.db) }` following the exact same pattern as the existing `KillSwitchRepo()` method on line 40. Place it after the `UserRepo()` method (line 53).

- [x] T007 [P] Create the DTO file `backend/internal/dto/scheduled_change.go`. Define three structs following the existing naming conventions in `backend/internal/dto/requests.go`:
    - `CreateScheduledChangeRequest` with fields: `TargetType string` (json:"target_type", validate:"required,oneof=FLAG CHANGE_REQUEST"), `TargetID string` (json:"target_id", validate:"required,uuid"), `Action string` (json:"action", validate:"required,oneof=ENABLE DISABLE APPLY"), `ScheduledFor string` (json:"scheduled_for", validate:"required").
    - `UpdateScheduledChangeRequest` with field: `ScheduledFor string` (json:"scheduled_for", validate:"required").
    - `ScheduledChangeResponse` with fields mapping all 13 model fields using `string` for UUIDs and `string` for timestamps (following the camelCase JSON naming convention used in `backend/internal/dto/responses.go`): `ID`, `ProjectID`, `EnvironmentID`, `TargetType`, `TargetID`, `Action`, `ScheduledFor`, `Status`, `CreatedBy`, `ExecutedAt` (pointer `*string`), `CancelledAt` (pointer `*string`), `CreatedAt`, `UpdatedAt`.

**Checkpoint**: Foundation ready — the scheduled_changes table, model, repository, DTOs, and Store interface are all in place. User story implementation can now begin.

---

## Phase 3: User Story 1 — Schedule Flag State Change (Priority: P1) 🎯 MVP

**Goal**: A Release Manager can schedule a flag to turn ON or OFF at a future UTC time. The flag state automatically updates when the scheduled time arrives. The schedule can be viewed, modified, and cancelled. Conflicting schedules for the same flag are rejected. An audit log entry is created when the scheduler executes.

**Independent Test**: Create a schedule for a flag via `POST /api/v1/environments/{envId}/scheduled-changes` with `target_type=FLAG`, wait for the scheduler tick after the scheduled time, and verify the flag state has changed and an audit log with `action=SCHEDULED_EXECUTION` exists.

### Backend Service Layer (US1)

- [x] T008 [US1] Create the service file `backend/internal/services/scheduled_change_service.go`. Define a `ScheduledChangeService` struct with dependencies: `store repository.Store` and `audit *AuditService`. Create a constructor `NewScheduledChangeService(store repository.Store, audit *AuditService) *ScheduledChangeService`. Define two package-level sentinel errors: `var ErrPendingScheduleExists = errors.New("a pending schedule already exists for this flag")` and `var ErrScheduleNotPending = errors.New("scheduled change is not in PENDING state")`. Implement the following methods:

    **`Create(ctx context.Context, sc *models.ScheduledChange) error`**:
    1. Validate that `sc.ScheduledFor` is strictly in the future (`time.Now().UTC()`). Return a descriptive error if not.
    2. Validate action/target_type consistency: if `target_type == FLAG`, action must be ENABLE or DISABLE; if `target_type == CHANGE_REQUEST`, action must be APPLY. Return error on mismatch.
    3. Generate a new UUID for `sc.ID` if it's nil. Set `sc.Status = ScheduleStatusPending`, `sc.CreatedAt` and `sc.UpdatedAt` to `time.Now().UTC()`.
    4. Call `s.store.ScheduledChangeRepo().GetPendingByTargetID(ctx, sc.TargetID)`. If a record is found (no error), return `ErrPendingScheduleExists`. If the error is `repository.ErrNotFound`, proceed. If it's another error, return it.
    5. Call `s.store.ScheduledChangeRepo().Create(ctx, sc)`.
    6. Log an audit entry via `s.audit.LogAction(ctx, &models.AuditLog{...})` with Action="CREATE", TargetType="SCHEDULED_CHANGE", TargetID=sc.ID.

    **`GetByID(ctx context.Context, id uuid.UUID) (*models.ScheduledChange, error)`**: Delegate to `s.store.ScheduledChangeRepo().GetByID(ctx, id)`.

    **`ListByEnvironment(ctx context.Context, envID uuid.UUID, status string, limit, offset int) ([]*models.ScheduledChange, int, error)`**: Delegate to `s.store.ScheduledChangeRepo().ListByEnvironment(ctx, envID, status, limit, offset)`.

    **`Cancel(ctx context.Context, id uuid.UUID, actorID uuid.UUID) (*models.ScheduledChange, error)`**:
    1. Fetch the schedule via `GetByID`. Return error if not found.
    2. If `sc.Status != ScheduleStatusPending`, return `ErrScheduleNotPending`.
    3. Call `s.store.ScheduledChangeRepo().MarkCancelled(ctx, id, time.Now().UTC())`.
    4. Log audit entry with Action="CANCEL", TargetType="SCHEDULED_CHANGE".
    5. Re-fetch and return the updated record.

    **`UpdateScheduledFor(ctx context.Context, id uuid.UUID, newTime time.Time, actorID uuid.UUID) (*models.ScheduledChange, error)`**:
    1. Validate `newTime` is in the future. Return error if not.
    2. Fetch the schedule, verify status is PENDING (return `ErrScheduleNotPending` if not).
    3. Call `s.store.ScheduledChangeRepo().UpdateScheduledFor(ctx, id, newTime)`.
    4. Log audit entry with Action="UPDATE", TargetType="SCHEDULED_CHANGE".
    5. Re-fetch and return the updated record.

- [x] T009 [US1] Create the background scheduler file `backend/internal/services/scheduler.go`. This is the core of the feature. Define a `Scheduler` struct with fields: `store repository.Store`, `audit *AuditService`, `cacheClient *cache.Client` (import from `github.com/flagmanagment/backend/internal/cache`), `logger zerolog.Logger` (import from `github.com/rs/zerolog`), `interval time.Duration`, `workerCount int`, `stopCh chan struct{}`. Create a constructor `NewScheduler(store repository.Store, audit *AuditService, cacheClient *cache.Client, logger zerolog.Logger) *Scheduler` that sets `interval` to `30 * time.Second`, `workerCount` to 20, and initializes `stopCh`.

    **`Start(ctx context.Context)`** method:
    1. Log `logger.Info().Msg("scheduler started")`.
    2. Create a `time.NewTicker(s.interval)`.
    3. Run a startup catchup immediately by calling `s.processDueSchedules(ctx)` once before entering the loop.
    4. Enter a `for { select { case <-ticker.C: ... case <-s.stopCh: ... case <-ctx.Done(): ... } }` loop.
    5. On each tick, call `s.processDueSchedules(ctx)`.
    6. On stop or context cancellation, stop the ticker and return.

    **`Stop()`** method: Close `s.stopCh`.

    **`processDueSchedules(ctx context.Context)`** method:
    1. Call `s.store.ScheduledChangeRepo().GetDueSchedules(ctx, time.Now().UTC(), 100)`.
    2. If error, log it and return (do not crash).
    3. If no schedules, return silently.
    4. Log `logger.Info().Int("count", len(schedules)).Msg("processing due scheduled changes")`.
    5. Create a buffered channel `jobs := make(chan *models.ScheduledChange, len(schedules))`.
    6. Create a `sync.WaitGroup`. Launch `min(s.workerCount, len(schedules))` goroutines, each reading from `jobs` and calling `s.executeOne(ctx, sc)`.
    7. Send all schedules to the `jobs` channel, then close it.
    8. Wait for all workers to complete.

    **`executeOne(ctx context.Context, sc *models.ScheduledChange)`** method:
    1. Log the execution attempt: `logger.Info().Str("id", sc.ID.String()).Str("target_type", string(sc.TargetType)).Msg("executing scheduled change")`.
    2. Switch on `sc.TargetType`:
        - **`TargetTypeFlag`**: Fetch the flag state via `s.store.FlagStateRepo().GetByEnvAndFlag(ctx, sc.EnvironmentID, sc.TargetID)`. If error, log and mark as failed. Otherwise, set `state.Enabled = (sc.Action == ActionEnable)` and `state.UpdatedAt = time.Now().UTC()`. Call `s.store.FlagStateRepo().Update(ctx, state)`. If error, log and return.
        - **`TargetTypeChangeRequest`**: This is Phase 4 work — for now, log a warning `"CHANGE_REQUEST scheduling not yet implemented"` and skip. (The US2 phase will fill this in.)
    3. Call `s.store.ScheduledChangeRepo().MarkExecuted(ctx, sc.ID, time.Now().UTC())`. If error, log it.
    4. Create an audit log entry using `s.audit.LogAction(ctx, &models.AuditLog{...})` with `ActorID: models.SystemActorUUID`, `Action: "SCHEDULED_EXECUTION"`, `TargetType: string(sc.TargetType)`, `TargetID: sc.TargetID`, `ProjectID: &sc.ProjectID`, `EnvironmentID: &sc.EnvironmentID`.
    5. If `s.cacheClient != nil`, call `s.cacheClient.PublishRulesetUpdate(ctx, sc.EnvironmentID.String(), time.Now().UTC().String())` to notify connected SDKs of the state change.
    6. Log success: `logger.Info().Str("id", sc.ID.String()).Msg("scheduled change executed successfully")`.
    7. Wrap the entire method in a `defer func() { if r := recover(); r != nil { s.logger.Error()... } }()` to prevent panics from crashing the scheduler.

### Backend API Layer (US1)

- [x] T010 [US1] Create the API handler file `backend/internal/api/scheduled_changes.go`. Define a `ScheduledChangeHandler` struct with fields: `store repository.Store`, `scService *services.ScheduledChangeService`, `rbac *RBACMiddleware`, `cacheClient *cache.Client`, `validate *validator.Validate`. Create a constructor `NewScheduledChangeHandler(store repository.Store, scService *services.ScheduledChangeService, rbac *RBACMiddleware, cacheClient *cache.Client) *ScheduledChangeHandler` that also initializes `validate: validator.New()`.

    **`RegisterRoutes(r chi.Router)`**: Register all 5 endpoints with RBAC middleware, following the exact same pattern as `backend/internal/api/change_requests.go` lines 32-37:
    ```
    r.With(h.rbac.RequireRole("RELEASE_MANAGER")).Post("/environments/{envId}/scheduled-changes", h.Create)
    r.With(h.rbac.RequireRole("VIEWER")).Get("/environments/{envId}/scheduled-changes", h.List)
    r.With(h.rbac.RequireRole("VIEWER")).Get("/scheduled-changes/{id}", h.GetByID)
    r.With(h.rbac.RequireRole("RELEASE_MANAGER")).Patch("/scheduled-changes/{id}", h.Update)
    r.With(h.rbac.RequireRole("RELEASE_MANAGER")).Delete("/scheduled-changes/{id}", h.Cancel)
    ```

    **`Create(w http.ResponseWriter, r *http.Request)`**:
    1. Parse `envId` from URL using `chi.URLParam(r, "envId")` and `uuid.Parse()`. Return 400 on error.
    2. Extract `actorID` from context using `r.Context().Value(UserIDKey).(string)` and `uuid.Parse()`. Return 401 if missing (follow exact pattern from `backend/internal/api/change_requests.go` lines 94-99).
    3. Decode request body into `dto.CreateScheduledChangeRequest`. Validate with `h.validate.Struct(req)`. Return 400 on validation error.
    4. Parse `req.ScheduledFor` as `time.Parse(time.RFC3339, req.ScheduledFor)`. Return 400 on parse error.
    5. Parse `req.TargetID` as UUID. Return 400 on error.
    6. Look up the environment to get `projectID`: call `h.store.EnvironmentRepo().GetByID(ctx, envID)`. Return 404 if not found.
    7. Build a `models.ScheduledChange` struct from the parsed fields, setting `ProjectID` from the environment lookup.
    8. Call `h.scService.Create(ctx, &sc)`. Handle errors: if `services.ErrPendingScheduleExists`, return HTTP 409 with `RespondWithError(w, http.StatusConflict, err.Error())`. Otherwise return 500.
    9. On success, `RespondWithJSON(w, http.StatusCreated, sc)`.

    **`List(w http.ResponseWriter, r *http.Request)`**:
    1. Parse `envId`. Read optional `status` query param via `r.URL.Query().Get("status")`.
    2. Get pagination using `GetPagination(r)` and `TokenToOffset(pagination.PageToken)` (same pattern as `backend/internal/api/change_requests.go` line 48-49).
    3. Call `h.scService.ListByEnvironment(ctx, envID, status, pagination.PageSize, offset)`.
    4. Build `dto.PaginatedResponse{Data: schedules}`. If `offset+len(schedules) < total`, set `NextPageToken` using `OffsetToToken()`.
    5. `RespondWithJSON(w, http.StatusOK, resp)`.

    **`GetByID(w http.ResponseWriter, r *http.Request)`**:
    1. Parse `id` from URL. Call `h.scService.GetByID(ctx, id)`.
    2. Handle `repository.ErrNotFound` → 404. Other errors → 500.
    3. `RespondWithJSON(w, http.StatusOK, sc)`.

    **`Update(w http.ResponseWriter, r *http.Request)`**:
    1. Parse `id`. Extract `actorID` from context.
    2. Decode body into `dto.UpdateScheduledChangeRequest`. Validate.
    3. Parse `req.ScheduledFor` as `time.Parse(time.RFC3339, ...)`.
    4. Call `h.scService.UpdateScheduledFor(ctx, id, parsedTime, actorID)`.
    5. Handle errors: `ErrScheduleNotPending` → 400, `ErrNotFound` → 404, other → 500.
    6. `RespondWithJSON(w, http.StatusOK, updated)`.

    **`Cancel(w http.ResponseWriter, r *http.Request)`**:
    1. Parse `id`. Extract `actorID` from context.
    2. Call `h.scService.Cancel(ctx, id, actorID)`.
    3. Handle errors: `ErrScheduleNotPending` → 400, `ErrNotFound` → 404, other → 500.
    4. `RespondWithJSON(w, http.StatusOK, cancelled)`.

### Backend Wiring (US1)

- [x] T011 [US1] Modify `backend/cmd/server/main.go` to wire the scheduled changes feature. Make these additions following the existing wiring patterns on lines 65-91:
    1. **Import** `"context"` at the top if not present.
    2. After `crService` initialization (line 76), add: `scService := services.NewScheduledChangeService(store, auditService)`.
    3. After `slackHandler` initialization (line 90), add: `scHandler := api.NewScheduledChangeHandler(store, scService, rbacMiddleware, cacheClient)`.
    4. Inside the `r.Group(func(r chi.Router) { r.Use(api.UserAuthMiddleware) ... })` block (around line 100-112), add `scHandler.RegisterRoutes(r)` on a new line after `slackHandler.RegisterRoutes(r)` (after line 110).
    5. **Create and start the Scheduler**: Before the gRPC server goroutine (line 130), add:
       ```go
       scheduler := services.NewScheduler(store, auditService, cacheClient, logger)
       ctx, cancel := context.WithCancel(context.Background())
       defer cancel()
       go scheduler.Start(ctx)
       ```
       Import `"github.com/flagmanagment/backend/internal/services"` if not already imported (it is, on line 27).

### Frontend — TypeScript Types & API Client (US1)

- [x] T012 [P] [US1] Create the TypeScript types file `frontend/src/types/scheduledChange.ts`. Define and export the following types following the naming conventions in `frontend/src/types/index.ts`:
    ```typescript
    export type ScheduledChangeTargetType = 'FLAG' | 'CHANGE_REQUEST';
    export type ScheduledChangeAction = 'ENABLE' | 'DISABLE' | 'APPLY';
    export type ScheduledChangeStatus = 'PENDING' | 'EXECUTED' | 'CANCELLED';

    export interface ScheduledChange {
      id: string;
      project_id: string;
      environment_id: string;
      target_type: ScheduledChangeTargetType;
      target_id: string;
      action: ScheduledChangeAction;
      scheduled_for: string; // ISO-8601 UTC
      status: ScheduledChangeStatus;
      created_by: string;
      executed_at?: string | null;
      cancelled_at?: string | null;
      created_at: string;
      updated_at: string;
    }

    export interface CreateScheduledChangeRequest {
      target_type: ScheduledChangeTargetType;
      target_id: string;
      action: ScheduledChangeAction;
      scheduled_for: string; // ISO-8601 UTC
    }

    export interface UpdateScheduledChangeRequest {
      scheduled_for: string; // ISO-8601 UTC
    }
    ```

- [x] T013 [P] [US1] Create the API service file `frontend/src/services/scheduledChangesApi.ts`. Follow the exact same pattern as `frontend/src/services/changeRequestApi.ts` and `frontend/src/services/killSwitchApi.ts`. Import `apiClient` from `./apiClient`. Import the types from `../types/scheduledChange`. Implement these functions:
    - `createScheduledChange(envId: string, data: CreateScheduledChangeRequest): Promise<ScheduledChange>` — POST to `/environments/${envId}/scheduled-changes`.
    - `listScheduledChanges(envId: string, status?: string): Promise<{ data: ScheduledChange[]; nextPageToken?: string }>` — GET `/environments/${envId}/scheduled-changes` with optional `?status=` query param.
    - `getScheduledChange(id: string): Promise<ScheduledChange>` — GET `/scheduled-changes/${id}`.
    - `updateScheduledChange(id: string, data: UpdateScheduledChangeRequest): Promise<ScheduledChange>` — PATCH `/scheduled-changes/${id}`.
    - `cancelScheduledChange(id: string): Promise<ScheduledChange>` — DELETE `/scheduled-changes/${id}`.
    Each function should use `apiClient.post()`, `apiClient.get()`, `apiClient.patch()`, or `apiClient.delete()` and return `response.data`. Look at `frontend/src/services/killSwitchApi.ts` for the exact pattern of how apiClient is used.

### Frontend — UI Components (US1)

- [x] T014 [US1] Create the badge component `frontend/src/components/flags/ScheduledChangeBadge.tsx`. This component receives a `scheduledChange: ScheduledChange | null` prop. When non-null, it displays a small badge/chip showing: a clock icon (🕐 emoji or an SVG), the action text ("Scheduled: ENABLE" or "Scheduled: DISABLE"), and the formatted local time (convert the UTC `scheduled_for` string to the user's local timezone using `new Date(scheduledChange.scheduled_for).toLocaleString()`). When null, render nothing. Style it with an amber/orange background to indicate pending action, using inline styles or a CSS module. The badge should be small enough to sit inline on a flag row. Add a tooltip (HTML `title` attribute) showing the full UTC time.

- [x] T015 [US1] Create the schedule dialog component `frontend/src/components/flags/ScheduleDialog.tsx`. This component is a modal dialog for creating and cancelling scheduled flag changes. Props: `isOpen: boolean`, `onClose: () => void`, `flagId: string`, `flagName: string`, `environmentId: string`, `existingSchedule: ScheduledChange | null`, `onSuccess: () => void`. The dialog should:
    1. If `existingSchedule` is null (no pending schedule): Show a form with a `<input type="datetime-local">` for selecting the date/time. Add two buttons: "Schedule Enable" and "Schedule Disable". On submit, convert the local datetime to UTC ISO-8601 using `new Date(inputValue).toISOString()`, then call `createScheduledChange(environmentId, { target_type: 'FLAG', target_id: flagId, action: selectedAction, scheduled_for: utcString })`. On success, call `onSuccess()` and `onClose()`. Handle 409 conflict errors by showing an alert "A pending schedule already exists for this flag". Handle other errors generically.
    2. If `existingSchedule` is provided (PENDING schedule exists): Show the current scheduled time (in local timezone), a "Modify" button that reveals a new datetime input, and a "Cancel Schedule" button. The cancel button calls `cancelScheduledChange(existingSchedule.id)` and on success calls `onSuccess()` and `onClose()`. The modify flow calls `updateScheduledChange(existingSchedule.id, { scheduled_for: newUtcString })`.
    Use a simple HTML `<dialog>` element or a div overlay with backdrop, following the project's existing component styling patterns. Import functions from `../services/scheduledChangesApi`.

- [x] T016 [US1] Modify `frontend/src/pages/ProjectDetail.tsx` to integrate the scheduled changes UI into the flag list. This page already displays flag states for each environment. Make these additions:
    1. Import `ScheduledChangeBadge` from `../components/flags/ScheduledChangeBadge` and `ScheduleDialog` from `../components/flags/ScheduleDialog`.
    2. Import `listScheduledChanges` from `../services/scheduledChangesApi` and the `ScheduledChange` type.
    3. Add state: `const [scheduledChanges, setScheduledChanges] = useState<Record<string, ScheduledChange>>({})` to map `flag_id → ScheduledChange`.
    4. Add state: `const [scheduleDialogOpen, setScheduleDialogOpen] = useState(false)` and `const [selectedFlagForSchedule, setSelectedFlagForSchedule] = useState<{id: string, name: string} | null>(null)`.
    5. Add a `useEffect` (or equivalent data fetch) that calls `listScheduledChanges(selectedEnvironmentId, 'PENDING')` when the environment changes, and maps the results into the `scheduledChanges` record keyed by `target_id`.
    6. In the flag list rendering (wherever each flag row/card is rendered), add `<ScheduledChangeBadge scheduledChange={scheduledChanges[flag.id] || null} />` next to the flag name/toggle.
    7. Add a "Schedule" button (🕐 icon or text) on each flag row that sets `selectedFlagForSchedule` and opens the dialog.
    8. Render `<ScheduleDialog isOpen={scheduleDialogOpen} onClose={() => setScheduleDialogOpen(false)} flagId={selectedFlagForSchedule?.id} flagName={selectedFlagForSchedule?.name} environmentId={selectedEnvironmentId} existingSchedule={scheduledChanges[selectedFlagForSchedule?.id] || null} onSuccess={() => { /* re-fetch scheduled changes */ }} />` at the bottom of the component.

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. You can validate using the curl commands in `specs/014-scheduled-flags/quickstart.md` Scenarios 1, 2, 3, and 5.

---

## Phase 4: User Story 2 — Schedule Change Request Application (Priority: P2)

**Goal**: A Release Manager can schedule an approved Change Request to be automatically applied at a future time. The scheduler dispatches the CR application using the existing `ChangeRequestService.Approve` flow.

**Independent Test**: Create an approved CR, schedule it for a near-future time via `POST /api/v1/environments/{envId}/scheduled-changes` with `target_type=CHANGE_REQUEST`, wait, and verify the CR status is `APPLIED`.

### Backend — Scheduler CR Dispatch (US2)

- [x] T017 [US2] Modify `backend/internal/services/scheduler.go` to add the Scheduler a reference to `ChangeRequestService` and implement the CHANGE_REQUEST execution path. Make these changes:
    1. Add a `crService *ChangeRequestService` field to the `Scheduler` struct.
    2. Update the `NewScheduler` constructor to accept `crService *ChangeRequestService` as a parameter and store it.
    3. In the `executeOne` method, replace the `TargetTypeChangeRequest` placeholder log with actual logic:
        - Call `cr, err := s.store.ChangeRequestRepo().GetByID(ctx, sc.TargetID)` to fetch the change request.
        - Verify `cr.Status == models.StatusApproved`. If not approved, log a warning `"skipping scheduled CR application: CR is not in APPROVED state"` and call `s.store.ScheduledChangeRepo().MarkCancelled(ctx, sc.ID, time.Now().UTC())` to auto-cancel the schedule, then return.
        - Reuse the change request application logic: call the existing `ChangeRequestService.Approve` method with the system actor UUID as the reviewer, OR directly apply the changes by extracting the flag state update code from `backend/internal/services/change_request_service.go` lines 104-155 into a helper. The simplest approach: call `s.crService.Approve(ctx, sc.TargetID, models.SystemActorUUID, "Applied by scheduler")`. This re-applies the full approve flow including applying flag state changes and creating an approval record.
        - Note: The existing `Approve` method checks `cr.Status != StatusPending` and returns `ErrChangeRequestNotPending`. For scheduled execution, the CR status is `APPROVED` (not `PENDING`), so you need to add a new method `ApplyScheduled(ctx context.Context, crID uuid.UUID) error` to `ChangeRequestService` that skips the status-pending check and directly applies the proposed changes. Model this method after the existing `Approve` logic (lines 85-155) but enter the transaction starting at the "apply proposed changes" step (line 104), and update status to `APPLIED`.

- [x] T018 [US2] Add the `ApplyScheduled(ctx context.Context, crID uuid.UUID) error` method to `backend/internal/services/change_request_service.go`. This method:
    1. Fetches the CR via `s.store.ChangeRequestRepo().GetByID(ctx, crID)`.
    2. Verifies `cr.Status == models.StatusApproved`. Return `errors.New("change request must be APPROVED to schedule application")` if not.
    3. Opens a transaction via `s.store.WithTx(ctx, func(txStore repository.Store) error { ... })`.
    4. Inside the transaction: (a) Updates CR status to `APPLIED` via `txStore.ChangeRequestRepo().UpdateStatus(ctx, crID, string(models.StatusApplied), nil)`, (b) applies the proposed changes to flag state using the same logic from the existing `Approve` method lines 105-152 (extract `flag_id` and `enabled` from `cr.ProposedChanges`, fetch or create the flag state, update it).
    5. After the transaction: logs an audit entry with `ActorID: models.SystemActorUUID`, `Action: "SCHEDULED_APPLY"`, `TargetType: "CHANGE_REQUEST"`.

- [x] T019 [US2] Update the scheduler wiring in `backend/cmd/server/main.go` to pass `crService` to the `NewScheduler` constructor. Change `services.NewScheduler(store, auditService, cacheClient, logger)` to `services.NewScheduler(store, auditService, crService, cacheClient, logger)`. Make sure the `NewScheduler` function signature in `scheduler.go` matches.

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently. You can validate US2 using the curl commands in `specs/014-scheduled-flags/quickstart.md` Scenario 4.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Tests, edge case handling, error resilience, and documentation that affect both user stories.

- [x] T020 [P] Create unit tests for the scheduled change service in `backend/internal/services/scheduled_change_service_test.go`. Use `testify/assert` and `testify/mock` (or `go-sqlmock`). Test at minimum:
    1. `Create` with a valid future schedule → success.
    2. `Create` with a past `scheduled_for` → error.
    3. `Create` when a pending schedule already exists → `ErrPendingScheduleExists`.
    4. `Create` with invalid action/target_type combo (e.g., target_type=FLAG + action=APPLY) → error.
    5. `Cancel` on a PENDING schedule → success, status becomes CANCELLED.
    6. `Cancel` on an EXECUTED schedule → `ErrScheduleNotPending`.
    7. `UpdateScheduledFor` with valid future time on PENDING → success.
    8. `UpdateScheduledFor` on a CANCELLED schedule → `ErrScheduleNotPending`.
    Follow the test patterns in `backend/internal/services/change_request_service_test.go` and `backend/internal/services/audit_test.go`.

- [x] T021 [P] Create unit tests for the scheduler in `backend/internal/services/scheduler_test.go`. Test:
    1. `executeOne` with target_type=FLAG and action=ENABLE: mock FlagStateRepo.GetByEnvAndFlag to return a state with Enabled=false, verify FlagStateRepo.Update is called with Enabled=true, verify ScheduledChangeRepo.MarkExecuted is called, verify AuditRepo.Create is called with action="SCHEDULED_EXECUTION".
    2. `executeOne` with target_type=FLAG and action=DISABLE: similar but verify Enabled=false.
    3. `executeOne` when FlagStateRepo.GetByEnvAndFlag returns ErrNotFound: verify the error is logged and the method does not panic.
    4. `processDueSchedules` when GetDueSchedules returns empty list: verify no executeOne calls.
    Use `testify/mock` and mock the repository interfaces. The scheduler uses `zerolog` for logging; create a test logger with `zerolog.New(io.Discard)`.

- [x] T022 [P] Add edge-case validation to the Create handler in `backend/internal/api/scheduled_changes.go`: verify the `target_id` actually exists. In the `Create` handler, after parsing the request, add a check: if `target_type == "FLAG"`, call `h.store.FlagRepo().GetByID(ctx, targetID)` and return 404 if the flag doesn't exist. If `target_type == "CHANGE_REQUEST"`, call `h.store.ChangeRequestRepo().GetByID(ctx, targetID)` and return 404 if not found. This prevents scheduling changes for non-existent entities.

- [x] T023 [P] Add a `ScheduledChange` export to `frontend/src/types/index.ts` by adding `export * from './scheduledChange';` at the bottom of the file, so all scheduled change types are available via the barrel export.

- [x] T024 Run `specs/014-scheduled-flags/quickstart.md` validation scenarios 1 through 5 against a running local stack to verify end-to-end functionality. Ensure Docker Compose is up, run the curl commands from quickstart.md, and confirm all expected outcomes match.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 (migration files must exist before model references them conceptually, though Go code doesn't directly depend on the SQL file)
- **User Story 1 (Phase 3)**: Depends on Phase 2 completion — MUST have model, repository, store interface, and DTOs in place
- **User Story 2 (Phase 4)**: Depends on Phase 3 completion — uses the Scheduler created in T009 and the API from T010
- **Polish (Phase 5)**: Can start after Phase 3; T020-T023 are parallelizable; T024 requires Phases 3 and 4

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Phase 2. No dependencies on other stories. This IS the MVP.
- **User Story 2 (P2)**: Depends on the Scheduler from US1 (T009). Adds the CHANGE_REQUEST dispatch path to the existing scheduler.

### Within Each User Story

- Models before services
- Services before API handlers
- API handlers before frontend components
- Backend wiring (main.go) after handlers and services are created
- Frontend types and API client can be built in parallel with backend (marked [P])

### Parallel Opportunities

- **Phase 1**: T001 and T002 can run in parallel (different files)
- **Phase 2**: T003 and T007 can run in parallel; T004-T006 are sequential (they modify the same files/interfaces)
- **Phase 3**: T012 and T013 can run in parallel with backend tasks (frontend vs backend); T014 depends on T012/T013
- **Phase 5**: T020, T021, T022, T023 are all independent and can run in parallel

---

## Parallel Example: User Story 1

```bash
# Backend foundation (sequential):
Task T008: "Create ScheduledChangeService in backend/internal/services/scheduled_change_service.go"
Task T009: "Create Scheduler in backend/internal/services/scheduler.go"
Task T010: "Create API handler in backend/internal/api/scheduled_changes.go"
Task T011: "Wire in backend/cmd/server/main.go"

# Frontend (parallel with backend after T012/T013 have no backend deps):
Task T012: "Create TypeScript types in frontend/src/types/scheduledChange.ts"    [P]
Task T013: "Create API client in frontend/src/services/scheduledChangesApi.ts"   [P]
# Then sequential:
Task T014: "Create ScheduledChangeBadge component"
Task T015: "Create ScheduleDialog component"
Task T016: "Integrate into ProjectDetail.tsx"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (migrations)
2. Complete Phase 2: Foundational (model, repo, store, DTOs)
3. Complete Phase 3: User Story 1 (service, scheduler, API, frontend)
4. **STOP and VALIDATE**: Test with quickstart.md Scenarios 1, 2, 3, 5
5. Deploy/demo if ready — this alone delivers the core scheduling capability

### Incremental Delivery

1. Phase 1 + Phase 2 → Foundation ready
2. Add Phase 3 (US1) → Test independently → Deploy/Demo (MVP! Flag scheduling works)
3. Add Phase 4 (US2) → Test independently → Deploy/Demo (CR scheduling added)
4. Add Phase 5 → Full test suite, edge cases hardened, E2E validated
5. Each phase adds value without breaking previous phases

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- The system sentinel UUID `00000000-0000-0000-0000-000000000001` is used for all scheduler-triggered audit entries
- All timestamps are UTC; the frontend handles local timezone conversion at display time
- The scheduler polls every 30 seconds; worst-case execution latency is ~30s after the scheduled time
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
