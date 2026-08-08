# Data Model: Sequential Dependencies

## Schema Modifications

### 1. `feature_flags` Table

The existing `feature_flags` table requires a new self-referencing foreign key column to track dependency relationships.

- **`parent_flag_id`** (UUID, Nullable)
  - **Foreign Key**: References `feature_flags(id)`.
  - **Constraint**: `ON DELETE RESTRICT` (A parent flag cannot be deleted if it has dependent flags).
  - **Index**: Create an index on `parent_flag_id` to quickly fetch all dependents of a specific flag for UI visualization.

## Entity Relationships

- **Feature Flag to Feature Flag (1:N)**: A parent flag can have multiple dependent flags. A dependent flag can have only ONE parent flag.
- **Environment Flag State**: No schema changes required here. The dependency relationship is defined at the project-wide *Feature Flag* level, meaning the dependency graph structure applies across all environments. The *evaluation* of that dependency uses the specific environment's state.

## Validation Rules (API Layer)

1. **Cycle Prevention**: Before saving a `parent_flag_id`, the API must traverse the parent chain. If the target parent (or any of its ancestors) is the flag being updated, the API must return a `400 Bad Request` indicating a circular dependency.
2. **Depth Limit**: To guarantee sub-millisecond evaluation latency, dependency chains should be hard-capped at a maximum depth of 3 levels (e.g., A -> B -> C -> D).
3. **Cross-Project Isolation**: The `parent_flag_id` must reference a flag within the same `project_id`.

## State Transitions

- **Delete/Archive Parent**: An attempt to archive or delete a flag that is referenced as a `parent_flag_id` by any active flag must fail with a validation error instructing the user to remove the dependency first.
