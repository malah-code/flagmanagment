terraform {
  required_version = ">= 1.5.0"
  required_providers {
    flagmanagment = {
      source  = "malah-code/flagmanagment"
      version = "~> 1.0.0"
    }
  }
}

provider "flagmanagment" {
  api_url = var.flagmanagment_api_url
  api_key = var.flagmanagment_api_key
}

# 1. Project Creation
resource "flagmanagment_project" "ecommerce" {
  name        = var.project_name
  description = "E-Commerce Flagship Application feature flag and targeting configuration"
}

# 2. Environments (Staging & Production)
resource "flagmanagment_environment" "staging" {
  project_id   = flagmanagment_project.ecommerce.id
  name         = "Staging"
  description  = "Pre-production testing environment"
  is_protected = false
}

resource "flagmanagment_environment" "production" {
  project_id   = flagmanagment_project.ecommerce.id
  name         = "Production"
  description  = "High-availability customer-facing production environment"
  is_protected = true
}

# 3. Feature Flags
resource "flagmanagment_feature_flag" "new_checkout" {
  project_id  = flagmanagment_project.ecommerce.id
  key         = "new-checkout-flow"
  name        = "New Checkout Experience"
  description = "Enables one-click Stripe & Apple Pay accelerated checkout"
  type        = "boolean"
}

resource "flagmanagment_feature_flag" "search_algorithm" {
  project_id  = flagmanagment_project.ecommerce.id
  key         = "search-ranking-algorithm"
  name        = "Search Ranking Algorithm Variant"
  description = "Multi-variate algorithmic recommendation ranking"
  type        = "string"
}

# 4. Environment-Specific Flag State Configuration
resource "flagmanagment_flag_state" "checkout_staging" {
  flag_id        = flagmanagment_feature_flag.new_checkout.id
  environment_id = flagmanagment_environment.staging.id
  is_enabled     = true
  default_value  = "true"
}

resource "flagmanagment_flag_state" "checkout_production" {
  flag_id        = flagmanagment_feature_flag.new_checkout.id
  environment_id = flagmanagment_environment.production.id
  is_enabled     = true
  default_value  = "false"
  rollout_percentage = 20
}

# 5. Service Account for CI/CD or Microservices
resource "flagmanagment_service_account" "evaluator_bot" {
  project_id  = flagmanagment_project.ecommerce.id
  name        = "ci-evaluation-agent"
  role        = "viewer"
  description = "Machine token for edge SDK flag evaluation and telemetry syncing"
}

# 6. Data Source Verification / Queries
data "flagmanagment_project" "lookup" {
  id = flagmanagment_project.ecommerce.id
}
