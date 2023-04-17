output "service_endpoints" {
  value       = module.stack.service_endpoints
  description = "The URL endpoint for the frontend website service"
  sensitive   = false
}

output "service_ecrs" {
  value       = module.stack.service_ecrs
  description = "The services ECR locations for their docker containers"
  sensitive   = false
}