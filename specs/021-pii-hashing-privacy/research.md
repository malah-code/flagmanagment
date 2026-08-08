# Phase 0: Research & Context

## Hashing Strategy
- **Decision**: MurmurHash3 (32-bit) will be used for high-performance deterministic bucketing (`EvaluateRolloutSplit`), and SHA-256 will be used for persistent PII hashing (`HashPII`).
- **Rationale**: Balances the strict <1ms evaluation performance requirement (Constitution IV) with cryptographic security for persistent audit logs/analytics (Constitution VII).
- **Alternatives considered**: SHA-256 for both (rejected due to slight performance overhead on hot-path bucketing and conflict with Constitution). MurmurHash3 for both (rejected as MurmurHash3 is non-cryptographic and unsafe for persistent PII storage).

## Salting Mechanism
- **Decision**: Each Environment will automatically generate a cryptographically secure 32-byte salt upon creation, stored in the `environments` table.
- **Rationale**: Environment-level salting isolates environments, preventing cross-environment correlation of users, satisfying the Constitution's strict environment isolation rule (II).

## Data Retention Cleanup
- **Decision**: A cron job/scheduler process will run periodically to purge `EvaluationAnalytics` records older than the configured retention period (default 30 days).
- **Rationale**: Required by Constitution VII and feature spec SC-005. Running it out-of-band prevents performance impact on the evaluation API.
