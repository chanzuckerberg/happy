resource "okta_auth_server" "auth_server" {
  audiences   = [local.label]
  description = "Auth server for ${local.label}"
  name        = local.label
  issuer_mode = "ORG_URL"
}

# Need a scope that will open custom claims: https://developer.okta.com/docs/reference/api/oidc/#well-known-oauth-authorization-server
resource "okta_auth_server_scope" "scope" {
  auth_server_id   = okta_auth_server.auth_server.id
  metadata_publish = "ALL_CLIENTS"
  name             = var.tags.service
  consent          = "IMPLICIT"
}

resource "okta_auth_server_policy" "policy" {
  auth_server_id = okta_auth_server.auth_server.id
  priority       = 1
  name           = "Default Policy"
  description    = "Default Policy for your Authorization Server"
  client_whitelist = [
    okta_app_oauth.app.id,
  ]
}

resource "okta_auth_server_policy_rule" "rule" {
  auth_server_id       = okta_auth_server.auth_server.id
  policy_id            = okta_auth_server_policy.policy.id
  name                 = "Default Policy Rule"
  priority             = 1
  scope_whitelist      = ["*"]
  grant_type_whitelist = ["authorization_code"]
  group_whitelist      = ["EVERYONE"]
}
