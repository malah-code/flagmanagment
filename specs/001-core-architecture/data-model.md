# Data Model: FlagManagment Core Architecture & Repository Bootstrap

**Feature**: `001-core-architecture`
**Date**: 2026-07-18

---

> **Note**: This is the foundation phase. The entities defined here represent the
> **infrastructure skeleton** — services, configuration, and build targets. The
> domain data model (projects, environments, flags, RBAC) is covered in Feature 2
> (Data Model & State Management).

## Entities

### Service

A deployable unit of the FlagManagment platform. Each service runs as an independent
container in the Docker Compose orchestration.

| Field | Type | Description |
|-------|------|-------------|
| Name | String | Unique identifier (e.g., `backend`, `frontend`, `postgres`, `redis`) |
| Image | String | Container image reference (e.g., `ghcr.io/flagmanagment/backend:latest`) |
| Ports | List<PortMapping> | Exposed host:container port pairs |
| HealthCheck | HealthCheckConfig | Endpoint or command to verify service readiness |
| DependsOn | List<ServiceRef> | Services that must be healthy before this service starts |
| Environment | Map<String, String> | Environment variables injected at runtime |
| Volumes | List<VolumeMount> | Persistent or bind-mounted filesystem paths |

**Relationships**:
- A Service has zero or more DependsOn references to other Services
- The Backend service depends on PostgreSQL and Redis
- The Frontend service depends on the Backend service (for health API)

**Validation Rules**:
- Name MUST be unique across all services in the compose file
- At least one HealthCheck MUST be defined per service
- Circular dependencies MUST NOT exist in the DependsOn graph

---

### PortMapping

Maps a host port to a container port.

| Field | Type | Description |
|-------|------|-------------|
| Host | Integer | Port exposed on the host machine |
| Container | Integer | Port inside the container |
| Protocol | Enum (tcp, udp) | Network protocol (default: tcp) |

**Validation Rules**:
- Host port MUST be unique across all services (no port conflicts)
- Host port MUST be configurable via environment variable

---

### HealthCheckConfig

Defines how to verify a service is ready to accept traffic.

| Field | Type | Description |
|-------|------|-------------|
| Type | Enum (http, cmd) | HTTP endpoint check or shell command check |
| Target | String | URL path (for http) or command string (for cmd) |
| Interval | Duration | Time between check attempts (default: 10s) |
| Timeout | Duration | Maximum time to wait for a response (default: 5s) |
| Retries | Integer | Number of consecutive failures before marking unhealthy (default: 3) |

---

### EnvironmentConfiguration

A set of environment variables that controls runtime behavior.

| Field | Type | Description |
|-------|------|-------------|
| VariableName | String | Environment variable name (prefixed with `FM_`) |
| DefaultValue | String | Default value used when not explicitly set |
| Description | String | Human-readable purpose of this variable |
| Service | ServiceRef | Which service(s) consume this variable |
| Sensitive | Boolean | Whether this value should be masked in logs (default: false) |

**Validation Rules**:
- All variable names MUST use the `FM_` prefix for namespacing
- Sensitive values (passwords, keys) MUST NOT appear in structured log output
- Every variable MUST have a documented default in `.env.example`

---

### BuildTarget

Defines a target platform for container image builds.

| Field | Type | Description |
|-------|------|-------------|
| OS | String | Operating system (e.g., `linux`) |
| Architecture | String | CPU architecture (e.g., `amd64`, `arm64`) |
| Variant | String (optional) | Architecture variant (e.g., `v8` for ARMv8) |
| PlatformString | String | Docker platform identifier (e.g., `linux/amd64`) |

**Supported Targets**:

| PlatformString | Use Case |
|----------------|----------|
| `linux/amd64` | Servers, Linux workstations, x86_64 CI runners |
| `linux/arm64` | Apple Silicon Macs, Raspberry Pi 4, ARM servers |

---

## Entity Relationships

```text
Service ──has──> PortMapping (1:N)
Service ──has──> HealthCheckConfig (1:1)
Service ──depends_on──> Service (N:M, acyclic)
Service ──consumes──> EnvironmentConfiguration (N:M)
BuildTarget ──produces──> Service image (N:M)
```

---

## State Transitions

### Service Lifecycle (Docker Compose)

```text
Created → Starting → Waiting (depends_on health checks)
                          ↓
                     Healthy ←→ Unhealthy
                          ↓
                     Stopped → Removed
```

- **Created**: Container defined but not started
- **Starting**: Container process initializing
- **Waiting**: Blocked on dependent services reaching healthy state
- **Healthy**: All health checks passing
- **Unhealthy**: One or more health checks failing (backend reports degraded if DB/Redis down)
- **Stopped**: Container gracefully shut down
- **Removed**: Container and associated resources cleaned up
