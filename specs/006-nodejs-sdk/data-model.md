# Data Model: Node.js SDK

## Entities

### `FlagRule`
Represents the state and logic for a single feature flag as downloaded from the server.
- `key` (string): The unique flag identifier.
- `type` (string): 'BOOLEAN', 'MULTIVARIATE', etc.
- `enabled` (boolean): Master kill switch for the flag.
- `defaultVariation` (string): The fallback value if targeting fails.
- `targetingRules` (JSON): The advanced targeting JSON payload (rollouts, segments).

### `EvaluationContext`
The context passed by the user during evaluation.
- `identity` (string): The unique user ID for hashing/rollouts.
- `attributes` (Record<string, unknown>): Additional custom properties for segment targeting.

### `EvaluationResult`
The result returned by the local evaluator.
- `value` (boolean | string | number): The resolved flag value.
- `reason` (string): Why the value was chosen ('TARGETING_MATCH', 'DEFAULT', 'DISABLED').

## Internal State

### `RuleStore`
- `flags` (Map<string, FlagRule>): The in-memory cache of the ruleset snapshot.
- `version` (string): The current ruleset version/hash for synchronization.
