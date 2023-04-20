# write the oidc provider values that happy api needs to the hapi chamber service so they will be available to happy api
locals {
  param_suffix = replace("${var.tags.service}_${var.tags.env}_${var.tags.project}", "-", "_")
}

module "params" {
  source  = "git@github.com:chanzuckerberg/cztack//aws-ssm-params-writer?ref=tsmith/aws_provider"
  service = "hapi"
  project = "happy"
  # all happy environments (dev, staging, prod) will be utilizing the prod API
  env   = "prod"
  owner = var.tags.owner

  parameters = {
    "oidc_provider_${local.param_suffix}" : "${okta_auth_server.auth_server.issuer}|${local.label}"
  }

  providers = {
    aws = aws.czi-si
  }
}
