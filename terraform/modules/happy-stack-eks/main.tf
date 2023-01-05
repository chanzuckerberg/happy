data "kubernetes_secret" "integration_secret" {
  metadata {
    name      = "integration-secret"
    namespace = var.k8s_namespace
  }
}

locals {
  # kubernetes_secret resource is always marked sensitive, which makes things a little difficult
  # when decoding pieces of the integration secret later. Mark the whole thing as nonsensitive and only
  # output particular fields as sensitive in this modules outputs (for instance, the RDS password)
  secret       = jsondecode(nonsensitive(data.kubernetes_secret.integration_secret.data.integration_secret))
  external_dns = local.secret["external_zone_name"]
  internal_dns = local.secret["internal_zone_name"]

  service_definitions = { for k, v in var.services : k => merge(v, {
    host_match   = v.service_type == "INTERNAL" ? try(join(".", ["${var.stack_name}-${k}", "internal", local.external_dns]), "") : try(join(".", ["${var.stack_name}-${k}", local.external_dns]), "")
    service_name = "${var.stack_name}-${k}"
  }) }

  task_definitions = { for k, v in var.tasks : k => merge(v, {
    task_name = "${var.stack_name}-${k}"
  }) }

  external_endpoints = concat([for k, v in local.service_definitions :
    v.service_type == "EXTERNAL" ?
    {
      "EXTERNAL_${upper(k)}_ENDPOINT" = try(join("", ["https://", v.service_name, "."]), "")
    }
    : {
      "INTERNAL_${upper(k)}_ENDPOINT" = try(join("", ["https://", v.service_name, "."]), "")
    }
  ])

  private_endpoints = concat([for k, v in local.service_definitions :
    {
      "PRIVATE_${upper(k)}_ENDPOINT" = "http://${v.service_name}.${var.k8s_namespace}.svc.cluster.local:${v.port}"
    }
  ])

  flat_external_endpoints = zipmap(
    flatten(
      [for item in local.external_endpoints : keys(item)]
    ),
    flatten(
      [for item in local.external_endpoints : values(item)]
    )
  )

  flat_private_endpoints = zipmap(
    flatten(
      [for item in local.private_endpoints : keys(item)]
    ),
    flatten(
      [for item in local.private_endpoints : values(item)]
    )
  )

  service_endpoints = merge(local.flat_external_endpoints, {})

  db_env_vars = merge(flatten(
    [for dbname, dbcongif in local.secret["dbs"] : [
      for varname, value in dbcongif : { upper(replace("${dbname}_${varname}", "/[^a-zA-Z0-9_]/", "_")) : value }
    ]]
  )...)
}

module "services" {
  for_each              = local.service_definitions
  source                = "../happy-service-eks"
  image                 = join(":", [local.secret["ecrs"][each.key]["url"], lookup(var.image_tags, each.key, var.image_tag)])
  container_name        = each.value.name
  stack_name            = var.stack_name
  desired_count         = each.value.desired_count
  service_name          = each.value.service_name
  service_type          = each.value.service_type
  memory                = each.value.memory
  cpu                   = each.value.cpu
  health_check_path     = each.value.health_check_path
  k8s_namespace         = var.k8s_namespace
  cloud_env             = local.secret["cloud_env"]
  certificate_arn       = local.secret["certificate_arn"]
  oauth_certificate_arn = local.secret["oauth_certificate_arn"]
  host_match            = each.value.host_match
  service_port          = each.value.port
  deployment_stage      = var.deployment_stage
  service_endpoints     = local.service_endpoints
  aws_iam_policy_json   = each.value.aws_iam_policy_json
  eks_cluster           = local.secret["eks_cluster"]
  additional_env_vars   = local.db_env_vars
}

module "tasks" {
  for_each          = local.task_definitions
  source            = "../happy-task-eks"
  task_name         = each.value.task_name
  image             = each.value.image
  cpu               = each.value.cpu
  memory            = each.value.memory
  cmd               = each.value.cmd
  remote_dev_prefix = var.stack_prefix
  deployment_stage  = var.deployment_stage
  k8s_namespace     = var.k8s_namespace
  stack_name        = var.stack_name
}
