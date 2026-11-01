# FlagManagment Terraform Blueprint & Examples

This directory provides an Infrastructure as Code (IaC) blueprint for provisioning and managing FlagManagment projects, environments, flags, targeting rollouts, and service accounts declaratively.

## Architecture

This blueprint manages:
1. **Projects**: Multi-tenant workspace definition (`flagmanagment_project`).
2. **Environments**: Isolated stages (`flagmanagment_environment`) with protection rules.
3. **Feature Flags**: Boolean, String, JSON, and Numeric feature flags (`flagmanagment_feature_flag`).
4. **Flag States & Rollouts**: Granular per-environment activation and percentage rollouts (`flagmanagment_flag_state`).
5. **Service Accounts**: Non-human machine accounts with scoped RBAC for CI/CD or edge SDKs (`flagmanagment_service_account`).
6. **Data Sources**: Query and reference existing projects and flags.

---

## Quickstart

### 1. Build & Install the Provider Locally (Optional during development)

```bash
cd ../../providers/terraform
go build -o terraform-provider-flagmanagment main.go
```

Or configure local dev overrides in `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "malah-code/flagmanagment" = "/path/to/flagmanagment/providers/terraform"
  }
  direct {}
}
```

### 2. Configure Environment Variables

```bash
export TF_VAR_flagmanagment_api_url="http://localhost:8080"
export TF_VAR_flagmanagment_api_key="your-flagmanagment-admin-token"
```

### 3. Initialize & Apply

```bash
cd examples/terraform
terraform init
terraform plan
terraform apply
```

### 4. Cleanup

```bash
terraform destroy
```
