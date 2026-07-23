# Quickstart Validation Guide: Data Model & State Management

**Feature**: 002-data-model-state
**Date**: 2026-07-20

---

## Prerequisites

- Docker Compose stack running (`make up`)
- Go 1.22+ installed locally
- `golang-migrate` CLI installed (`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`)

---

## Scenario 1: Migration Up — Schema Provisioning

**Goal**: Verify all migrations execute cleanly against a fresh PostgreSQL instance.

```bash
# From repository root
export DATABASE_URL="postgres://flagmgmt:flagmgmt_dev@localhost:5432/flagmanagment?sslmode=disable"

# Run all migrations up
migrate -path backend/migrations -database "$DATABASE_URL" up

# Verify: list all tables created
docker compose exec postgres psql -U flagmgmt -d flagmanagment -c "\dt"
```

**Expected Result**: 9 tables listed: `projects`, `environments`, `feature_flags`, `environment_flag_states`, `change_requests`, `change_request_approvals`, `audit_logs`, `roles`, `user_roles`, plus `schema_migrations`.

---

## Scenario 2: Migration Down — Clean Rollback

**Goal**: Verify all migrations can be fully rolled back without leaving artifacts.

```bash
# Roll back all migrations
migrate -path backend/migrations -database "$DATABASE_URL" down -all

# Verify: only schema_migrations should remain (or no tables)
docker compose exec postgres psql -U flagmgmt -d flagmanagment -c "\dt"
```

**Expected Result**: No application tables remain. `schema_migrations` table may or may not remain (implementation-dependent).

---

## Scenario 3: Constraint Validation — Foreign Keys & Unique Indexes

**Goal**: Verify database-level integrity constraints work correctly.

```bash
docker compose exec postgres psql -U flagmgmt -d flagmanagment <<'SQL'
-- Insert a project
INSERT INTO projects (id, name, key) VALUES
  ('a0000000-0000-0000-0000-000000000001', 'Test Project', 'test-project');

-- Insert an environment
INSERT INTO environments (id, project_id, name, key, api_key_hash) VALUES
  ('b0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001',
   'Development', 'dev', 'abc123hash');

-- Try duplicate project key (should fail with unique violation)
INSERT INTO projects (id, name, key) VALUES
  ('a0000000-0000-0000-0000-000000000002', 'Duplicate', 'test-project');

-- Try orphaned environment (should fail with FK violation)
INSERT INTO environments (id, project_id, name, key, api_key_hash) VALUES
  ('b0000000-0000-0000-0000-000000000002', 'ffffffff-ffff-ffff-ffff-ffffffffffff',
   'Orphan', 'orphan', 'xyz789hash');
SQL
```

**Expected Result**: First two INSERTs succeed. Third INSERT fails with `duplicate key value violates unique constraint`. Fourth INSERT fails with `foreign key constraint violation`.

---

## Scenario 4: Flag State Evaluation Lookup Performance

**Goal**: Verify the `api_key_hash` → flag states lookup path is indexed and fast.

```bash
docker compose exec postgres psql -U flagmgmt -d flagmanagment <<'SQL'
-- Check index exists
SELECT indexname, indexdef FROM pg_indexes
WHERE tablename = 'environments' AND indexname LIKE '%api_key%';

-- Explain analyze: lookup environment by api_key_hash
EXPLAIN ANALYZE
SELECT id FROM environments WHERE api_key_hash = 'abc123hash';

-- Explain analyze: fetch all flag states for an environment
EXPLAIN ANALYZE
SELECT * FROM environment_flag_states WHERE environment_id = 'b0000000-0000-0000-0000-000000000001';
SQL
```

**Expected Result**: Both queries use index scans (not sequential scans). Execution time < 1ms on empty/small datasets.

---

## Scenario 5: Go Repository Layer Smoke Test

**Goal**: Verify the Go repository layer can perform basic CRUD against PostgreSQL.

```bash
cd backend && go test -race -cover -run TestRepository ./internal/repository/...
```

**Expected Result**: All repository tests pass. Coverage ≥ 80%.

---

## Scenario 6: Audit Log Immutability Verification

**Goal**: Verify audit logs cannot be updated or deleted at the application layer.

```bash
docker compose exec postgres psql -U flagmgmt -d flagmanagment <<'SQL'
-- Insert an audit log entry
INSERT INTO audit_logs (id, actor_id, action, target_type, target_id) VALUES
  ('c0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000001',
   'project.created', 'project', 'a0000000-0000-0000-0000-000000000001');

-- Verify it exists
SELECT id, action FROM audit_logs WHERE id = 'c0000000-0000-0000-0000-000000000001';
SQL
```

**Expected Result**: INSERT succeeds. Application-layer code MUST NOT expose UPDATE or DELETE methods on audit logs (enforced by the repository interface contract — see `contracts/repository-interfaces.md`).
