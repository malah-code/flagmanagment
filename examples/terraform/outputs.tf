output "project_id" {
  description = "UUID of the created FlagManagment Project"
  value       = flagmanagment_project.ecommerce.id
}

output "staging_env_id" {
  description = "UUID of the Staging environment"
  value       = flagmanagment_environment.staging.id
}

output "production_env_id" {
  description = "UUID of the Production environment"
  value       = flagmanagment_environment.production.id
}

output "checkout_flag_id" {
  description = "UUID of the New Checkout Feature Flag"
  value       = flagmanagment_feature_flag.new_checkout.id
}

output "service_account_token" {
  description = "Client token generated for the CI/CD evaluation agent"
  value       = flagmanagment_service_account.evaluator_bot.api_key
  sensitive   = true
}
