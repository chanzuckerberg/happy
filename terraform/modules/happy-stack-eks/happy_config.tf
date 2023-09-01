data "happy_resolved_app_configs" "configs" {
  count       = var.skip_config_injection ? 0 : 1
  app_name    = length(var.app_name) == 0 ? local.secret["tags"]["project"] : var.app_name
  environment = var.deployment_stage
  stack       = var.stack_name
}

locals {
  stack_configs = var.skip_config_injection ? {} : { for v in data.happy_resolved_app_configs.configs[0].app_configs : v["key"] => v["value"] }
}

provider "happy" {
  api_base_url        = local.secret["hapi_config"]["base_url"]
  api_oidc_issuer     = local.secret["hapi_config"]["oidc_issuer"]
  api_oidc_authz_id   = local.secret["hapi_config"]["oidc_authz_id"]
  api_kms_key_id      = local.secret["hapi_config"]["kms_key_id"]
  api_assume_role_arn = local.secret["hapi_config"]["assume_role_arn"]
  api_oidc_scope      = local.secret["hapi_config"]["scope"]
}
