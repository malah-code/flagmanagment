# Contract: Terraform HCL Examples

This document outlines the HashiCorp Configuration Language (HCL) syntax and expected structures for the FlagManagment Terraform Provider.

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
  api_key = var.flagmanagment_api_key
  api_url = "https://api.flagmanagment.internal"
  bypass_change_requests = false
}

resource "flagmanagment_project" "main" {
  name        = "Frontend Application"
  description = "React frontend feature flags"
}

resource "flagmanagment_environment" "production" {
  project_id   = flagmanagment_project.main.id
  name         = "Production"
  is_protected = true
}

resource "flagmanagment_feature_flag" "new_checkout" {
  project_id = flagmanagment_project.main.id
  key        = "new-checkout-flow"
  name       = "New Checkout Flow"
  type       = "boolean"
}

resource "flagmanagment_flag_state" "new_checkout_prod" {
  environment_id  = flagmanagment_environment.production.id
  flag_id         = flagmanagment_feature_flag.new_checkout.id
  enabled         = true
  default_variant = "false"

  variants {
    name  = "true"
    value = "true"
  }
  variants {
    name  = "false"
    value = "false"
  }

  targeting_rules {
    name    = "Internal Users Only"
    variant = "true"
    
    conditions {
      attribute = "email"
      operator  = "CONTAINS"
      values    = ["@flagmanagment.com"]
    }
  }

  targeting_rules {
    name = "Gradual Rollout"
    rollout {
      percentages = {
        "true"  = 20
        "false" = 80
      }
    }
  }
}
```
