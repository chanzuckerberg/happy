module "happy_app" {
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth-head?ref=v0.249.0"

  okta = {
    label         = "${var.service_name}-${var.app_name}-${var.env}"
    redirect_uris = concat(["https://*.${var.app_name}.${var.env}.si.czi.technology/oauth2/idpresponse"], var.redirect_uris)
    login_uri     = var.login_uri == "" ? "https://oauth.${var.app_name}.${var.env}.si.czi.technology" : var.login_uri
    tenant        = "czi-prod"
  }

  grant_types                = var.grant_types
  app_type                   = var.app_type
  token_endpoint_auth_method = var.token_endpoint_auth_method
  omit_secret                = var.omit_secret

  tags = {
    owner   = "infra-eng@chanzuckerberg.com"
    service = var.service_name
    project = var.app_name
    env     = var.env
  }
  aws_ssm_paths     = var.aws_ssm_paths
  wildcard_redirect = "SUBDOMAIN"
}

resource "okta_app_group_assignments" "happy_app" {
  app_id = module.happy_app.app.id
  dynamic "group" {
    for_each = toset([for k, v in data.okta_group.teams : v.id])
    content {
      id = group.value
    }
  }
}

data "okta_group" "teams" {
  for_each = var.teams
  name     = each.value
}
