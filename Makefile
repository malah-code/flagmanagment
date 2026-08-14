.PHONY: up down test lint build logs clean build-setup health fmt

# Start all services (with prerequisite validation)
up:
	@./scripts/bootstrap.sh

# Stop all services
down:
	docker compose down

# Stop and remove all data volumes (clean slate)
clean:
	docker compose down -v

# Run all tests
test: test-backend

test-backend:
	cd backend && go test -race -cover -coverprofile=coverage.out ./...

test-frontend:
	cd frontend && npm test -- --coverage

# Run all linters
lint: lint-backend lint-frontend

lint-backend:
	cd backend && golangci-lint run ./...

lint-frontend:
	cd frontend && npx -y eslint src/ --ext .ts,.tsx
	cd frontend && npx -y prettier --check "src/**/*.{ts,tsx,css}"

build-setup:
	@docker buildx inspect flagmanagment-builder > /dev/null 2>&1 || \
	  docker buildx create --name flagmanagment-builder --use --driver docker-container
	@docker buildx use flagmanagment-builder

# Multi-architecture Docker build (does not push)
build: build-setup
	docker buildx build --platform linux/amd64,linux/arm64 -t flagmanagment/backend:latest -f backend/Dockerfile backend/
	docker buildx build --platform linux/amd64,linux/arm64 -t flagmanagment/frontend:latest -f frontend/Dockerfile frontend/

# View logs
logs:
	docker compose logs -f

# Check health
health:
	@curl -s http://localhost:8080/healthz | python3 -m json.tool 2>/dev/null || echo "Backend not running"

# Auto-fix formatting across all code
fmt:
	cd frontend && npx -y prettier --write "src/**/*.{ts,tsx,css}"
	cd backend && gofmt -w .

# Run database migrations up
migrate-up:
	migrate -path backend/migrations -database "postgres://flagmgmt:flagmgmt_dev@localhost:5432/flagmanagment?sslmode=disable" up

# Run database migrations down
migrate-down:
	migrate -path backend/migrations -database "postgres://flagmgmt:flagmgmt_dev@localhost:5432/flagmanagment?sslmode=disable" down -all

# Seed database with realistic sample projects, environments, and feature flags
seed:
	@node scripts/seed.js

# Wipe all existing projects/flags and seed fresh enterprise sample data
seed-reset: reseed
reseed:
	@node scripts/reset-seed.js
