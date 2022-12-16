data "aws_region" "current" {}

locals {
  tags = {
    happy_env          = var.deployment_stage,
    happy_stack_name   = var.stack_name,
    happy_service_name = "", // TODO: fill this is later
    happy_service_type = "", // TODO: fill this is later
    happy_region       = data.aws_region.current.name,
    happy_image        = local.image
    happy_last_applied = timestamp(),
  }
}
