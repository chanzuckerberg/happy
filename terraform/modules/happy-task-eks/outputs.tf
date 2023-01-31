output "task_definition_arn" {
  value       = kubernetes_cron_job.task_definition.metadata[0].name
  description = "Task definition name"
}
