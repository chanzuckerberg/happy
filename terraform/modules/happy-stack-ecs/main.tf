data "aws_secretsmanager_secret_version" "config" {
  secret_id = var.happy_config_secret
}

locals {
  secret             = jsondecode(nonsensitive(data.aws_secretsmanager_secret_version.config.secret_string))
  vpc_id             = local.secret["vpc_id"]
  cloud_env           = local.secret["cloud_env"]
  security_groups    = local.secret["security_groups"]
  zone               = local.secret["zone_id"]
  cluster            = local.secret["cluster_arn"]
  ecs_execution_role = lookup(local.secret, "ecs_execution_role", "")

  image = join(":", [local.secret["ecrs"][var.app_name]["url"], lookup(var.image_tags, var.app_name, var.image_tag)])

  external_dns = local.secret["external_zone_name"]
  internal_dns = local.secret["internal_zone_name"]
  dns_zone_id  = local.secret["zone_id"]

  alb_key      = var.require_okta ? "private_albs" : "public_albs"
  listener_arn = local.secret[local.alb_key][var.app_name]["listener_arn"]
  alb_zone     = local.secret[local.alb_key][var.app_name]["zone_id"]
  alb_dns      = local.secret[local.alb_key][var.app_name]["dns_name"]

  ecs_role_arn  = local.secret["service_roles"]["ecs_role"]
  ecs_role_name = element(split("/", local.secret["service_roles"]["ecs_role"]), length(split("/", local.secret["service_roles"]["ecs_role"])) - 1)
  url           = try(join("", ["https://", module.dns[0].dns_prefix, ".", local.external_dns]), var.url)

  db_env_vars = flatten([
    for dbname, dbcongif in local.secret["dbs"] :
    [
      for varname, value in dbcongif :
      {
        "name" : upper(replace("${dbname}_${varname}", "/[^a-zA-Z0-9_]/", "_")),
        "value" : value
      }
    ]
  ])
}

# TODO: remove this
module "dns" {
  source                = "../happy-dns-ecs"
  custom_stack_name     = var.stack_name
  app_name              = var.app_name
  alb_dns               = local.alb_dns
  canonical_hosted_zone = local.alb_zone
  zone                  = local.internal_dns
  tags                  = local.secret["tags"]
}

module "service" {
  for_each            = var.services
  source              = "../happy-service-ecs"

  service_name        = each.value.name
  service_type        = each.value.service_type
  desired_count       = each.value.desired_count
  service_port        = each.value.port
  memory              = each.value.memory
  cpu                 = each.value.cpu
  health_check_path   = each.value.health_check_path

  execution_role      = local.ecs_execution_role
  custom_stack_name   = var.stack_name
  app_name            = var.app_name
  image               = local.image
  cluster             = local.cluster
  listener            = local.listener_arn
  security_groups     = local.security_groups
  task_role           = { arn : local.ecs_role_arn, name : local.ecs_role_name }
  deployment_stage    = var.deployment_stage
  host_match          = "${var.stack_name}-${each.value.name}.${local.external_dns}"
  priority            = var.priority
  launch_type         = var.launch_type
  additional_env_vars = merge(var.additional_env_vars, local.db_env_vars)
  tags                = local.secret["tags"]
}
