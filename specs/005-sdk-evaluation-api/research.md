# Research: SDK Evaluation API

## Unknown 1: SDK Streaming Protocol
- **Context**: The feature spec requires a mechanism for real-time delta updates to connected SDKs.
- **Decision**: gRPC with Protobuf streams.
- **Rationale**: The project constitution explicitly mandates `gRPC/Protobuf (internal and SDK streaming)` under Technology Stack Constraints.
- **Alternatives considered**: Server-Sent Events (SSE) over HTTP/1.1 (rejected due to constitution constraint, though it was assumed in the spec).

## Unknown 2: Thin Client Evaluation Protocol
- **Context**: Thin clients (mobile, frontend) need to evaluate flags without downloading the full ruleset.
- **Decision**: REST/JSON over HTTPS via `/api/v1/sdk/evaluate`.
- **Rationale**: The project constitution mandates `REST/JSON over HTTPS (external)` for external APIs. Thin clients typically run in browsers or on mobile devices where REST/JSON is highly compatible and easy to integrate without pulling in heavy gRPC dependencies.
- **Alternatives considered**: gRPC-Web (rejected as it requires an Envoy proxy and adds frontend complexity not yet established in this project).

## Unknown 3: Caching Strategy
- **Context**: Evaluation and ruleset endpoints must be sub-millisecond. The constitution mandates Redis 7+.
- **Decision**: Store the pre-compiled full environment ruleset as a single JSON/Protobuf payload in Redis, keyed by environment ID. Use Redis Pub/Sub to notify the Go application instances to push updates to connected gRPC streams. Thin client `/evaluate` will fetch rules from an in-memory application cache (e.g., Go `sync.Map` or `ristretto`) backed by Redis.
- **Rationale**: A single Redis fetch per request for rulesets is fast, but an in-memory cache in the Go application (updated via Redis Pub/Sub) guarantees the <1ms local evaluation target for thin-client requests processed on the server.
- **Alternatives considered**: Pure Redis caching without application-level in-memory cache (rejected as network hop to Redis might exceed the strict <1ms evaluation budget under high load).
