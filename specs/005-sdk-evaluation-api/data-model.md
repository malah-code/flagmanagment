# Data Model: SDK Evaluation API

## Entities

### `RulesetSnapshot` (Payload Model - Not DB)
Represents the complete state of an environment's flags and rules.
- `environmentId`: UUID
- `version`: string (hash representing the current state of rules)
- `flags`: List of `FlagRule` objects

### `FlagRule` (Payload Model)
- `key`: string
- `type`: enum (BOOLEAN, STRING, NUMBER, JSON)
- `defaultState`: boolean (enabled/disabled)
- `defaultVariation`: string/value
- `rules`: List of targeting rules (ordered by priority)

### `EvaluationContext` (Payload Model)
The context passed by thin clients to evaluate a flag.
- `environmentKey`: string (derived from auth token)
- `flagKey`: string
- `context`: JSON map (e.g., `{"userId": "123", "email": "test@example.com"}`)

## State Transitions & Workflows

### SDK Bootstrapping
1. SDK authenticates with Environment Token.
2. API validates token against Database/Redis cache.
3. API fetches the pre-compiled `RulesetSnapshot` from Redis.
4. Payload returned to SDK.

### Flag Update (Management API to SDK)
1. User updates a flag via Dashboard -> Management API.
2. Management API updates PostgreSQL DB.
3. Management API recalculates `RulesetSnapshot` and pushes to Redis.
4. Management API publishes an event to Redis Pub/Sub (`environment:<id>:updates`).
5. All connected gRPC streams listening on this topic receive the delta payload and push to their respective SDKs.

## Validation Rules
- Environment Tokens must be valid and active.
- PII inside `EvaluationContext` must be hashed using MurmurHash3 before evaluating bucketing logic.
