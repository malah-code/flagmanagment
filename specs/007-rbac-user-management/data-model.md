# Data Model: Granular RBAC and User Management

## Entities

### `users`
Represents an individual who can log in to the dashboard.
- `id` (UUID, Primary Key)
- `email` (String, Unique, Not Null)
- `password_hash` (String, Not Null) - Bcrypt hash of the user's password
- `created_at` (Timestamp)
- `updated_at` (Timestamp)

### `project_role_assignments`
Maps a user to a specific role within a project.
- `id` (UUID, Primary Key)
- `user_id` (UUID, Foreign Key -> users.id)
- `project_id` (UUID, Foreign Key -> projects.id)
- `role` (Enum: `VIEWER`, `EDITOR`, `ADMIN`)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)
*Constraint: Unique index on `(user_id, project_id)`*

### `audit_logs`
Immutable record of all state-mutating actions (append-only).
- `id` (UUID, Primary Key)
- `user_id` (UUID, Foreign Key -> users.id, Nullable if system action)
- `project_id` (UUID, Foreign Key -> projects.id, Nullable)
- `action` (String) - e.g., `FLAG_CREATED`, `FLAG_TOGGLED`, `ENVIRONMENT_CREATED`
- `resource_type` (String) - e.g., `FLAG`, `ENVIRONMENT`, `PROJECT`
- `resource_id` (UUID) - ID of the mutated resource
- `old_state` (JSONB, Nullable) - State before mutation
- `new_state` (JSONB, Nullable) - State after mutation (sanitized of PII)
- `created_at` (Timestamp) - When the action occurred

## State Transitions
* **Roles**: Users can be granted, updated, or revoked roles by Project Admins.
* **Audit Logs**: Insert-only. No updates or deletes permitted.
