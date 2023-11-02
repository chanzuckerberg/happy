locals {
  oidc_config_secret_name = "${var.stack_name}-oidc-config"
  issuer_domain           = try(local.secret["oidc_config"]["idp_url"], "todofindissuer.com")
  issuer_url              = "https://${local.issuer_domain}"
  oidc_config = {
    issuer                = local.issuer_url
    authorizationEndpoint = "${local.issuer_url}/oauth2/v1/authorize"
    tokenEndpoint         = "${local.issuer_url}/oauth2/v1/token"
    userInfoEndpoint      = "${local.issuer_url}/oauth2/v1/userinfo"
    secretName            = local.oidc_config_secret_name
  }
}

resource "kubernetes_secret" "oidc_config" {
  metadata {
    name      = local.oidc_config_secret_name
    namespace = var.enable_service_mesh ? "nginx-encrypted-ingress" : var.k8s_namespace
  }

  data = {
    clientID     = try(local.secret["oidc_config"]["client_id"], "")
    clientSecret = try(local.secret["oidc_config"]["client_secret"], "")
  }
}