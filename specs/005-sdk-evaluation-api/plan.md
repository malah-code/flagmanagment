# Implementation Plan: SDK Evaluation API

**Branch**: `005-sdk-evaluation-api` | **Date**: 2026-07-21 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/005-sdk-evaluation-api/spec.md`

## Summary

Implement the high-performance SDK Client API flag evaluation endpoints leveraging Redis and in-memory caching to achieve sub-millisecond local evaluations. Thick clients (server-side SDKs) will fetch a complete ruleset via gRPC and subscribe to real-time delta updates. Thin clients will evaluate flags server-side via a REST endpoint.

## Technical Context

**Language/Version**: Go 1.23+

**Primary Dependencies**: `grpc-go`, `protobuf`, `go-redis/v9`, `ristretto` (or similar Go in-memory cache)

**Storage**: Redis 7+ (Primary cache & Pub/Sub), PostgreSQL (Source of truth)

**Testing**: `go test`, `testcontainers-go` (for Redis integration testing)

**Target Platform**: Linux Server / Docker Container

**Project Type**: Web Service (Backend REST & gRPC API)

**Performance Goals**: < 10ms for ruleset snapshot, < 1ms for thin-client evaluate (when cached locally on the application server), sub 500ms push delta updates.

**Constraints**: Must use MurmurHash3 for targeting buckets (already implemented in core engine, if any), must NOT hit PostgreSQL on the hot path for evaluations.

**Scale/Scope**: Support 10k concurrent gRPC streaming connections per instance.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. API-First Contract Design**: gRPC Protobuf and OpenAPI specifications define the contracts.
- **IV. Local Evaluation Performance**: gRPC stream ensures immediate updates; Redis provides high-throughput ruleset serving.
- **VII. PII Protection**: Thin-client `evaluate` endpoint must hash PII before logging or metric extraction.
- **Technology Stack Constraints**: Go, Redis, REST/JSON, gRPC/Protobuf are all adhered to exactly as mandated.

## Project Structure

### Documentation (this feature)

```text
specs/005-sdk-evaluation-api/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (sdk.proto, evaluate-api.yaml)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
backend/
├── api/
│   ├── proto/
│   │   └── sdk/v1/sdk.proto      # The gRPC contract
│   └── openapi/
│       └── evaluate-api.yaml     # The REST contract
├── cmd/
│   └── server/main.go            # Initialize Redis & gRPC server
├── internal/
│   ├── cache/                    # Redis client and in-memory cache wrappers
│   ├── sdk/
│   │   ├── service.go            # gRPC service implementation
│   │   └── stream.go             # gRPC streaming and Pub/Sub manager
│   └── handlers/
│       └── sdk_evaluate.go       # Thin client REST handler
└── pkg/
    └── gen/
        └── sdk/v1/               # Protoc generated Go code
```

**Structure Decision**: Integrated directly into the existing Go backend web service. gRPC service will run alongside the existing HTTP server (either on a different port, e.g. 9090, or multiplexed if using a library like `cmux`).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

*(No violations found. The architecture perfectly aligns with the constitution constraints.)*
