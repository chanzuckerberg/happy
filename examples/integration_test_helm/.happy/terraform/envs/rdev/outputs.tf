output "service_urls" {
  value       = module.stack.service_endpoints
  description = "The URL endpoint for the frontend website service"
  sensitive   = false
}

output "service_ecrs" {
  value       = module.stack.service_ecrs
  description = "The services ECR locations for their docker containers"
  sensitive   = false
}

output "k8s_namespace" {
  value = data.kubernetes_namespace.happy-namespace.metadata.0.name
}
