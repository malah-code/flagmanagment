# Research: One-Click Flag Environment Promotions

## Decisions

### 1. Cross-Environment Copying
- **Decision**: Read source `EnvironmentFlagState` and write to target `EnvironmentFlagState`. If target is protected, construct a `ChangeRequest` populated with the target's existing state as `old` and source's state as `new`.
- **Rationale**: Reuses the robust Change Request infrastructure.

### 2. Frontend Promotion Flow
- **Decision**: Add a "Promote" action button on the Flag Details / Environments list, opening a modal to select the target environment.
- **Rationale**: Standard UX pattern for cross-environment actions.
