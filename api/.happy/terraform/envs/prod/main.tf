module "stack" {
  source              = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/happy-ecs-stack?ref=dbd73dd6c0ac86fc6c9af940de9398bc201fb8fb"
  app_name            = "happy-api"
  happy_config_secret = var.happy_config_secret
  image_tag           = var.image_tag
  image_tags          = jsondecode(var.image_tags)
  priority            = var.priority
  stack_name          = var.stack_name
  deployment_stage    = "prod"
  require_okta        = true
  stack_prefix        = "/${var.stack_name}"
  wait_for_steady_state = var.wait_for_steady_state
}
