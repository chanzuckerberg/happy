# write the oidc provider values that happy api needs to the hapi chamber service so they will be available to happy api
module "params" {
  source  = "github.com/chanzuckerberg/cztack//aws-ssm-params-writer?ref=v0.43.1"
  service = "happy"
  project = "hapi"
  env     = var.tags.env
  owner   = var.tags.owner

  parameters = {
    "oidc_provider_${var.tags.service}_${var.tags.env}_${var.tags.project}" : "${okta_auth_server.auth_server.issuer}|${okta_auth_server.auth_server.audiences[0]}"
  }
}
