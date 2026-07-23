# Research: Granular RBAC and User Management

## Authentication & Password Hashing
- **Decision**: Use `bcrypt` for password hashing in Go and JWT for session tokens.
- **Rationale**: `bcrypt` is the industry standard for securely hashing passwords with built-in salting. JWT (JSON Web Tokens) enables stateless authentication across the frontend dashboard and backend REST APIs without needing session lookup in Redis (though Redis can be used for token blacklisting if needed).
- **Alternatives considered**: `argon2` (slightly more secure, but `bcrypt` is simpler and fully sufficient for MVP). Stateful cookie sessions (less flexible for API consumers).

## Role-Based Access Control (RBAC) Architecture
- **Decision**: Implement a `project_role_assignments` table in PostgreSQL mapping `user_id`, `project_id`, and `role` (enum: Viewer, Editor, Admin).
- **Rationale**: The spec mandates project-level roles. A relational table easily handles this many-to-many mapping. The backend API middleware will intercept requests, look up the user's role for the requested `project_id`, and compare it against the required role for the endpoint (e.g., `POST /projects/{id}/environments` requires `Editor` or `Admin`).
- **Alternatives considered**: Global roles only (rejected per spec). Casbin/OPA for complex policy evaluation (overkill for fixed Viewer/Editor/Admin roles).

## Immutable Audit Logging
- **Decision**: Create an `audit_logs` PostgreSQL table (append-only) with JSONB columns for `old_state` and `new_state`.
- **Rationale**: PostgreSQL JSONB is perfect for storing unstructured historical data. Triggers or application-level logging can insert records asynchronously.
- **Alternatives considered**: Event Sourcing (too complex for MVP). NoSQL datastore (violates Constitution constraint to use PostgreSQL 16+ as primary datastore).

## PII Sanitization in Audit Logs
- **Decision**: Application-level scrubbing of sensitive keys (e.g., `api_key`) before inserting into the `audit_logs` table.
- **Rationale**: The spec requires that no plain-text API keys appear in the payload. Intercepting the payload in the Go service layer and redacting fields (replacing with `***`) before saving ensures compliance.
- **Alternatives considered**: Database-level redaction (harder to maintain across schema changes).
