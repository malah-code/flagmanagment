# Quickstart: Running the Comprehensive Test Suite

**Branch**: `038-comprehensive-test-plan` | **Date**: 2026-10-31

This guide documents how to execute, validate, and benchmark the FlagManagment test suites across all layers.

---

## 1. Prerequisites

Ensure local development prerequisites are running:
```bash
# 1. Start core datastores (PostgreSQL 16 & Redis 7)
docker compose up -d postgres redis

# 2. Verify Go & Node.js environments
go version        # Expected: >= go1.25.0
node -v           # Expected: >= v20.0.0
```

---

## 2. Layer-by-Layer Test Execution

### 2.1 Backend Unit & Integration Tests (Go)
Validates API endpoints, RBAC, change requests, scheduler, and database logic:
```bash
cd backend
go test -v -race -cover ./...
```
*To generate coverage HTML profile:*
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 2.2 Frontend Component & Unit Tests (Vitest)
Validates UI components, Remote Config builder, Login error flows, and Context switching:
```bash
cd frontend
npm test                 # Runs vitest run
npm run test:coverage    # Runs vitest run --coverage
```

### 2.3 End-to-End User Journey Tests (Playwright)
Executes browser tests simulating user and admin workflows against live services:
```bash
# Start backend and frontend in Docker
docker compose up -d --build

# Run Playwright E2E tests
cd e2e
npx playwright test
```

### 2.4 Server-Side SDK In-Memory Performance Benchmark
Validates **Constitution Principle IV** (sub-millisecond evaluation latency):
```bash
# Run Go evaluation benchmark
cd backend
go test -bench=BenchmarkEvaluateFlag -benchmem ./internal/sdk/...
```
*Assertion: Latency MUST report `< 1,000,000 ns/op (1.0 ms)`.*

---

## 3. Unified Suite Execution

To run all tiers sequentially with consolidated pass/fail reporting and threshold verification:
```bash
# From project root:
./scripts/run-all-tests.sh --all --json
```

---

## 4. Syncing the Test Execution Tracker

After executing tests, synchronize the canonical tracker document:
```bash
python3 scripts/sync-test-tracker.py --results test-results.json
```
This updates [`docs/TEST_EXECUTION_TRACKER.md`](../../docs/TEST_EXECUTION_TRACKER.md) with real-time pass/fail status across the 158 test cases.

