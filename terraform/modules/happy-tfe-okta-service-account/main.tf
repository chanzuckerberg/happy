# generate a signing key in terraform so that this
# can be fully automatic and we don't have to worry about
# passing keys around.
resource "aws_kms_key" "service_user" {
  description              = local.label
  key_usage                = "SIGN_VERIFY"
  customer_master_key_spec = "RSA_4096"
}
resource "aws_kms_alias" "service_user" {
  name_prefix   = "alias/${local.label}"
  target_key_id = aws_kms_key.service_user.key_id
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

# exact configuration needed to make a service user account
# the output of this module will contain the KMS key's ID which can be used
# to sign JWT for this service user.
module "service_user" {
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth?ref=v0.255.0"

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
  grant_type_whitelist       = ["client_credentials"]

  tags = {
    owner   = "infra-eng@chanzuckerberg.com"
    service = "${var.service_name}-oauth"
    project = var.app_name
    env     = var.env
  }
  aws_ssm_paths = var.aws_ssm_paths
  jwks          = toset([local.jwks])
  # we set at least one role so that an authorization server is created
  # the authorization server is required for creating a service account
  rbac_role_mapping = merge({
    base : []
  }, var.rbac_role_mapping)
}
