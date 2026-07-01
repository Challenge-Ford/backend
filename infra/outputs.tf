output "ecr_repository_urls" {
  description = "ECR repository URLs for the backend components"
  value       = module.ecr.repository_urls
}

output "api_service_name" {
  description = "Backend API ECS service name"
  value       = module.api.service_name
}

output "worker_service_name" {
  description = "Backend worker ECS service name"
  value       = module.worker.service_name
}

output "api_domain" {
  description = "Public hostname the API is served on"
  value       = var.api_domain
}
