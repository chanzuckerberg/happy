output "kms_key_id" {
  value = aws_kms_key.service_user.key_id
}

output "app" {
  value = {
    id   = okta_app_oauth.app.id
    name = okta_app_oauth.app.name
  }
}

output "authz_server" {
  value = okta_auth_server.auth_server
}

locals {
  auth_server_id = split("oauth2/", okta_auth_server.auth_server.issuer)[1]
}

output "oidc_config" {
  value = {
    client_id     = okta_app_oauth.app.client_id
    client_secret = okta_app_oauth.app.client_secret
    idp_url       = okta_auth_server.auth_server.issuer
    authz_id      = local.auth_server_id
    scope         = okta_auth_server_scope.scope.name
    config_uri    = "https://${okta_app_oauth.app.client_id}:${okta_app_oauth.app.client_secret}@${var.okta_tenant}.okta.com/oauth2/${local.auth_server_id}"
  }
}
