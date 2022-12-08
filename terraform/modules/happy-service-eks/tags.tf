locals {
  tags = {
    happy_env          = var.deployment_stage,
    happy_stack_name   = var.stack_name,
    happy_service_name = var.service_name,
    happy_region       = data.aws_region.current.name,
    happy_image        = var.image,
    happy_service_type = var.service_type,
    happy_last_applied = timestamp(),
  }
}
