# Audit Logs API Contracts

## 1. Get Audit Logs

Retrieves a paginated list of audit logs, filterable by project, environment, and date range.

**Endpoint**: `GET /api/v1/audit-logs`

**Query Parameters**:
- `project_id` (optional, UUID)
- `environment_id` (optional, UUID)
- `actor_id` (optional, UUID)
- `action_type` (optional, String)
- `start_time` (optional, ISO8601 Timestamp)
- `end_time` (optional, ISO8601 Timestamp)
- `limit` (optional, Int, default 50)
- `offset` (optional, Int, default 0)

**Response** (`200 OK`):
```json
{
  "data": [
    {
      "id": "uuid",
      "actor_id": "uuid",
      "target_project_id": "uuid",
      "target_environment_id": "uuid",
      "target_flag_id": "uuid",
      "action_type": "FLAG_UPDATED",
      "previous_state": { "status": "OFF" },
      "new_state": { "status": "ON" },
      "actor_ip": "192.168.1.1",
      "created_at": "2026-08-08T12:00:00Z"
    }
  ],
  "meta": {
    "total": 1250,
    "limit": 50,
    "offset": 0
  }
}
```

## 2. Export Audit Logs (CSV)

Streams a CSV file containing filtered audit logs.

**Endpoint**: `GET /api/v1/audit-logs/export`

**Query Parameters**: (Same as `Get Audit Logs` above, minus limit/offset)

**Response** (`200 OK`):
- `Content-Type: text/csv`
- `Content-Disposition: attachment; filename="audit_logs_export.csv"`

```csv
id,created_at,actor_id,action_type,target_project_id,target_environment_id,target_flag_id
123e4567-e89b-12d3-a456-426614174000,2026-08-08T12:00:00Z,user-123,FLAG_UPDATED,proj-123,,flag-123
```

## 3. Webhook Delivery Payload

The HTTP POST payload sent to the configured SIEM Webhook URL.

**Method**: `POST`
**Headers**:
- `Content-Type: application/json`
- `X-FlagManagment-Signature`: (HMAC-SHA256 signature if `secret_key` is configured)

**Body**:
```json
{
  "event_id": "uuid-for-webhook-delivery",
  "event_type": "audit.FLAG_UPDATED",
  "timestamp": "2026-08-08T12:00:00Z",
  "data": {
    "id": "uuid-of-audit-log",
    "actor_id": "uuid",
    "target_project_id": "uuid",
    "target_environment_id": "uuid",
    "target_flag_id": "uuid",
    "action_type": "FLAG_UPDATED",
    "previous_state": { "status": "OFF" },
    "new_state": { "status": "ON" },
    "actor_ip": "192.168.1.1",
    "created_at": "2026-08-08T12:00:00Z"
  }
}
```
