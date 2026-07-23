# Contract: Health Check Endpoint

**Feature**: `001-core-architecture`
**Date**: 2026-07-18

---

## Overview

The health check endpoint is the only API contract defined in this foundation
phase. It serves two purposes:

1. **Docker Compose orchestration** — Used in `depends_on` health checks to
   sequence service startup.
2. **Dashboard integration** — The minimal landing page fetches this endpoint
   to display service connectivity status.

---

## Endpoint

### `GET /healthz`

Returns the composite health status of the backend and its dependencies.

**Request**: No body, no query parameters, no authentication.

**Response** (200 OK — all dependencies healthy):

```json
{
  "status": "healthy",
  "version": "0.1.0",
  "uptime_seconds": 3421,
  "checks": {
    "postgres": {
      "status": "healthy",
      "latency_ms": 2
    },
    "redis": {
      "status": "healthy",
      "latency_ms": 1
    }
  }
}
```

**Response** (503 Service Unavailable — one or more dependencies unhealthy):

```json
{
  "status": "unhealthy",
  "version": "0.1.0",
  "uptime_seconds": 3421,
  "checks": {
    "postgres": {
      "status": "healthy",
      "latency_ms": 2
    },
    "redis": {
      "status": "unhealthy",
      "error": "connection refused"
    }
  }
}
```

---

## Response Schema

| Field | Type | Description |
|-------|------|-------------|
| `status` | Enum: `healthy`, `unhealthy` | Composite status — `unhealthy` if any check fails |
| `version` | String | Application version (semver) |
| `uptime_seconds` | Integer | Seconds since the backend process started |
| `checks` | Object | Map of dependency name → check result |
| `checks.<name>.status` | Enum: `healthy`, `unhealthy` | Individual dependency status |
| `checks.<name>.latency_ms` | Integer (optional) | Round-trip time for the health check probe |
| `checks.<name>.error` | String (optional) | Error message, present only when unhealthy |

---

## Behavior Rules

- The endpoint MUST respond within **1 second** (SC-004).
- If PostgreSQL is unreachable, the backend MUST return 503 with `postgres.status = "unhealthy"` — it MUST NOT crash or panic.
- If Redis is unreachable, the backend MUST return 503 with `redis.status = "unhealthy"` — it MUST NOT crash or panic.
- The endpoint MUST NOT require authentication (it is used for infrastructure health probing).
- The `version` field MUST match the version embedded at build time (via `ldflags`).

---

## Usage by Docker Compose

```yaml
backend:
  healthcheck:
    test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/healthz"]
    interval: 10s
    timeout: 5s
    retries: 3
```

---

## Usage by Frontend Dashboard

The minimal landing page fetches `GET /healthz` every 10 seconds and displays:

- Overall status badge (green/red)
- Individual dependency status (PostgreSQL, Redis)
- Backend version and uptime

No additional API contracts are defined in this foundation phase. Business APIs
(projects, environments, flags) are covered in Feature 3 (API Contracts &
Streaming Protocol).
