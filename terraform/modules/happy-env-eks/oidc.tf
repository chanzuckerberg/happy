module "happy_okta_app" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-tfe-okta-app?ref=happy-tfe-okta-app-v3.1.0"

  app_name = var.tags.project
  env      = var.tags.env
  # backward compatibility
  # todo: remove var.okta_teams for var.oidc_config.teams
  teams                      = coalesce(var.okta_teams, var.oidc_config.teams)
  base_domain                = data.aws_route53_zone.base_zone.name
  redirect_uris              = var.oidc_config.redirect_uris
  login_uri                  = var.oidc_config.login_uri
  grant_types                = var.oidc_config.grant_types
  app_type                   = var.oidc_config.app_type
  token_endpoint_auth_method = var.oidc_config.token_endpoint_auth_method
}
