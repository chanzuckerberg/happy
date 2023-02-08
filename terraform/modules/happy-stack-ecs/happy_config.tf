data "happy_resolved_app_configs" "configs" {
  app_name    = local.secret["tags"]["project"]
  environment = var.deployment_stage
  stack       = var.stack_name
}

locals {
  stack_configs = { for v in data.happy_resolved_app_configs.configs.app_configs : v["key"] => v["value"] }
}
