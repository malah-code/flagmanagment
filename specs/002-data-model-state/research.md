# Research: Data Model & State Management

**Feature**: 002-data-model-state
**Date**: 2026-07-20

---

## §1 Migration Tool: golang-migrate

**Decision**: Use `golang-migrate/migrate` v4 (Go library + CLI).

**Rationale**:
- Constitution mandates `golang-migrate or equivalent versioned migration tool` (Technology Stack Constraints).
- golang-migrate is the most widely adopted Go migration library with 14k+ GitHub stars.
- Supports both CLI execution (`migrate -path ./migrations -database $DB_URL up`) and programmatic use via Go's `migrate.New()`.
- Supports PostgreSQL natively via `pgx` or `lib/pq` driver.
- File-based migrations use numbered `.up.sql` / `.down.sql` files — simple, auditable, version-controlled.

**Alternatives Considered**:
- **goose** — Similar feature set but less community adoption. `golang-migrate` has better Docker integration.
- **atlas** — Modern declarative approach but adds complexity beyond what Phase 1 needs. Consider for future schema drift detection.
- **sqlx-migrate** — Tied to the `sqlx` library which we are not using.

**Implementation Notes**:
- Migration files live in `backend/migrations/`.
- File naming: `NNNNNN_description.up.sql` and `NNNNNN_description.down.sql` (6-digit zero-padded sequential).
- Migrations run automatically at application startup via embedded Go library OR via CLI in CI/CD.
- For development: run via `make migrate-up` / `make migrate-down` targets.
- The `migrate` CLI is installed inside the dev Docker container for interactive use.

---

## §2 UUID Strategy

**Decision**: Use UUID v4 (random) for all entity primary keys. Generate in Go using `google/uuid`.

**Rationale**:
- UUID v4 provides global uniqueness without coordination — essential for distributed or multi-tenant deployments.
- PostgreSQL's native `UUID` type stores UUIDs as 128-bit values (16 bytes) — efficient storage and indexing.
- `google/uuid` is the de facto Go UUID library (well-maintained, zero dependencies beyond standard library).

**Alternatives Considered**:
- **UUID v7 (timestamp-ordered)** — Better for B-tree index locality. However, the Go ecosystem support is still maturing and `google/uuid` v1.6+ supports it. Can be adopted later if index performance becomes an issue.
- **ULID** — Sortable, but adds a non-standard dependency and PostgreSQL doesn't have a native ULID type.
- **Serial/BIGSERIAL** — Not suitable for distributed systems; reveals record count; not globally unique.

**Implementation Notes**:
- Go: `uuid.New()` generates UUID v4. Store as `uuid.UUID` in Go structs.
- PostgreSQL: Column type `UUID` with `DEFAULT gen_random_uuid()` as fallback.
- All foreign keys reference `UUID` columns.

---

## §3 JSONB Schema for Targeting Rules

**Decision**: Use a structured JSONB schema for `targeting_rules` with documented operators and nesting.

**Rationale**:
- JSONB provides schema flexibility for targeting rules that evolve independently of the relational schema.
- PostgreSQL JSONB supports indexing via GIN indexes for containment queries.
- Validation happens at the application layer (Go) before persistence.

**Schema Definition**:

```json
{
  "rules": [
    {
      "id": "rule-uuid",
      "priority": 1,
      "conditions": [
        {
          "attribute": "user.country",
          "operator": "IN",
          "values": ["US", "CA", "UK"]
        },
        {
          "attribute": "user.plan",
          "operator": "EQUALS",
          "value": "enterprise"
        }
      ],
      "logic": "AND",
      "variation_id": "variation-uuid",
      "percentage": null
    },
    {
      "id": "rule-uuid-2",
      "priority": 2,
      "conditions": [],
      "logic": "AND",
      "variation_id": null,
      "percentage": {
        "buckets": [
          { "variation_id": "var-a", "weight": 50 },
          { "variation_id": "var-b", "weight": 50 }
        ],
        "seed": "flag-key-hash-seed"
      }
    }
  ],
  "default_variation_id": "variation-uuid-default"
}
```

**Supported Operators**: `EQUALS`, `NOT_EQUALS`, `CONTAINS`, `NOT_CONTAINS`, `REGEX`, `IN`, `NOT_IN`, `GREATER_THAN`, `LESS_THAN`, `SEMVER_EQUALS`, `SEMVER_GREATER_THAN`, `SEMVER_LESS_THAN`.

