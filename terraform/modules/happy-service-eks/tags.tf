locals {
  tags = merge(var.tags, {
    happy_env          = var.deployment_stage,
    happy_stack_name   = var.stack_name,
    happy_service_name = var.routing.service_name,
    happy_region       = data.aws_region.current.name,
    happy_image_tag    = var.image_tag
    happy_service_type = var.routing.service_type,
    happy_last_applied = timestamp(),
  })

  routing_tags = merge(var.tags, {
    happy_env          = var.deployment_stage,
    happy_stack_name   = var.stack_name,
    happy_region       = data.aws_region.current.name,
    happy_last_applied = timestamp(),
    happy_service_type = var.routing.service_type,
    happy_service_name = var.routing.service_name,
  })
}
