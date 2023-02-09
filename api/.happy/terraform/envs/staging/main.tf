module "stack" {
  source                = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-ecs?ref=fd91423"
  app_name              = "hapi"
  happy_config_secret   = var.happy_config_secret
  image_tag             = var.image_tag
  image_tags            = jsondecode(var.image_tags)
  priority              = var.priority
  stack_name            = var.stack_name
  deployment_stage      = "staging"
  require_okta          = false
  stack_prefix          = "/${var.stack_name}"
  wait_for_steady_state = var.wait_for_steady_state
  chamber_service       = "happy-staging-hapi"
  service_port          = 3001
}
