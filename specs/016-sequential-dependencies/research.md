# Research: Sequential Dependencies

## Decision: Circular Dependency Detection

- **Decision**: Use a depth-first search (DFS) cycle detection algorithm at the API layer when a feature flag's `parent_flag_id` is updated.
- **Rationale**: Feature flags in this system form a directed graph. Preventing cycles (e.g., A -> B -> A) at creation/update time ensures the evaluation engine never enters an infinite loop, satisfying Constitution Principle IV (sub-millisecond evaluation latency). 
- **Alternatives considered**: 
  - *Database-level recursive CTEs*: Rejected due to potential performance hits and complexity across different RDBMS systems. In-memory graph traversal on the API server is faster for the relatively small number of flags in a single environment.
  - *SDK-level detection*: Rejected because it pushes complexity to all client SDKs and risks breaking applications at runtime rather than failing safely at configuration time.

## Decision: SDK Evaluation Order

- **Decision**: The Server-Side SDK local evaluation engine will lazily evaluate the parent flag only when the dependent flag is requested. If the parent evaluates to a non-matching state (e.g., OFF), the engine short-circuits and returns the dependent flag's fallback state.
- **Rationale**: Avoids unnecessary evaluations and maintains the sub-millisecond constraint. The SDK already downloads a full snapshot of all flags; it simply follows the `parent_flag_id` pointer in memory.
- **Alternatives considered**:
  - *Pre-computing states on the backend*: Rejected because targeting rules often depend on dynamic context (e.g., User ID) provided only at runtime by the application. Evaluation must remain local to the SDK.
