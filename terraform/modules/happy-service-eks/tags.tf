locals {
  happy_last_applied = timestamp()
  tags = merge(var.tags, {
    happy_env          = var.deployment_stage,
    happy_stack_name   = var.stack_name,
    happy_service_name = var.routing.service_name,
    happy_region       = data.aws_region.current.name,
    happy_image_tag    = var.image_tag
    happy_service_type = var.routing.service_type,
    happy_last_applied = local.happy_last_applied,
  })

  # the tags have to be exactly the same across all ingresses in the same ingress group
  # otherwise, you gets errors like
  # {"level":"error","ts":"2023-09-26T04:01:46Z","msg":"Reconciler error","controller":"ingress","object":{"name":"stack-jheath2-fond-albacore"},"namespace":"","name":"stack-jheath2-fond-albacore","reconcileID":"1c80d91d-4643-4b6c-9784-48f2d8e40976","error":"conflicting tag happy_service_name: jheath2-internal-api | jheath2-frontend"}
  # in context based routing, different service will share the same ingress group
  # so don't include their service specific tag information or the ALB won't be created
  routing_tags = merge(var.tags, {
    happy_env          = var.deployment_stage,
    happy_stack_name   = var.stack_name,
    happy_region       = data.aws_region.current.name,
    happy_last_applied = local.happy_last_applied,
  })
}
