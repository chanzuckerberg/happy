output "service_endpoints" {
  value       = local.service_endpoints
  description = "The URL endpoints for services"
  sensitive   = false
}

output "task_arns" {
  value       = { for name, task in module.tasks : name => task.task_definition_arn }
  description = "ARNs for all the tasks"
}

output "dashboard" {
  value = {
    id  = var.create_dashboard ? module.datadog_dashboard[0].id : ""
    url = var.create_dashboard ? module.datadog_dashboard[0].url : ""
  }
}

output "service_ecrs" {
  value = local.service_ecrs
}
