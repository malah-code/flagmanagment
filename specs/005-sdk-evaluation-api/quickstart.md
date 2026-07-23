# Quickstart Validation Guide: SDK Evaluation API

This guide walks you through verifying the SDK flag evaluation endpoints locally.

## Prerequisites
- Docker Compose running (for PostgreSQL and Redis).
- Go backend running.
- A valid Environment SDK Token (obtained from the dashboard).
- Wait for gRPC UI tools (like `grpcurl`) or `curl` to be available.

## Setup
1. Ensure the backend is running with Redis enabled.
   ```bash
   docker compose up -d redis db
   go run cmd/server/main.go
   ```
2. Using the Dashboard or API, create a Project, Environment, and a Feature Flag (e.g., `test-flag` = true). Note the Environment SDK token.

## Validation Scenarios

### Scenario 1: Server-side Evaluation (Thin Client REST API)
1. Execute a `curl` POST request to the `/api/v1/sdk/evaluate` endpoint.
2. Provide the Environment SDK token in the `Authorization` header.
   ```bash
   curl -X POST http://localhost:8080/api/v1/sdk/evaluate \
     -H "Authorization: Bearer <ENVIRONMENT_TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"flagKey": "test-flag", "context": {"userId": "123"}}'
   ```
3. **Expected**: The API returns a 200 OK with `{"value": "true", "reason": "DEFAULT"}` (or similar).

### Scenario 2: Thick Client Snapshot Fetch (gRPC)
1. Using `grpcurl`, request the initial ruleset snapshot.
   ```bash
   grpcurl -plaintext \
     -d '{"environment_token": "<ENVIRONMENT_TOKEN>"}' \
     localhost:9090 flagmanagement.sdk.v1.SDKService/FetchSnapshot
   ```
2. **Expected**: The gRPC response contains the `version` string and the list of `flags` with their default variations and targeting rules in JSON format.

### Scenario 3: Real-time Delta Updates (gRPC Stream)
1. Establish a streaming connection using `grpcurl`.
   ```bash
   grpcurl -plaintext \
     -d '{"environment_token": "<ENVIRONMENT_TOKEN>", "last_known_version": ""}' \
     localhost:9090 flagmanagement.sdk.v1.SDKService/StreamRulesets
   ```
2. Keep the connection open.
3. In a separate terminal or the Dashboard, modify `test-flag` (e.g., disable it).
4. **Expected**: The `grpcurl` terminal immediately receives a `RulesetDelta` payload indicating the flag was updated.
