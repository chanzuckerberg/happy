locals {
  tags = merge(var.tags, {
    happy_env          = var.deployment_stage,
    happy_stack_name   = var.stack_name,
    happy_service_name = var.task_name,
    happy_region       = data.aws_region.current.name,
    happy_image_tag    = var.image
    happy_service_type = var.is_cron_job ? "cronjob" : "task",
    happy_last_applied = timestamp(),
  })
}
