locals {
  tags = {
    happy-env           = var.deployment_stage,
    happy-stack-name    = var.stack_name,
    happy-service-name  = var.service_name,
    happy-region        = data.aws_region.current.name,
    happy-image         = var.image,
    happy-service-type  = var.service_type,
    happy-last-applied  = timestamp(),
  }
}
