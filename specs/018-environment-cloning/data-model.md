# Data Model: Environment Cloning

The environment cloning feature relies on the existing entities but requires deep-copying associated records.

## Entities Involved in Cloning

### `environments`
- **Fields Copied**: `ProjectID`.
- **Fields Generated**: `ID` (new UUID), `Name` (provided by user request), `APIKeyHash` (newly generated token), `IsProtected` (defaults to `false` for ephemeral environments).

### `environment_flag_states`
- **Fields Copied**: `FeatureFlagID`, `BooleanState`, `TargetingRules` (JSONB), `RemoteConfig` (JSONB).
- **Fields Updated**: `EnvironmentID` (set to the newly generated environment UUID).
- **Process**: For every record in `environment_flag_states` where `EnvironmentID == SourceEnvironmentID`, insert a duplicate record referencing the new EnvironmentID.

### `audit_logs` (Side Effect)
- A single new record is inserted to track the cloning event.
- **ActionType**: `ENVIRONMENT_CLONED`
- **TargetEnvironmentID**: The ID of the newly created environment.
- **PreviousState**: `{ "source_environment_id": "<ID>" }`
- **NewState**: The created environment metadata.
