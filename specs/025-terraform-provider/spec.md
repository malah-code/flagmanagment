# Feature Specification: Terraform Provider

**Feature Branch**: `[025-terraform-provider]`

**Created**: 2026-08-08

**Status**: Draft

**Input**: User description: "Terraform provider" (as specified in Phase 3 of PRD)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Declarative Infrastructure as Code (Priority: P1)

As a DevOps Engineer, I want to define and manage my FlagManagment projects, environments, and feature flags using HashiCorp Configuration Language (HCL) so that I can version control my feature flag infrastructure alongside my application deployments.

**Why this priority**: Enterprise teams manage infrastructure deterministically. Without a Terraform provider, feature flags require manual UI clicks or custom API scripts, violating GitOps workflows.

**Independent Test**: Can be fully tested by running `terraform apply` with a basic `.tf` file and verifying that the project, environments, and flags are correctly created in the FlagManagment database.

**Acceptance Scenarios**:

1. **Given** an empty FlagManagment instance, **When** I apply a Terraform configuration defining a project and environments, **Then** the resources are correctly created and state is tracked by Terraform.
2. **Given** existing resources managed by Terraform, **When** I modify a flag's default variant in the `.tf` file and run `terraform apply`, **Then** the change is executed as a Change Request (if protected) or immediately applied (if unprotected).

---

### User Story 2 - Managing Governance and RBAC via Terraform (Priority: P2)

As a Platform Engineer, I want to provision API keys, service accounts, and Role-Based Access Control (RBAC) permissions using Terraform so that I can securely automate access provisioning.

**Why this priority**: Security teams require auditable, code-reviewed infrastructure for API credentials and permissions, especially for CI/CD integrations.

**Independent Test**: Can be tested by provisioning a service account and API key via Terraform, then using that API key to query the FlagManagment REST API.

**Acceptance Scenarios**:

1. **Given** a new environment, **When** I use Terraform to create a Service Account with "Read-Only" role, **Then** the corresponding token is generated and outputted.

---

### User Story 3 - Telemetry and Rollback Triggers Integration (Priority: P3)

As an SRE, I want to define telemetry triggers and automated rollback actions on specific flags via Terraform.

**Why this priority**: Links FlagManagment's kill-switches with existing infrastructure-as-code deployments for Datadog or Prometheus monitors.

**Independent Test**: Test by provisioning a flag with attached telemetry triggers using Terraform, and verifying in the FlagManagment UI that the triggers exist.

**Acceptance Scenarios**:

1. **Given** a feature flag configuration in Terraform, **When** I add a `telemetry_trigger` block and apply, **Then** the trigger configuration is persisted to the environment flag state.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The Terraform provider MUST support CRUD operations for `flagmanagment_project` resources.
- **FR-002**: The Terraform provider MUST support CRUD operations for `flagmanagment_environment` resources, including managing the `is_protected` attribute.
- **FR-003**: The Terraform provider MUST support CRUD operations for `flagmanagment_feature_flag` resources (defining key, type, and parent dependencies).
- **FR-004**: The Terraform provider MUST support CRUD operations for `flagmanagment_flag_state` resources, allowing configuration of the `enabled` state, `default_variant`, and multivariate `variants`.
- **FR-005**: The Terraform provider MUST support defining complex `targeting_rules` (including rollout percentages and contextual operators) within the flag state resource.
- **FR-006**: The Terraform provider MUST support CRUD operations for `flagmanagment_service_account` and API key generation.
- **FR-007**: When mutating resources in a "Protected" environment, the provider MUST automatically create a pending Change Request and wait for approval (or fail gracefully depending on provider configuration), rather than forcing the state change immediately.

### Key Entities

- **Project Resource**: Represents the top-level container for environments and flags.
- **Environment Resource**: Represents the deployment tier (e.g., Staging, Prod).
- **Feature Flag Resource**: Represents the generic flag definition (Key, Type).
- **Flag State Resource**: Represents the environment-specific configuration of a flag (Targeting Rules, Rollouts).
- **Service Account Resource**: Represents machine identities for API access.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: DevOps engineers can provision a completely empty FlagManagment cluster to a production-ready state (with projects, environments, and base flags) entirely via `terraform apply` in under 3 minutes.
- **SC-002**: Terraform plan outputs match the expected state of the FlagManagment API with zero drift upon initial import or clean apply.
- **SC-003**: The provider binary compiles and operates successfully on Linux (amd64, arm64), macOS (amd64, arm64), and Windows.
- **SC-004**: End-to-end Terraform acceptance tests (using the HashiCorp `terraform-plugin-testing` framework) pass with 100% success against a live FlagManagment container.

## Assumptions

- The FlagManagment REST/gRPC API is stable and feature-complete enough to support full CRUD for all Terraform-managed resources.
- HashiCorp's `terraform-plugin-framework` (Go) will be used to build the provider to align with modern HashiCorp standards and the existing Go backend stack.
- Terraform runs are generally performed by CI/CD systems with sufficient access privileges to mutate FlagManagment infrastructure.
- For protected environments, Terraform applies might require a two-step process (apply creates Change Request, manual approval in UI, second apply updates state) OR the provider is run with a bypass token. This behavior will be configured via provider settings.
