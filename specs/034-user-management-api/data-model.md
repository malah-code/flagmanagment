# Data Model: User Management & Config

## Entities

### User
- `id` (UUID, Primary Key)
- `email` (String, Unique, Indexed)
- `name` (String)
- `role` (Enum: `global_admin`, `project_editor`, `read_only_auditor`)
- `status` (Enum: `active`, `pending`)
- `created_at` (Timestamp)
- `last_active_at` (Timestamp)

### ProjectAccess
- `id` (UUID, Primary Key)
- `user_id` (UUID, Foreign Key)
- `project_id` (UUID, Foreign Key)
- `role` (Enum: `project_editor`, `read_only_auditor`)

### Invitation
- `id` (UUID, Primary Key)
- `email` (String)
- `token_hash` (String)
- `role` (String)
- `project_ids` (JSONB Array)
- `expires_at` (Timestamp)
- `created_by` (UUID, Foreign Key -> User)

### SystemConfig
- `key` (String, Primary Key) e.g., `smtp_config`
- `value` (JSONB)
- `updated_at` (Timestamp)

*Note: The `value` for `smtp_config` will store `{ host, port, username, encrypted_password, sender_email }`.*
