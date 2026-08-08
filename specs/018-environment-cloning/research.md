# Research: API-Driven Environment Cloning

## Needs Clarification Resolutions

**How is a cloned environment named?**
- **Decision**: The clone API will accept a `name` payload (e.g., `PR-123-Test`) and an optional `description`.
- **Rationale**: CI/CD pipelines need to name the environments dynamically based on the build context.

**How do we guarantee transactional consistency during cloning?**
- **Decision**: The clone operation must be wrapped in a single database transaction using `pgx.BeginTx`. If any part of the clone fails (e.g., copying the environment flags), the entire transaction rolls back.
- **Rationale**: Prevents orphaned environments or incomplete flag states which could corrupt the data model or cause testing anomalies.

**What happens to audit logs of the source environment?**
- **Decision**: Audit logs are NOT cloned. The act of cloning is logged as a single "Environment Cloned" audit event in the project, but historical flag changes from the source environment remain tied exclusively to the source environment.
- **Rationale**: Audit logs represent immutable history of actions. Copying them would falsify the history of the new environment.

## Technology Choices & Patterns

**SDK Token Generation**
- **Decision**: The clone operation will generate a brand new cryptographically secure token utilizing the existing `GenerateAPIKey` service.
- **Rationale**: Strict isolation (Constitution Principle II) requires that environments never share tokens, preventing accidental evaluation bleed.
