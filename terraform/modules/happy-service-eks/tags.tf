locals {
  tags = {
    happy-env           = var.deployment_stage,
    happy-stack-name    = var.custom_stack_name,
    happy-service-name  = var.app_name,
    happy-region        = data.aws_region.current.name,
    happy-image         = var.image,
    happy-service-type  = var.service_type,
    happy-last-applied  = timestamp(),
  }
}
