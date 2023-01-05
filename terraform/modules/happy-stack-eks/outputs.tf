output "service_endpoints" {
  value       = local.service_endpoints[0]
  description = "The URL endpoints for services"
}

output "task_arns" {
  value       = { for name, task in module.tasks : name => task.task_definition_arn }
  description = "ARNs for all the tasks"
}
