data "aws_secretsmanager_secret_version" "config" {
  secret_id = var.happy_config_secret
}

locals {
  secret  = jsondecode(data.aws_secretsmanager_secret_version.config.secret_string)
  alb_key = var.require_okta ? "private_albs" : "public_albs"

  app_name              = var.app_name
  custom_stack_name     = var.stack_name
  priority              = var.priority
  deployment_stage      = var.deployment_stage
  remote_dev_prefix     = var.stack_prefix
  wait_for_steady_state = var.wait_for_steady_state

  vpc_id             = local.secret["cloud_env"]["vpc_id"]
  subnets            = local.secret["cloud_env"]["private_subnets"]
  security_groups    = local.secret["security_groups"]
  cluster            = local.secret["cluster_arn"]
  ecs_execution_role = lookup(local.secret, "ecs_execution_role", "")
  datadog_api_key    = try(local.secret["datadog_api_key"], "")

  image = join(":", [local.secret["ecrs"][local.app_name]["url"], lookup(var.image_tags, local.app_name, var.image_tag)])

  external_dns = local.secret["external_zone_name"]
  internal_dns = local.secret["internal_zone_name"]

  dns_prefix = local.custom_stack_name == local.app_name ? local.app_name : "${local.custom_stack_name}-${local.app_name}"
  fqdn       = join(".", [local.dns_prefix, local.external_dns])

  listener_arn = local.secret[local.alb_key][local.app_name]["listener_arn"]
  alb_zone     = local.secret[local.alb_key][local.app_name]["zone_id"]
  alb_dns      = local.secret[local.alb_key][local.app_name]["dns_name"]

  ecs_role_arn  = local.secret["service_roles"]["ecs_role"]
  ecs_role_name = element(split("/", local.secret["service_roles"]["ecs_role"]), length(split("/", local.secret["service_roles"]["ecs_role"])) - 1)

  stack_resource_prefix = local.app_name

  db_env_vars = merge(flatten(
    [for dbname, dbcongif in local.secret["dbs"] : [
      for varname, value in dbcongif : { upper(replace("${dbname}_${varname}", "/[^a-zA-Z0-9_]/", "_")) : value }
    ]]
  )...)
}

module "dns" {
  count                 = var.require_okta ? 1 : 0
  source                = "../happy-dns-ecs"
  dns_prefix            = local.dns_prefix
  alb_dns               = local.alb_dns
  canonical_hosted_zone = local.alb_zone
  zone                  = local.internal_dns
}

module "service" {
  source                = "../happy-service-ecs"
  stack_resource_prefix = local.stack_resource_prefix
  execution_role        = local.ecs_execution_role
  memory                = var.memory
  cpu                   = var.cpu
  custom_stack_name     = local.custom_stack_name
  app_name              = local.app_name
  vpc                   = local.vpc_id
  image                 = local.image
  cluster               = local.cluster
  desired_count         = var.desired_count
  listener              = local.listener_arn
  subnets               = local.subnets
  security_groups       = local.security_groups
  task_role             = { arn : local.ecs_role_arn, name : local.ecs_role_name }
  service_port          = var.service_port
  deployment_stage      = local.deployment_stage
  host_match            = local.fqdn
  priority              = local.priority
  remote_dev_prefix     = local.remote_dev_prefix
  wait_for_steady_state = local.wait_for_steady_state
  launch_type           = var.launch_type
  additional_env_vars   = merge(local.db_env_vars, local.stack_configs)
  chamber_service       = var.chamber_service
  tags                  = local.tags
  datadog_api_key       = local.datadog_api_key
  fail_fast             = var.fail_fast
}

