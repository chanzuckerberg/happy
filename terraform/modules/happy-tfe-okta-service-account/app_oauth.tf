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
  label = "${var.tags.project}-${var.tags.env}-${var.tags.service}-service-account"
}

resource "okta_app_oauth" "app" {
  label                      = local.label
  type                       = "service"
  grant_types                = ["client_credentials"]
  redirect_uris              = []
  login_uri                  = ""
  omit_secret                = false
  response_types             = ["token"]
  login_mode                 = "DISABLED"
  login_scopes               = []
  token_endpoint_auth_method = "private_key_jwt"

  jwks {
    kid = local.jwks.kid
    kty = local.jwks.kty
    e   = local.jwks.e
    n   = local.jwks.n
  }
}
