# FlagManagment

FlagManagment — Cloud-native feature flag management and remote configuration platform.

**Official Repository:** [github.com/malah-code/flagmanagment](https://github.com/malah-code/flagmanagment)

## Prerequisites

- Docker 24+ and Docker Compose
- Git
- Make (optional, but recommended)
- Go 1.22+ (for local development)
- Node.js 20+ (for local development)

## Quick Start

1. Start all services locally:
   ```bash
   make up
   ```

2. Verify the backend health:
   ```bash
   curl -s http://localhost:8080/healthz | python3 -m json.tool
   ```

3. Open the Dashboard:
   Navigate to [http://localhost:3000](http://localhost:3000) to see the system health overview.

## Development

- **Run all services with hot-reload**:
  ```bash
  make up
  ```
  The backend uses `air` for hot-reloading and the frontend uses `vite`. Changes to source code will automatically apply.

- **Run Tests**:
  ```bash
  make test
  ```

- **Run Linters**:
  ```bash
  make lint
  ```

## Architecture

The project consists of 4 core services:
1. **Backend**: Go API service (Port 8080)
2. **Frontend**: React/TypeScript Dashboard (Port 3000)
3. **PostgreSQL**: Relational database for persistent storage (Port 5432)
4. **Redis**: In-memory cache and pub/sub (Port 6379)

For more architectural details, see [Core Architecture Specifications](specs/001-core-architecture/).

## License

BSL / Fair-Source (TBD)
