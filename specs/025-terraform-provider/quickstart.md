# Quickstart: Terraform Provider Validation

This guide explains how to validate the FlagManagment Terraform Provider once built.

## Prerequisites
- Terraform CLI >= 1.5.0 installed
- FlagManagment server running locally (e.g., via Docker Compose) on `http://localhost:8080`
- An existing API Key for the FlagManagment admin user.

## Setup

1. **Build the provider:**
   Navigate to the provider source code and compile it:
   ```bash
   cd providers/terraform
   make build
   ```

2. **Configure local Terraform overrides:**
   To test the provider locally without publishing to the Terraform Registry, create a `~/.terraformrc` file:
   ```hcl
   provider_installation {
     dev_overrides {
       "malah-code/flagmanagment" = "/path/to/flagmanagment/providers/terraform/bin"
     }
     direct {}
   }
   ```

3. **Create a test configuration:**
   Create a `main.tf` file:
   ```hcl
   terraform {
     required_providers {
       flagmanagment = {
         source  = "malah-code/flagmanagment"
       }
     }
   }
   
   provider "flagmanagment" {
     api_url = "http://localhost:8080"
     api_key = "test-admin-key"
   }
   
   resource "flagmanagment_project" "test" {
     name = "TF Test Project"
   }
   ```

## Validation Steps

1. **Initialize Terraform:**
   ```bash
   terraform init
   ```
   *Expected Outcome*: Success message (ignoring the dev_overrides warning).

2. **Plan the changes:**
   ```bash
   terraform plan
   ```
   *Expected Outcome*: Terraform shows that 1 resource (`flagmanagment_project.test`) will be created.

3. **Apply the changes:**
   ```bash
   terraform apply -auto-approve
   ```
   *Expected Outcome*: Apply completes successfully. Query the FlagManagment API (`GET /projects`) to confirm the project exists.

4. **Verify idempotency:**
   Run `terraform plan` again.
   *Expected Outcome*: "No changes. Your infrastructure matches the configuration."

5. **Destroy the infrastructure:**
   ```bash
   terraform destroy -auto-approve
   ```
   *Expected Outcome*: Project is successfully deleted from FlagManagment.
