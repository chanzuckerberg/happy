resource "aws_kms_key" "service_user" {
  description              = local.label
  key_usage                = "SIGN_VERIFY"
  customer_master_key_spec = "RSA_4096"
}
data "aws_kms_public_key" "service_user" {
  key_id = aws_kms_key.service_user.key_id
}
data "jwks_from_key" "jwks" {
  key = data.aws_kms_public_key.service_user.public_key_pem
  kid = aws_kms_key.service_user.key_id
}
locals {
  jwks  = jsondecode(data.jwks_from_key.jwks.jwks)
  label = "${var.service_name}-${var.app_name}-service-account"
}

module "happy_app" {
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth?ref=heathj/jwks"

  okta = {
    label         = local.label
    redirect_uris = concat(["https://oauth.${var.app_name}.si.czi.technology/oauth2/callback"], var.redirect_uris)
    login_uri     = var.login_uri == "" ? "https://oauth.${var.app_name}.si.czi.technology" : var.login_uri
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
    env     = var.app_name
  }
  aws_ssm_paths = var.aws_ssm_paths
  jwks          = local.jwks
  # we set at least on role so that an authorization server is created
  rbac_role_mapping = merge({
    base : []
  }, var.rbac_role_mapping)
}

resource "okta_app_group_assignments" "happy_app" {
  app_id    = module.happy_app.app.id
  group_ids = var.teams
}
