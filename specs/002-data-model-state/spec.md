# Feature Specification: Data Model & State Management

**Feature Branch**: `002-data-model-state`

**Created**: 2026-07-20

**Status**: Draft

**Input**: User description: "Design and implement the complete PostgreSQL schema for FlagManagment. Core tables: projects (UUID, Name, Timestamps), environments (UUID, ProjectID FK, Name, APIKeyHash, IsProtected, Timestamps), feature_flags (UUID, ProjectID FK, Key, Type enum, ParentFlagID FK nullable, Timestamps), and environment_flag_states (junction table with EnvironmentID FK, FeatureFlagID FK, BooleanState, TargetingRules JSONB, RemoteConfig JSONB, Timestamps). Governance tables: change_requests, change_request_approvals, audit_logs. RBAC tables: roles, user_roles. Include JSONB schema definitions for contextual targeting rules and remote configuration payloads. Define database migration tooling (golang-migrate) and indexing strategies for high-read performance. Track last_evaluated_at timestamp on every flag for stale flag detection."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Database Schema Provisioning & Migration Pipeline (Priority: P1)

As a DevOps engineer or developer, I want the system to manage its database migrations automatically using versioned scripts so that the PostgreSQL schema can be created, updated, and verified across all environments reliably.

**Why this priority**: Without a version-controlled database schema and migration pipeline, no data can be stored, queried, or updated by the application services.

**Independent Test**: Can be fully tested by running `golang-migrate` up/down commands against a clean PostgreSQL instance and verifying that all tables, foreign keys, indexes, and constraints are created and teardown cleanly.

**Acceptance Scenarios**:

1. **Given** a clean PostgreSQL database, **When** migration scripts are executed up to latest, **Then** all 8 core tables (`projects`, `environments`, `feature_flags`, `environment_flag_states`, `change_requests`, `change_request_approvals`, `audit_logs`, `roles`, `user_roles`) are created with exact data types, constraints, and foreign key relations.
2. **Given** an initialized database, **When** a rollback migration is executed, **Then** schema changes revert cleanly without leaving orphaned tables or indexes.

---

### User Story 2 - Feature Flag State & Targeting Rule Persistence (Priority: P2)

As a product manager or developer, I want to store feature flag definitions, environment-specific states, contextual targeting rules, and remote configuration JSON payloads so that the SDK evaluation engine can retrieve exact configuration rules per environment.

**Why this priority**: Core functionality of a feature flag system relies on retrievingflag definitions, targeting rules (JSONB), and states for a specific project environment.

**Independent Test**: Can be tested by inserting sample projects, environments, flags, and flag states into the database and verifying relational integrity, JSONB rule structure, and constraint validation.

**Acceptance Scenarios**:

1. **Given** an existing project and environment, **When** a new feature flag (Boolean or Multivariate) is created and assigned environment states with JSONB targeting rules, **Then** the flag state record links correctly to environment and feature flag UUIDs and stores valid JSONB.
2. **Given** a parent flag and a child flag, **When** the child flag references the parent flag ID, **Then** the foreign key relationship enforces parent existence, and circular parent references are rejected by application/constraint rules.

---

### User Story 3 - Governance & Audit Trail Data Storage (Priority: P3)

As a security officer or compliance auditor, I want all change requests, approval records, and administrative audit events to be stored in dedicated tables so that full traceability and compliance auditing can be performed.

**Why this priority**: Required by Constitution Principle III (Governance by Default) and Principle VII (PII Protection & Compliance) for enterprise compliance and change tracking.

**Independent Test**: Can be tested by writing change request, approval, and audit log entries to PostgreSQL and verifying foreign key tracking, change delta JSONB fields, and timestamp indexing.

**Acceptance Scenarios**:

1. **Given** a change request created for a protected environment, **When** approvals are submitted and the request is applied, **Then** the `change_requests`, `change_request_approvals`, and `audit_logs` tables record all transition details with exact actor timestamps and state diff JSONB payloads.
2. **Given** an administrative audit log entry, **When** stored in the `audit_logs` table, **Then** user PII attributes (like identity hashes) are stored sanitized/hashed and sensitive tokens/API key plaintexts are omitted.

---

### User Story 4 - RBAC & Access Control Model Persistence (Priority: P4)

As a platform administrator, I want user roles and permission assignments to be stored relationally across global, project, and environment scopes so that access control rules can be evaluated quickly during API requests.

**Why this priority**: Supports granular access enforcement across multiple teams, projects, and environments.

**Independent Test**: Can be tested by creating roles, assigning user-role mappings for specific project/environment scopes, and querying user permissions.

**Acceptance Scenarios**:

1. **Given** defined system roles (System Admin, Project Owner, Release Manager, QA Engineer, Read-Only Auditor), **When** a user is assigned a role for a specific project scope, **Then** `user_roles` accurately links user ID, role ID, and project/environment scope UUIDs.

---

### Edge Cases

