# API Contracts: User Management & Config

## Users

### `GET /api/v1/users`
Returns a list of users and their access.
**Response**:
```json
{
  "users": [
    {
      "id": "uuid",
      "email": "user@example.com",
      "name": "Jane",
      "role": "global_admin",
      "status": "active",
      "projects": []
    }
  ]
}
```

### `POST /api/v1/users/invite`
Invites a new user.
**Request**:
```json
{
  "email": "new@example.com",
  "role": "project_editor",
  "project_ids": ["uuid-1", "uuid-2"]
}
```

### `PUT /api/v1/users/:id/access`
Updates user roles/projects.
**Request**:
```json
{
  "role": "read_only_auditor",
  "project_ids": ["uuid-1"]
}
```

## System Config

### `GET /api/v1/config/smtp`
Retrieves SMTP settings (password redacted).
**Response**:
```json
{
  "host": "smtp.example.com",
  "port": 587,
  "username": "apikey",
  "sender_email": "noreply@example.com"
}
```

### `PUT /api/v1/config/smtp`
Updates SMTP settings.
**Request**:
```json
{
  "host": "smtp.example.com",
  "port": 587,
  "username": "apikey",
  "password": "supersecretpassword",
  "sender_email": "noreply@example.com"
}
```

### `POST /api/v1/config/smtp/test`
Dispatches a test email.
**Request**: 
```json
{ 
  "to": "admin@example.com" 
}
```
