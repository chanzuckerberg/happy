output service_endpoints {
  value       = local.service_endpoints
  description = "The URL endpoints for services"
}

output delete_db_task {
  value       = try(module.tasks["deletion"].task_definition_arn, "")
  description = "ARN of the Deletion Task Definition"
}

output migrate_db_task {
  value       = try(module.tasks["migration"].task_definition_arn, "")
  description = "ARN of the Migration Task Definition"
}
