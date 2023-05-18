output "task_definition_arn" {
  value       = kubernetes_cron_job_v1.task_definition.metadata[0].name
  description = "Task definition name"
}
