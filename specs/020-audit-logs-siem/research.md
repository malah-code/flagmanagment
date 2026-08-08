# Phase 0: Research & Technical Decisions

## 1. Webhook Dispatching & Retry Logic

**Decision**: Use an in-memory buffered channel with worker goroutines for webhook dispatching, falling back to a lightweight in-memory retry mechanism using Go's `time.AfterFunc` for exponential backoff, rather than introducing a heavy external task queue dependency like Celery or Asynq immediately.

**Rationale**: 
The FlagManagment backend is built in Go and already utilizes goroutines heavily. Since webhook payloads (audit log JSON) are small, we can afford to keep them in memory for the short duration of the retry window (e.g., 3 retries over a few minutes). Introducing a persistent queuing system (like RabbitMQ or Redis-backed Asynq) specifically for webhooks adds infrastructure complexity. If a webhook delivery completely fails after 3 retries or the server crashes, the loss of that real-time webhook is acceptable because the source of truth (the Immutable Audit Log in PostgreSQL) is always safely written *before* the webhook is dispatched. The SIEM can always query the REST API to reconcile missed events.

**Alternatives considered**: 
- **Redis-backed Asynq queue**: More robust for guaranteeing delivery across server restarts, but adds more moving parts to the architecture. Can be adopted later if webhook volume scales massively.
- **Synchronous dispatch**: Unacceptable as it would block the primary administrative API request waiting for the SIEM endpoint to respond.

## 2. PII & Secret Sanitization Strategy

**Decision**: Implement a generic JSON scrubber that traverses the `PreviousState` and `NewState` JSON objects right before database insertion, recursively replacing the values of known sensitive keys (e.g., `api_key`, `token`, `password`, `email`) with `[REDACTED]`.

**Rationale**:
Targeting rules and flag states are stored as flexible `JSONB` in PostgreSQL. We cannot hardcode structs for every possible payload shape. A recursive map traversal in Go is fast enough for the infrequent nature of administrative mutations. 

**Alternatives considered**:
- **Sanitizing on read**: Increases read latency for audit logs and risks leaking secrets if the database is ever compromised. Sanitizing on write ensures the database itself is clean.

## 3. CSV Export Handling

**Decision**: Implement the CSV export endpoint using Go's `encoding/csv` writer paired with an `http.ResponseWriter` stream.

**Rationale**:
Streaming the CSV directly from the database cursor to the HTTP response prevents memory exhaustion on the backend when an admin exports a large date range of audit logs.

**Alternatives considered**:
- **Generating the CSV in-memory and returning it**: Will cause out-of-memory (OOM) crashes on large exports.
- **Asynchronous email delivery**: Overkill for an administrative dashboard. Streaming download is standard UX.
