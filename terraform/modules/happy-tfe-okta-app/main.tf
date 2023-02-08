locals {
  base_domain = (var.base_domain == "si.czi.technology" ?
    "${var.app_name}.${var.env}.${var.base_domain}" :
  var.base_domain)

  base_domain_redirect_uris = [for route in var.redirect_uri_base_domain_routes : "https://*.${local.base_domain}${route}"]
}
module "happy_app" {
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth-head?ref=v0.249.0"

  okta = {
    label         = "*.${local.base_domain}"
    redirect_uris = concat(local.base_domain_redirect_uris, var.redirect_uris)
    login_uri     = var.login_uri == "" ? "https://oauth.${local.base_domain}" : var.login_uri
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
    for_each = merge([for x, y in data.okta_groups.teams : { for k, v in y.groups : v.name => v }]...)
    content {
      id = group.value.id
    }
  }
}

data "okta_groups" "teams" {
  for_each = var.teams
  q        = each.value
}