- **Stale Flag Evaluation Timestamp**: What happens when a flag has never been evaluated by an SDK? The `last_evaluated_at` column MUST be NULL initially and updated asynchronously without blocking read evaluations.
- **JSONB Schema Validation Failure**: What happens if invalid JSON is written to `targeting_rules` or `remote_config`? Database CHECK constraints or application-level JSON schema validators MUST reject invalid JSON payloads before persistence.
- **Environment Deletion Cascade**: What happens when an environment is deleted? Cascade rules MUST clean up associated `environment_flag_states` while preserving `audit_logs` for historical compliance.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST represent entity IDs using standard 128-bit UUID (v4/v7) primary keys across all tables for global uniqueness and distributed safety.
- **FR-002**: System MUST implement migration scripts using `golang-migrate` format (`.up.sql` and `.down.sql`) under `backend/migrations/`.
- **FR-003**: System MUST store `projects` containing `id` (UUID), `name` (string), `key` (unique string), `description` (text), `created_at`, and `updated_at`.
- **FR-004**: System MUST store `environments` linked to a `project_id` (FK) with `id` (UUID), `name` (string), `key` (string), `api_key_hash` (string, unique), `is_protected` (boolean), `created_at`, and `updated_at`. Key MUST be unique per project.
- **FR-005**: System MUST store `feature_flags` linked to a `project_id` (FK) with `id` (UUID), `key` (string, unique per project), `name` (string), `description` (text), `type` (enum: `BOOLEAN`, `MULTIVARIATE`, `JSON`), `parent_flag_id` (FK nullable to `feature_flags.id`), `last_evaluated_at` (timestamp nullable), `created_at`, and `updated_at`.
- **FR-006**: System MUST store `environment_flag_states` joining `environment_id` (FK) and `feature_flag_id` (FK) with `id` (UUID), `enabled` (boolean), `targeting_rules` (JSONB), `remote_config` (JSONB), `variations` (JSONB nullable for multivariate), `created_at`, and `updated_at`. Primary lookup constraint MUST be unique on `(environment_id, feature_flag_id)`.
- **FR-007**: System MUST store `change_requests` with `id` (UUID), `project_id` (FK), `environment_id` (FK), `title` (string), `description` (text), `status` (enum: `PENDING`, `APPROVED`, `REJECTED`, `APPLIED`), `proposed_changes` (JSONB diff), `created_by` (UUID), `applied_by` (UUID nullable), `created_at`, and `updated_at`.
- **FR-008**: System MUST store `change_request_approvals` with `id` (UUID), `change_request_id` (FK), `approver_id` (UUID), `decision` (enum: `APPROVE`, `REJECT`), `comment` (text), and `created_at`.
- **FR-009**: System MUST store `audit_logs` with `id` (UUID), `project_id` (FK nullable), `environment_id` (FK nullable), `actor_id` (UUID), `action` (string), `target_type` (string), `target_id` (UUID), `previous_state` (JSONB nullable), `new_state` (JSONB nullable), `actor_ip` (string), and `created_at`.
- **FR-010**: System MUST store `roles` with `id` (UUID), `name` (string, unique), `description` (text), `permissions` (JSONB / text array), `created_at`, and `updated_at`.
- **FR-011**: System MUST store `user_roles` with `id` (UUID), `user_id` (UUID), `role_id` (FK), `project_id` (FK nullable), `environment_id` (FK nullable), `created_at`, and `updated_at`.
- **FR-012**: System MUST define high-performance secondary indexes for high-read evaluation paths:
  - `idx_env_flag_states_env_id`: Index on `environment_flag_states(environment_id)`
  - `idx_environments_api_key_hash`: Unique index on `environments(api_key_hash)`
  - `idx_feature_flags_project_key`: Unique index on `feature_flags(project_id, key)`
  - `idx_audit_logs_project_env_created`: Compound index on `audit_logs(project_id, environment_id, created_at DESC)`
- **FR-013**: System MUST support tracking `last_evaluated_at` on `feature_flags` without causing lock contention during high-throughput evaluations.
- **FR-014**: System MUST define JSONB schemas for `targeting_rules` supporting attribute matching (operators: `EQUALS`, `NOT_EQUALS`, `CONTAINS`, `REGEX`, `IN`, `NOT_IN`), percentage rollouts, and segment IDs.
- **FR-015**: System MUST enforce foreign key cascade deletes for dependent flag states when an environment or flag is deleted, while ensuring audit logs and change requests preserve referential integrity without cascading deletion of historical compliance data.

### Key Entities *(include if feature involves data)*

- **Project**: Represents a top-level organization workspace containing multiple environments and feature flags.
- **Environment**: Represents a deployment target (e.g., Dev, QA, Staging, Production) with its own SDK authentication key hash and protection settings.
- **FeatureFlag**: Represents a feature toggle definition (Boolean, Multivariate, or JSON payload) with an optional parent flag dependency.
- **EnvironmentFlagState**: Represents the active evaluation state, targeting rules (JSONB), and variant payload for a specific flag within a specific environment.
- **ChangeRequest**: Represents a proposed set of flag state mutations requiring review and approval before being applied to a protected environment.
- **ChangeRequestApproval**: Represents an individual approval or rejection decision logged by a reviewer on a change request.
- **AuditLog**: Represents an immutable, append-only record of administrative mutations and governance events.
- **Role & UserRole**: Represents RBAC permissions and assignments mapped to users across global, project, or environment scopes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Database schema migrations (`up` and `down`) complete execution in under 5 seconds on standard database instances.
- **SC-002**: Database lookup for all environment flag states by `api_key_hash` executes in under 2 milliseconds for environments containing 500+ flags.
- **SC-003**: Schema passes 100% of automated structural integrity and constraint validation tests (foreign keys, nullability, unique indexes).
- **SC-004**: Database supports concurrent read evaluation workloads without deadlock or table locking under 1,000 requests per second.

## Assumptions

- PostgreSQL 16+ is used as the primary relational datastore per Constitution constraints.
- `golang-migrate` CLI / Go library is used for versioned SQL file execution.
- Sensitive credentials (API keys) are hashed using SHA-256 before storing in `environments.api_key_hash`.
- User authentication IDs (`user_id`, `actor_id`) are formatted as UUIDs to integrate with external OIDC/OAuth identity providers.
