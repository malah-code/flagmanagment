# Data Model: Terraform Provider

The Terraform provider relies on the FlagManagment Management REST API.

## Provider Configuration
- `api_url` (String, Optional): The base URL for the FlagManagment instance. Defaults to `https://api.flagmanagment.com`.
- `api_key` (String, Required): The Service Account API key for authentication.

## Entity: Project (`flagmanagment_project`)
- `id` (String, Computed): The UUID of the project.
- `name` (String, Required): Human-readable name.
- `description` (String, Optional): Project description.

## Entity: Environment (`flagmanagment_environment`)
- `id` (String, Computed): The UUID of the environment.
- `project_id` (String, Required): The ID of the parent project.
- `name` (String, Required): e.g., "Production", "Staging".
- `is_protected` (Boolean, Optional): Defaults to `false`. If `true`, modifications require Change Requests.

## Entity: Feature Flag (`flagmanagment_feature_flag`)
- `id` (String, Computed): The UUID of the flag definition.
- `project_id` (String, Required): The ID of the parent project.
- `key` (String, Required): The unique flag identifier (e.g., `new-checkout-flow`).
- `name` (String, Required): Human-readable name.
- `type` (String, Required): `boolean` or `multivariate`.
- `parent_flag_id` (String, Optional): For sequential dependencies.

## Entity: Flag State (`flagmanagment_flag_state`)
- `environment_id` (String, Required): The target environment.
- `flag_id` (String, Required): The flag being configured.
- `enabled` (Boolean, Required): Global on/off state in this environment.
- `default_variant` (String, Required): The fallback variant name (e.g., `true`, `false`, or custom).
- `variants` (Block List, Optional): Definitions of variants and their values.
  - `name` (String, Required)
  - `value` (JSON String, Required)
- `targeting_rules` (Block List, Optional):
  - `name` (String, Required)
  - `variant` (String, Optional): The variant to serve if matched.
  - `rollout` (Block, Optional): Percentage-based rollout configuration.
  - `conditions` (Block List, Required):
    - `attribute` (String, Required)
    - `operator` (String, Required): `EQUALS`, `CONTAINS`, `REGEX`, `IN`
    - `values` (List of Strings, Required)

## Entity: Service Account (`flagmanagment_service_account`)
- `id` (String, Computed): UUID.
- `name` (String, Required): Identifier.
- `project_id` (String, Optional): Scope.
- `role` (String, Required): `ADMIN`, `EDITOR`, `VIEWER`.
- `token` (String, Computed, Sensitive): The generated API key token.