---

## §4 JSONB Schema for Remote Configuration

**Decision**: Use a flexible JSONB schema for `remote_config` that stores arbitrary key-value configuration payloads.

**Rationale**:
- Remote configuration payloads are inherently schemaless — different flags serve different config structures (color themes, feature limits, A/B test parameters).
- JSONB provides efficient storage and querying without requiring schema migration for every new config shape.

**Schema Definition**:

```json
{
  "type": "JSON",
  "value": {
    "maxRetries": 3,
    "timeout_ms": 5000,
    "featureEnabled": true,
    "uiConfig": {
      "theme": "dark",
      "layout": "grid"
    }
  },
  "schema_version": "1.0"
}
```

---

## §5 Indexing Strategy

**Decision**: Create targeted B-tree and unique indexes for high-read evaluation paths.

**Rationale**:
- SDK flag evaluation is the hottest read path — `api_key_hash` lookup into `environment_flag_states` by `environment_id`.
- Audit log queries filter by `project_id`, `environment_id`, and sort by `created_at DESC`.
- Unique indexes enforce business rules at the database level (no duplicate flag keys per project, no duplicate API keys).

**Index Plan**:

| Index Name | Table | Columns | Type | Purpose |
|------------|-------|---------|------|---------|
| `idx_environments_api_key_hash` | `environments` | `api_key_hash` | UNIQUE B-tree | SDK authentication lookup |
| `idx_environments_project_key` | `environments` | `(project_id, key)` | UNIQUE B-tree | Unique env key per project |
| `idx_feature_flags_project_key` | `feature_flags` | `(project_id, key)` | UNIQUE B-tree | Unique flag key per project |
| `idx_env_flag_states_lookup` | `environment_flag_states` | `(environment_id, feature_flag_id)` | UNIQUE B-tree | Primary evaluation lookup |
| `idx_env_flag_states_env_id` | `environment_flag_states` | `environment_id` | B-tree | Bulk flag state fetch per env |
| `idx_audit_logs_query` | `audit_logs` | `(project_id, environment_id, created_at DESC)` | B-tree | Audit log dashboard queries |
| `idx_change_requests_env` | `change_requests` | `(environment_id, status)` | B-tree | Pending CR lookup |
| `idx_user_roles_user` | `user_roles` | `(user_id, project_id)` | B-tree | RBAC permission check |

---

## §6 API Key Hashing

**Decision**: SHA-256 hash of the raw API key stored in `environments.api_key_hash`.

**Rationale**:
- API keys are high-entropy random tokens (not passwords) — SHA-256 is sufficient without salting for lookup speed.
- Constitution Principle VII mandates that plaintext API keys are never stored.
- SHA-256 is deterministic — enables O(1) hash lookup without iterating salts (unlike bcrypt which is intentionally slow).

**Implementation Notes**:
- Go: `crypto/sha256` → `hex.EncodeToString(hash[:])`.
- SDK connects with raw key → backend hashes it → looks up `api_key_hash` in database.
- API key is shown to user exactly once at creation time, then discarded.

---

## §7 Stale Flag Detection

**Decision**: `last_evaluated_at` column on `feature_flags` updated asynchronously via background batch writer.

**Rationale**:
- Constitution Principle IV requires sub-millisecond evaluation with zero network calls. Synchronous `UPDATE` on every evaluation would add write contention.
- Batch writer collects evaluation timestamps in memory and flushes to PostgreSQL every N seconds (configurable, default 60s).
- Dashboard queries `WHERE last_evaluated_at < NOW() - INTERVAL '30 days' AND enabled = true` to find stale flags.

---

## §8 Repository Pattern

**Decision**: Use a thin repository layer (`internal/repository/`) wrapping `database/sql` directly.

**Rationale**:
- The repository pattern provides a clean boundary between business logic and database access.
- Using `database/sql` directly (no ORM) keeps the codebase simple, explicit, and performant — aligned with Go idioms.
- Repositories are interface-driven for easy mocking in unit tests.

**Alternatives Considered**:
- **sqlx** — Adds struct scanning convenience but introduces a dependency. Can adopt later if boilerplate becomes excessive.
- **GORM / ent** — Full ORMs add complexity and hide SQL. Not appropriate for a system that needs fine-grained query control.
- **sqlc** — Generates Go code from SQL. Excellent choice but adds a code generation step. Consider for Feature 003 when API queries become complex.
