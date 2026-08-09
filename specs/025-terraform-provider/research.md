# Research: Terraform Provider

## Decision 1: Terraform Plugin Framework vs SDKv2
**Decision**: We will use HashiCorp's `terraform-plugin-framework`.
**Rationale**: `terraform-plugin-framework` is the modern, recommended framework by HashiCorp. It provides better support for complex data types, nested blocks, and context-aware validation, which are critical for configuring complex targeting rules and multivariate variants in FlagManagment. SDKv2 is legacy.
**Alternatives considered**: HashiCorp `terraform-plugin-sdk/v2` (rejected due to legacy status).

## Decision 2: API Client
**Decision**: Create a dedicated management API client in Go (e.g., inside `internal/client/`) for the provider, rather than using the existing `sdk/go/` SDK.
**Rationale**: The existing `sdk/go/` SDK is an *evaluation* SDK (Server-Side) designed to connect via gRPC to evaluate flags. The Terraform provider needs to interact with the *management* REST API to create projects, environments, and modify flag state. These are fundamentally different operations. 
**Alternatives considered**: Extending the `sdk/go/` SDK (rejected because it mixes evaluation runtime dependencies with infrastructure management dependencies).

## Decision 3: Handling Protected Environments and Change Requests
**Decision**: We will add a provider-level configuration `bypass_change_requests` (bool, default `false`). If false, the provider will fail or require a two-step apply process when modifying protected environments, emitting a warning that a Change Request was created. If true (and the API Key has sufficient permissions), it will bypass the Change Request workflow.
**Rationale**: Terraform expects declarative, immediate consistency. A pending Change Request means the state in FlagManagment won't match the Terraform state until manually approved. The provider will handle this by returning the "Pending" status in the state, and waiting for the user to approve it out-of-band, then re-running `terraform apply`.
**Alternatives considered**: Forcing Terraform to poll until the change request is approved (rejected because approvals could take hours/days).

## Decision 4: Authentication
**Decision**: The provider will authenticate using a Service Account API Key passed via the `FLAGMANAGMENT_API_KEY` environment variable or the `api_key` provider configuration.
**Rationale**: Standard practice for Terraform providers (e.g., AWS, Datadog) to authenticate via Service Accounts.
**Alternatives considered**: OAuth2 device flow (rejected as it's not suitable for CI/CD).
