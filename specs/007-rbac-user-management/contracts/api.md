# API Contracts: RBAC and User Management

## `POST /api/v1/auth/login`
Authenticates a user and returns a JWT session token.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK)**:
```json
{
  "token": "eyJhbGciOiJIUzI1...",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com"
  }
}
```

**Response (401 Unauthorized)**:
```json
{
  "error": "Invalid email or password"
}
```

---

## `GET /api/v1/projects/{project_id}/audit-logs`
Retrieves paginated audit logs for a project (Requires `ADMIN` or `VIEWER` role in the project).

**Headers**:
`Authorization: Bearer <token>`

**Query Parameters**:
- `limit` (default: 50)
- `page_token` (optional)

**Response (200 OK)**:
```json
{
  "logs": [
    {
      "id": "log_uuid",
      "user_id": "user_uuid",
      "action": "FLAG_TOGGLED",
      "resource_type": "FLAG",
      "resource_id": "flag_uuid",
      "old_state": { "enabled": false },
      "new_state": { "enabled": true },
      "created_at": "2026-07-22T19:00:00Z"
    }
  ],
  "next_page_token": "..."
}
```

**Response (403 Forbidden)**:
```json
{
  "error": "Forbidden: You do not have permission to view audit logs for this project."
}
```
