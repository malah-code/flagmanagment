# Implementation Plan: Terraform Provider

**Branch**: `[025-terraform-provider]` | **Date**: 2026-08-08 | **Spec**: [specs/025-terraform-provider/spec.md](spec.md)

**Input**: Feature specification from `/specs/025-terraform-provider/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Build a HashiCorp Configuration Language (HCL) compatible Terraform provider for FlagManagment to allow declarative infrastructure-as-code management of projects, environments, flags, and service accounts.

## Technical Context

**Language/Version**: Go 1.21+

**Primary Dependencies**: `hashicorp/terraform-plugin-framework`, custom internal REST API client for FlagManagment Management API.

**Storage**: Terraform State files

**Testing**: `hashicorp/terraform-plugin-testing`, standard Go testing for the API client.

**Target Platform**: Multi-arch Terraform CLI plugin (Linux, macOS, Windows).

**Project Type**: Terraform Provider CLI Plugin.

**Performance Goals**: Fast parallel resource creation; sub-second plan outputs.

**Constraints**: Must handle protected environments correctly (graceful error/warnings on pending change requests or bypass capability).

**Scale/Scope**: ~10 resources and data sources mapping to the FlagManagment management API.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **API-First Contract Design**: N/A for exposing an API, but the provider acts as an API consumer.
- **Environment Isolation**: Enforced by separate resources `flagmanagment_environment` and `flagmanagment_flag_state`.
- **Governance by Default**: Handled via `is_protected` attributes on environments and Change Request handling mechanisms built into the provider configuration.
- **Test-First Quality Gates**: Acceptance tests using `terraform-plugin-testing` are mandatory.

## Project Structure

### Documentation (this feature)

```text
specs/025-terraform-provider/
├── plan.md              
├── research.md          
├── data-model.md        
├── quickstart.md        
├── contracts/           
│   └── terraform-hcl.md
└── tasks.md             
```

### Source Code (repository root)

```text
providers/terraform/
├── main.go
├── internal/
│   ├── provider/
│   │   ├── provider.go
│   │   ├── project_resource.go
│   │   ├── environment_resource.go
│   │   ├── feature_flag_resource.go
│   │   ├── flag_state_resource.go
│   │   └── service_account_resource.go
│   └── client/
│       ├── client.go
│       ├── project.go
│       └── ...
├── GNUmakefile
└── go.mod
```

**Structure Decision**: The Terraform provider logic will be isolated in the `providers/terraform` directory, containing its own `go.mod` because Terraform providers require heavy `hashicorp` dependencies that should not pollute the main backend or SDK modules. A dedicated API client will be written under `internal/client/` to interact with the Management API, separating it from the evaluation SDK (`sdk/go/`).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

*(No violations)*
