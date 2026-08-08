# Data Model: PII Hashing & Data Privacy

## Entity Modifications

### Environment
The `Environment` entity will be updated to include a `Salt` property to enable environment-scoped deterministic hashing.

**New Fields:**
- `Salt` (string/varchar): A 64-character hex string representing 32 cryptographically random bytes generated upon environment creation.

### AuditLog / EvaluationAnalytics
Existing tables that store analytics/telemetry (e.g. `audit_logs` or evaluation events) do not need schema changes, but the system must ensure that attributes containing PII (like `email`, `user_id`, `phone`) are replaced with their SHA-256 salted hash equivalent before being written to these tables.

### TargetingRule
Rules containing exact matches for PII strings (e.g., "email equals user@example.com") must hash the literal value (`user@example.com`) in the rule definition before saving it to the database, ensuring no plaintext PII exists in `flags` JSONB payload.

## State Transitions
N/A - Hashing occurs immutably on write/evaluation.
