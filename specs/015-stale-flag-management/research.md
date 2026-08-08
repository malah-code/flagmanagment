# Research & Technical Decisions: Stale Flag Detection

**Feature**: `015-stale-flag-management`

## Decision 1: Aggregation of Evaluation Timestamps
- **Decision**: The backend will aggregate flag evaluation timestamps in-memory using an asynchronous buffer (or Redis) and batch-flush updates to PostgreSQL periodically (e.g., every 10 seconds).
- **Rationale**: A synchronous DB update on every flag evaluation (which can be millions per second across SDKs) would instantly saturate the PostgreSQL connection pool and CPU. Batch flushing ensures `last_evaluated_at` is accurate to within a few seconds without impacting read throughput.
- **Alternatives Considered**: 
  - *Direct DB writes*: Rejected due to unacceptable performance impact (violates Constitution Principle IV).
  - *Redis-only storage*: Rejected because historical reporting requires persistent joins with the `feature_flags` table.

## Decision 2: Stale Detection Mechanism
- **Decision**: Implement a background worker (similar to `ScheduledChangeService`) that runs periodically (e.g., once an hour or daily) to scan flag states and evaluation metrics against the environment's `StaleFlagPolicy`.
- **Rationale**: Staleness is a lagging indicator. Real-time detection is unnecessary. A periodic background job efficiently updates lifecycle statuses to `STALE` in bulk without affecting the hot path of flag evaluation.
- **Alternatives Considered**: 
  - *On-read evaluation*: Checking staleness every time the dashboard loads. Rejected as it moves expensive analytical queries to user-facing API latency.

## Decision 3: Flag Lifecycle State Persistence
- **Decision**: Add a `lifecycle_state` ENUM (`ACTIVE`, `STALE`, `DEPRECATED`, `ARCHIVED`) to the `environment_flag_states` table, alongside existing `boolean_state`.
- **Rationale**: A flag's staleness is environment-specific (it might be stale in Dev but still active in Prod). Placing the lifecycle state directly on the environment state junction table ensures the SDK payload generator can easily filter out `ARCHIVED` flags during delta updates.
- **Alternatives Considered**: 
  - *Global flag status*: Rejected because it violates the Environment Isolation principle (Constitution Principle II).
