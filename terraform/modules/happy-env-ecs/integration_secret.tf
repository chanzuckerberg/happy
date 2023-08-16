locals {
  standard_secrets = {
    kind               = "ecs"
    zone_id            = var.base_zone
    external_zone_name = local.env_domain
    internal_zone_name = try(module.ecs-multi-domain-oauth-proxy[0].proxy.zones.internal_name, "")
    cluster_arn        = module.ecs-cluster.arn
    ecs_execution_role = aws_iam_role.task_execution_role.arn
    cloud_env          = var.cloud-env
    vpc_id             = var.cloud-env.vpc_id          # Deprecated, use "cloud_env" value directly
    private_subnets    = var.cloud-env.private_subnets # Deprecated, use "cloud_env" value directly
    public_subnets     = var.cloud-env.public_subnets  # Deprecated, use "cloud_env" value directly
    tags               = var.tags
    security_groups    = [aws_security_group.happy_env_sg.id]
    # NOTE - this is the old and busted, still here for reverse-compatibility with older happy envs
    batch_queues = merge({ for name, batch in module.batch : name => { "queue_arn" : batch.batch.queue, "role_arn" : batch.role.arn } },
    { for name, batch in module.batch-swipe : name => { "queue_arn" : batch.batch.envs["EC2"].queue_arn, "role_arn" : batch.role.arn } })

    # NOTE - this is the newer hotness - SWIPE ONLY, still here for reverse-compatibility with older happy envs
    batch_envs = { for name, batch in module.batch-swipe : name => { "envs" : batch.batch.envs, "role" : batch.role } }

    # NOTE - these are the newest outputs - SWIPE from the module, please use these
    swipe_sfn_arns                                = { for name, swipe in module.swipe : name => swipe.sfn_arns }
    swipe_sfn_notification_queue_arns             = { for name, swipe in module.swipe : name => swipe.sfn_notification_queue_arns }
    swipe_sfn_notification_dead_letter_queue_arns = { for name, swipe in module.swipe : name => swipe.sfn_notification_dead_letter_queue_arns }

    s3_buckets = { for name, bucket in module.s3_bucket : name => { "name" : bucket.name, "arn" : bucket.arn } }
    public_albs = {
      for name, lb in aws_lb.lb-public :
      name => {
        "arn" : lb.arn,
        "dns_name" : lb.dns_name,
        "zone_id" : lb.zone_id,
        "listener_arn" : aws_lb_listener.public-https[name].arn
      }
    }

    private_albs = {
      for name, lb in aws_lb.lb-private :
      name => {
        "arn" : lb.arn,
        "dns_name" : lb.dns_name,
        "zone_id" : lb.zone_id,
        "listener_arn" : aws_lb_listener.private-lb-listener[name].arn
      }
    }

    ecrs = { for name, ecr in module.ecrs : name => { "url" : ecr.repository_url } }
    dbs = {
      for name, db in module.db :
      name => {
        "database_user" : db.master_username,
        "database_password" : db.master_password,
        "database_host" : db.endpoint,
        "database_name" : db.database_name,
        "database_port" : db.port
      }
    }
    hapi_config = {
      base_url        = var.hapi_base_url
      oidc_issuer     = module.happy_service_account.oidc_config.client_id
      oidc_authz_id   = module.happy_service_account.oidc_config.authz_id
      scope           = module.happy_service_account.oidc_config.scope
      kms_key_id      = module.happy_service_account.kms_key_id
      assume_role_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/tfe-si"
    }

    dynamo_locktable_name = aws_dynamodb_table.locks.id
    datadog_api_key       = var.datadog_api_key
  }

  # TODO: this only works if all additional_secrets values are maps!
  merged_secrets = { for key, value in var.additional_secrets : key => merge(lookup(local.standard_secrets, key, {}), value) }
  secret_string = merge(
    local.standard_secrets,
    local.merged_secrets,
  )
}

resource "aws_secretsmanager_secret" "happy_env_secret" {
  name = "happy/env-${var.name}-config"
  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "happy_env_secret_version" {
  secret_id     = aws_secretsmanager_secret.happy_env_secret.id
  secret_string = jsonencode(local.secret_string)
}
