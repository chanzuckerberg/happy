
module "params" {
  source  = "github.com/chanzuckerberg/cztack//aws-ssm-params-writer?ref=v0.43.1"
  project = var.happy_namespace.project
  env     = var.happy_namespace.env
  service = var.happy_namespace.service
  owner   = var.tags.owner

  parameters = {
    (var.aws_ssm_paths.client_id)     = okta_app_oauth.app.client_id
    (var.aws_ssm_paths.client_secret) = (okta_app_oauth.app.client_secret == null || okta_app_oauth.app.client_secret == "") ? local.default_client_secret : okta_app_oauth.app.client_secret
    # If RBAC is defined, use custom issuer (without https:// prefix)
    # If not, use the default value
    (var.aws_ssm_paths.okta_idp_url) = (local.should_create_auth_server == 1) ? local.auth_server_issuer : local.default_issuer


    // see https://registry.terraform.io/providers/oktadeveloper/okta/latest/docs/resources/auth_server#issuer
    (var.aws_ssm_paths.config_uri) = "https://${okta_app_oauth.app.client_id}:${okta_app_oauth.app.client_secret}@${var.okta.tenant}.okta.com/oauth2/${local.auth_server_id}"
  }
}
