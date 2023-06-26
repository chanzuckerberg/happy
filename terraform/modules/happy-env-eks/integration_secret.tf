data "aws_caller_identity" "current" {}

locals {
  standard_secrets = {
    kind               = "k8s"
    vpc_id             = var.cloud-env.vpc_id
    zone_id            = var.base_zone_id
    external_zone_name = data.aws_route53_zone.base_zone.name

    cloud_env       = var.cloud-env
    eks_cluster     = var.eks-cluster
    tags            = var.tags
    certificate_arn = module.cert.arn
    ci_roles        = var.github_actions_roles
    ecrs            = { for name, ecr in module.ecrs : name => { "url" : ecr.repository_url } }
    dbs = {
      for name, db in module.dbs :
      name => {
        "database_user" : db.master_username,
        "database_password" : db.master_password,
        "database_host" : db.endpoint,
        "database_name" : db.database_name,
        "database_port" : db.port,
      }
    }
    oidc_config = module.happy_okta_app.oidc_config
    hapi_config = {
      base_url        = var.hapi_base_url
      oidc_issuer     = module.happy_service_account.oidc_config.client_id
      oidc_authz_id   = module.happy_service_account.oidc_config.authz_id
      scope           = module.happy_service_account.oidc_config.scope
      kms_key_id      = module.happy_service_account.kms_key_id
      assume_role_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/tfe-si"
    }
    dynamo_locktable_name = aws_dynamodb_table.locks.id
  }

  merged_secrets = { for key, value in var.additional_secrets : key => merge(lookup(local.standard_secrets, key, {}), value) }
  secret_string = merge(
    local.standard_secrets,
    local.merged_secrets
  )
}

resource "kubernetes_secret" "happy_env_secret" {
  metadata {
    name      = "integration-secret"
    namespace = kubernetes_namespace.happy.id
  }
  data = {
    "integration_secret" = jsonencode(local.secret_string)
  }
}
