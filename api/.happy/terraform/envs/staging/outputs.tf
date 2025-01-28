output "service_ecrs" {
  sensitive = false
  value     = module.stack.service_ecrs
}
output "service_endpoints" {
  description = "The URL endpoints for services"
  sensitive   = false
  value       = module.stack.service_endpoints
}
output "task_arns" {
  description = "ARNs for all the tasks"
  sensitive   = false
  value       = module.stack.task_arns
}
