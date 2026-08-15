# Phase 0 Research: Flag Creation & Targeting Workflows

## Decisions & Architecture

### 1. Multi-Type Flag Representations (`BOOLEAN`, `MULTIVARIATE`, `JSON`)
- **Decision**: Represent variations as a structured JSONB column in PostgreSQL (`variations` array on `feature_flags` table) containing unique string variation keys/IDs, names, and concrete type values.
- **Rationale**: Provides schema flexibility for remote config JSON payloads while maintaining strong typing and high-speed in-memory evaluation for boolean and multivariate flags.
- **Alternatives Considered**: Dedicated SQL relational tables for variations (rejected due to excessive joins in evaluation hot path).

### 2. Contextual Targeting & Rule Evaluation Syntax
- **Decision**: Adopt OpenFeature standard rule syntax with condition objects `{ attribute, operator, values }` evaluated against user attributes.
- **Operators**: `EQUALS`, `NOT_EQUALS`, `CONTAINS`, `STARTS_WITH`, `ENDS_WITH`, `IN_LIST`, `GREATER_THAN`, `LESS_THAN`, `SEMVER_EQ`, `SEMVER_GTE`.
- **Rationale**: Enables seamless SDK evaluation and standard OpenFeature provider interoperability.

### 3. Percentage Rollouts and Deterministic Bucketing
- **Decision**: MurmurHash3 algorithm hashing `flagKey + ":" + environmentSalt + ":" + userContextKey` to map into a `[0, 99]` bucket range.
- **Rationale**: Ensures deterministic, sticky variation assignment across client and server evaluations without storing user assignment states in the database.

### 4. Kill Switch vs Standard Toggle
- **Decision**: Kill switch sets `is_enabled = false` and flags the environment flag state with an emergency metadata block `{ killed: true, killed_by: ..., reason: ... }`.
- **Rationale**: Differentiates emergency operational shutoffs from routine manual toggles in audit trails and monitoring dashboards.

### 5. Automated Scheduled Changes
- **Decision**: Background worker (`SchedulerService`) polling PostgreSQL for pending state transitions where `scheduled_at <= NOW()` and `status = 'PENDING'`, executed within atomic transactions with audit logs.
- **Rationale**: Simple, highly reliable, and eliminates dependency on external distributed job queues for local and self-hosted deployments.
