resource "okta_auth_server" "auth_server" {
  audiences   = [local.label]
  description = "Auth server for ${local.label}"
  name        = local.label
  issuer_mode = "ORG_URL"
}

resource "okta_auth_server_scope" "scope" {
  auth_server_id   = okta_auth_server.auth_server.id
  metadata_publish = "ALL_CLIENTS"
  name             = var.tags.service
  consent          = "IMPLICIT"
}

resource "okta_auth_server_policy" "policy" {
  auth_server_id = okta_auth_server.auth_server.id
  priority       = 1
  name           = "Service account"
  description    = "Only allow the service account's credentials access to this application."
  client_whitelist = [
    okta_app_oauth.app.id,
  ]
}

resource "okta_auth_server_policy_rule" "rule" {
  auth_server_id       = okta_auth_server.auth_server.id
  policy_id            = okta_auth_server_policy.policy.id
  name                 = "Service account client credentials only"
  priority             = 1
  scope_whitelist      = [okta_auth_server_scope.scope.name]
  grant_type_whitelist = ["client_credentials"]
  group_whitelist      = ["EVERYONE"]
}
