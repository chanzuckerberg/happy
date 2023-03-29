# write the oidc provider values that happy api needs to the hapi chamber service so they will be available to happy api
locals {
  param_suffix = replace("${var.tags.service}_${var.tags.env}_${var.tags.project}", "-", "_")
}

module "params" {
  source  = "github.com/chanzuckerberg/cztack//aws-ssm-params-writer?ref=v0.43.1"
  service = "hapi"
  project = "happy"
  env     = var.tags.env
  owner   = var.tags.owner

  parameters = {
    "oidc_provider_${local.param_suffix}" : "${okta_auth_server.auth_server.issuer}|${local.label}"
  }

  providers = {
    aws = aws.czi-si
  }
}
