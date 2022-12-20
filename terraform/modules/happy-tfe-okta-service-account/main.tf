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
  label = "${var.service_name}-${var.app_name}-${var.env}-service-account"
}

module "happy_app" {
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth?ref=heathj/jwks"

  okta = {
    label         = local.label
    redirect_uris = []
    login_uri     = ""
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
    env     = var.env
  }
  aws_ssm_paths = var.aws_ssm_paths
  jwks          = local.jwks
  # we set at least one role so that an authorization server is created
  # the authorization server is required for creating a service account
  rbac_role_mapping = merge({
    base : []
  }, var.rbac_role_mapping)
}
