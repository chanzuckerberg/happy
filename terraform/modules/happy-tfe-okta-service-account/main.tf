module "happy_apps" {
  for_each = var.envs
  source   = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth?ref=heathj/jwks"

  okta = {
    label         = "${var.service_name}-${var.app_name}-${each.value}-service-account"
    redirect_uris = concat(["https://oauth.${var.app_name}.${each.value}.si.czi.technology/oauth2/callback"], var.redirect_uris)
    login_uri     = var.login_uri == "" ? "https://oauth.${var.app_name}.${each.value}.si.czi.technology" : var.login_uri
    tenant        = "czi-prod"
  }

  grant_types                = ["client_credentials"]
  app_type                   = "service"
  token_endpoint_auth_method = "private_key_jwt"
  response_types             = ["token"]

  tags = {
    owner   = "infra-eng@chanzuckerberg.com"
    service = "${var.service_name}-oauth"
    project = var.app_name
    env     = each.value
  }
  aws_ssm_paths = var.aws_ssm_paths
  jwks = var.jwks
  # we set at least on role so that an authorization server is created
  rbac_role_mapping = merge({
    base: []
    }, var.rbac_role_mapping)
}

resource "okta_app_group_assignments" "happy_app" {
  for_each  = module.happy_apps
  app_id    = each.value.app.id
  group_ids = var.teams
}
