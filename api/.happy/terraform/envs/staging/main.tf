module "stack" {
  source                = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/happy-ecs-stack?ref=happy-ecs-stack-v0.155.0"
  app_name              = "happy-api"
  happy_config_secret   = var.happy_config_secret
  image_tag             = var.image_tag
  image_tags            = jsondecode(var.image_tags)
  priority              = var.priority
  stack_name            = var.stack_name
  deployment_stage      = "staging"
  require_okta          = true
  stack_prefix          = "/${var.stack_name}"
  wait_for_steady_state = var.wait_for_steady_state
  chamber_service       = "happy-staging-api"
}
