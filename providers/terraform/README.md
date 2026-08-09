# Terraform Provider for FlagManagment

The FlagManagment Terraform provider allows you to declaratively manage feature flag projects, environments, flags, targeting rules, service accounts, and telemetry triggers using HashiCorp Configuration Language (HCL).

## Usage

```hcl
terraform {
  required_providers {
    flagmanagment = {
      source  = "malah-code/flagmanagment"
      version = "~> 1.0.0"
    }
  }
}

provider "flagmanagment" {
  api_url = "http://localhost:8080"
  api_key = var.flagmanagment_api_key
}

resource "flagmanagment_project" "app" {
  name        = "E-Commerce App"
  description = "Production & Staging feature flags"
}

resource "flagmanagment_environment" "prod" {
  project_id   = flagmanagment_project.app.id
  name         = "Production"
  is_protected = true
}

resource "flagmanagment_feature_flag" "new_payment" {
  project_id = flagmanagment_project.app.id
  key        = "new-payment-gateway"
  name       = "New Payment Gateway"
  type       = "boolean"
}
```

## Developing the Provider

Build the provider binary:

```bash
make build
```

Run tests:

```bash
make test
```
