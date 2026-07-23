#!/bin/bash
set -e

if ! command -v docker &> /dev/null; then
    echo "ERROR: Docker is not installed. Install it from https://docs.docker.com/get-docker/"
    exit 1
fi

if ! docker compose version &> /dev/null; then
    echo "ERROR: Docker Compose v2 is required. Update Docker Desktop or install the compose plugin."
    exit 1
fi

if ! docker info > /dev/null 2>&1; then
    echo "ERROR: Docker daemon is not running. Start Docker Desktop or run 'sudo systemctl start docker'."
    exit 1
fi

DOCKER_VERSION=$(docker version --format '{{.Server.Version}}')
if [[ "$(printf '%s\n' "24.0" "$DOCKER_VERSION" | sort -V | head -n1)" != "24.0" ]]; then
    echo "WARNING: Docker 24.0+ recommended. You have $DOCKER_VERSION. Some features may not work."
fi

check_port() {
    local port=$1
    if command -v lsof &> /dev/null; then
        if lsof -i :$port > /dev/null 2>&1; then
            echo "ERROR: Port $port is already in use. Either stop the conflicting service or set FM_BACKEND_PORT/FM_FRONTEND_PORT/FM_DB_PORT/FM_REDIS_PORT in your .env file."
            exit 1
        fi
    elif command -v ss &> /dev/null; then
        if ss -tlnp | grep :$port > /dev/null 2>&1; then
            echo "ERROR: Port $port is already in use. Either stop the conflicting service or set FM_BACKEND_PORT/FM_FRONTEND_PORT/FM_DB_PORT/FM_REDIS_PORT in your .env file."
            exit 1
        fi
    fi
}

check_port 8080
check_port 3000
check_port 5432
check_port 6379

ARCH=$(uname -m)
if [[ "$ARCH" != "x86_64" && "$ARCH" != "aarch64" ]]; then
    echo "WARNING: Architecture $ARCH is not officially supported. Supported: x86_64, aarch64 (ARM64)."
fi

if [ ! -f .env ]; then
    if [ -f .env.example ]; then
        cp .env.example .env
        echo "INFO: Created .env from .env.example with default values."
    fi
fi

echo "Starting FlagManagment services..."
docker compose up --build -d

echo "Waiting for services to become healthy..."

TIMEOUT=120
ELAPSED=0
BACKEND_PORT=${FM_BACKEND_PORT:-8080}

while [ $ELAPSED -lt $TIMEOUT ]; do
    if curl -sf http://localhost:$BACKEND_PORT/healthz > /dev/null 2>&1; then
        echo ""
        curl -s http://localhost:$BACKEND_PORT/healthz
        echo ""
        echo "✅ FlagManagment is running!"
        echo "   Backend:    http://localhost:$BACKEND_PORT/healthz"
        echo "   Dashboard:  http://localhost:${FM_FRONTEND_PORT:-3000}"
        echo "   PostgreSQL: localhost:${FM_DB_PORT:-5432}"
        echo "   Redis:      localhost:${FM_REDIS_PORT:-6379}"
        echo "   Logs:       docker compose logs -f"
        echo "   Stop:       make down"
        exit 0
    fi
    sleep 5
    ELAPSED=$((ELAPSED + 5))
done

echo "WARNING: Backend health check did not pass within 120 seconds. Check logs with: docker compose logs backend"
exit 1
