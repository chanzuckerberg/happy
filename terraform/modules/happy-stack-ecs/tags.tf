data "aws_region" "current" {}

locals {
  tags = {
    happy_env           = var.deployment_stage,
    happy_stack_name    = var.stack_name,
    happy_service_name  = var.service_name, // missing
    happy_service_type  = var.service_type, // missing
    happy_region        = data.aws_region.current.name,
    happy_image         = local.image
    happy_last_applied  = timestamp(),
  }
}
