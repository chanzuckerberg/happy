data "happy_resolved_app_configs" "configs" {
  app_name    = len(var.app_name) == 0 ? local.secret["tags"]["project"] : var.app_name
  environment = var.deployment_stage
  stack       = var.stack_name
}

locals {
  stack_configs = { for v in data.happy_resolved_app_configs.configs.app_configs : v["key"] => v["value"] }
}
