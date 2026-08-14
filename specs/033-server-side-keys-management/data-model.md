# Data Model: Server-Side Keys Management

## Entity: `EnvironmentServerKey`

Represents a named server-side credential used by backend SDKs to authenticate with the FlagManagment engine for a specific environment.

### Database Table: `environment_server_keys`

| Field | Type | Modifiers | Description |
| :--- | :--- | :--- | :--- |
| `id` | UUID | PRIMARY KEY | Unique identifier for the key |
| `environment_id` | UUID | FOREIGN KEY | Links to `environments(id)` ON DELETE CASCADE |
| `name` | VARCHAR(255) | NOT NULL | Human-readable name (e.g. `billing-service`) |
| `key_hash` | VARCHAR(255) | NOT NULL, UNIQUE | SHA-256 hash of the generated API key (for fast auth lookup) |
| `encrypted_key`| VARCHAR(512) | NULLABLE | (Optional) AES-GCM encrypted version of the key to support 'Show/Copy' in UI. If null, key is display-once. |
| `created_at` | TIMESTAMP | NOT NULL | Creation timestamp |
| `last_used_at` | TIMESTAMP | NULLABLE | Timestamp of the last time this key was used to authenticate |

*Note: For the highest security matching our current posture, we can opt to omit `encrypted_key` and only return the plaintext key once upon creation, updating the frontend to show `••••••••` unconditionally with a message "Key hidden for security", aligning with standard AWS/GCP key practices if AES encryption is deemed too complex for this phase.*

### Validations
- `name` MUST be between 1 and 100 characters.
- `name` MUST be unique per `environment_id`.
