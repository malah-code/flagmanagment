variable "flagmanagment_api_url" {
  description = "Base URL of the FlagManagment API server"
  type        = string
  default     = "http://localhost:8080"
}

variable "flagmanagment_api_key" {
  description = "FlagManagment API Admin or Service Key"
  type        = string
  sensitive   = true
  default     = "test-admin-key"
}

variable "project_name" {
  description = "The name of the project to manage via Terraform"
  type        = string
  default     = "Flagship Storefront"
}
