# Server-Side Keys API Contracts

## `GET /api/v1/projects/:projectId/environments/:environmentId/server-keys`
List all server keys for a specific environment.

**Response**:
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "billing-service",
      "createdAt": "2026-08-10T12:00:00Z"
    }
  ]
}
```

## `POST /api/v1/projects/:projectId/environments/:environmentId/server-keys`
Create a new server-side environment key.

**Request Body**:
```json
{
  "name": "billing-service"
}
```

**Response**:
```json
{
  "id": "uuid",
  "name": "billing-service",
  "key": "env_abc123...", // ONLY returned once
  "createdAt": "2026-08-10T12:00:00Z"
}
```

## `DELETE /api/v1/projects/:projectId/environments/:environmentId/server-keys/:keyId`
Revoke and delete a server-side environment key.

**Response**: `204 No Content`
