data "kubernetes_secret" "integration_secret" {
  metadata {
    name      = "integration-secret"
    namespace = var.k8s_namespace
  }
}

resource "validation_error" "mix_of_internal_and_external_services" {
  condition = length(local.external_services) > 0 && length(local.internal_services) > 0 && var.routing_method == "DOMAIN"
  summary   = "Invalid mix of INTERNAL and EXTERNAL services"
  details   = "With DOMAIN routing, a mix of EXTERNAL and INTERNAL services is not permitted; only EXTERNAL and PRIVATE can be mixed"
}

resource "random_pet" "suffix" {}

locals {
  # kubernetes_secret resource is always marked sensitive, which makes things a little difficult
  # when decoding pieces of the integration secret later. Mark the whole thing as nonsensitive and only
  # output particular fields as sensitive in this modules outputs (for instance, the RDS password)
  secret       = jsondecode(nonsensitive(data.kubernetes_secret.integration_secret.data.integration_secret))
  external_dns = local.secret["external_zone_name"]

  s = { for k, v in var.services : k => merge(v, {
    external_stack_host_match   = try(join(".", [var.stack_name, local.external_dns]), "")
    external_service_host_match = try(join(".", ["${var.stack_name}-${k}", local.external_dns]), "")

    stack_host_match   = try(join(".", [var.stack_name, local.external_dns]), "")
    service_host_match = try(join(".", ["${var.stack_name}-${k}", local.external_dns]), "")
  }) }

  suffix = random_pet.suffix.id

  sd = { for k, v in local.s : k => merge(v, {
    external_host_match = var.routing_method == "CONTEXT" ? v.external_stack_host_match : v.external_service_host_match
    host_match          = var.routing_method == "CONTEXT" ? v.stack_host_match : v.service_host_match
    group_name          = var.routing_method == "CONTEXT" ? "stack-${var.stack_name}-${local.suffix}" : "service-${var.stack_name}-${k}-${local.suffix}"
    service_name        = "${var.stack_name}-${k}"
    service_port        = coalesce(v.service_port, v.port)
    bypasses = (v.service_type == "INTERNAL" ?
      merge({
        // by default, add an options bypass since this is used a lot by developers out of the gate
        options = {
          paths   = toset(["/*"])
          methods = toset(["OPTIONS"])
        } },
        v.bypasses
      ) :
    {})
  }) }

  // calculate the highest priority and build off of that
  highest_priority = max(0, [for k, v in local.sd : v.priority]...)
  // find all the services that used the default 0 priority
  unprioritized_service_definitions = [for k, v in local.sd : { (k) = v } if v.priority == 0]
  // make a range starting from the highest and going for every unprioritized service
  // ex: if the highest priority was 4 and was have 2 unprioritized services, they will be assigned priority 5 and 6
  priority_split = range(local.highest_priority + 1, local.highest_priority + length(local.unprioritized_service_definitions) + 1)
  // and reassign them
  reprioritized_service_definitions = { for p, def in zipmap(local.priority_split, local.unprioritized_service_definitions) : keys(def)[0] => merge(def[keys(def)[0]], { priority = p }) }
  prioritized_service_definitions   = { for k, v in local.sd : k => v if v.priority != 0 }
  service_definitions               = merge(local.prioritized_service_definitions, local.reprioritized_service_definitions)

  external_services = [for v in var.services : v if v.service_type == "EXTERNAL"]
  internal_services = [for v in var.services : v if v.service_type == "INTERNAL"]

  service_ecrs = { for k, v in module.services : k => v.ecr.repository_url }

  task_definitions = { for k, v in var.tasks : k => merge(v, {
    task_name = "${var.stack_name}-${k}"
    // substitute {service} references in task image with the appropriate ECR repo urls
    image = format(
      replace(v.image, "/{(${join("|", keys(local.service_ecrs))})}/", "%s"),
      [
        for repo in flatten(regexall("{(${join("|", keys(local.service_ecrs))})}", v.image)) :
        lookup(local.service_ecrs, repo)
      ]...
    )
  }) }


  external_endpoints = concat([for k, v in local.service_definitions :
    v.service_type == "EXTERNAL" ?
    {
      "EXTERNAL_${upper(replace(k, "-", "_"))}_ENDPOINT" = try(join("", ["https://", v.external_host_match]), "")
      "${upper(replace(k, "-", "_"))}_ENDPOINT"          = try(join("", ["https://", v.host_match]), "")
    }
    : {
      "EXTERNAL_${upper(replace(k, "-", "_"))}_ENDPOINT" = try(join("", ["https://", v.external_host_match]), "")
      "INTERNAL_${upper(replace(k, "-", "_"))}_ENDPOINT" = try(join("", ["https://", v.host_match]), "")
      "${upper(replace(k, "-", "_"))}_ENDPOINT"          = try(join("", ["https://", v.host_match]), "")
    }
  ])

  private_endpoints = concat([for k, v in local.service_definitions :
    {
      "PRIVATE_${upper(replace(k, "-", "_"))}_ENDPOINT" = "${lower(v.service_scheme)}://${v.service_name}.${var.k8s_namespace}.svc.cluster.local:${v.service_port}"
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

  service_endpoints = merge(local.flat_external_endpoints, local.flat_private_endpoints)

  db_env_vars = merge(flatten(
    [for dbname, dbcongif in local.secret["dbs"] : [
      for varname, value in dbcongif : { upper(replace("${dbname}_${varname}", "/[^a-zA-Z0-9_]/", "_")) : value }
    ]]
  )...)

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

  // each bypass we add to the load balancer adds one rule. In order to space out the rules
  // evenly for all services on the same load balancer, we need to know the max number of
  // bypasses and multiply by that. The "+ 1" at the end accounts for the actual service rule.
  priority_spread = max(0, [for i in local.service_definitions : length(i.bypasses)]...) + 1

  // If WAF information is set, pull it out so we can configure a WAF. Otherwise, ignore
  waf_config       = lookup(local.secret, "waf_config", {})
  regional_waf_arn = lookup(local.waf_config, "arn", null)
}

resource "kubernetes_secret" "oidc_config" {
  metadata {
    name      = local.oidc_config_secret_name
    namespace = var.routing.service_mesh ? "nginx-encrypted-ingress" : var.k8s_namespace
  }

  data = {
    clientID     = try(local.secret["oidc_config"]["client_id"], "")
    clientSecret = try(local.secret["oidc_config"]["client_secret"], "")
  }
}

module "services" {
  for_each = local.service_definitions
  source   = "../happy-service-eks"

  image_tag                        = lookup(var.image_tags, each.key, var.image_tag)
  container_name                   = each.value.name
  stack_name                       = var.stack_name
  desired_count                    = each.value.desired_count
  max_count                        = try(each.value.max_count, each.value.desired_count)
  scaling_cpu_threshold_percentage = each.value.scaling_cpu_threshold_percentage
  memory                           = each.value.memory
  cpu                              = each.value.cpu
  health_check_path                = each.value.health_check_path
  k8s_namespace                    = var.k8s_namespace
  cloud_env                        = local.secret["cloud_env"]
  certificate_arn                  = local.secret["certificate_arn"]
  deployment_stage                 = var.deployment_stage
  service_endpoints                = local.service_endpoints
  aws_iam_policy_json              = each.value.aws_iam_policy_json
  eks_cluster                      = local.secret["eks_cluster"]
  initial_delay_seconds            = each.value.initial_delay_seconds
  period_seconds                   = each.value.period_seconds
  platform_architecture            = each.value.platform_architecture
  sidecars                         = each.value.sidecars
  routing = {
    method              = var.routing_method
    host_match          = each.value.host_match
    group_name          = each.value.group_name
    priority            = each.value.priority * local.priority_spread
    path                = each.value.path
    service_name        = each.value.service_name
    port                = each.value.port
    service_port        = coalesce(each.value.service_port, each.value.port)
    scheme              = each.value.scheme
    service_scheme      = each.value.service_scheme
    success_codes       = each.value.success_codes
    service_type        = each.value.service_type
    service_mesh        = each.value.service_mesh
    allow_mesh_services = each.value.allow_mesh_services
    oidc_config         = local.oidc_config
    bypasses            = each.value.bypasses
    alb                 = each.value.alb
  }

  additional_env_vars                  = merge(local.db_env_vars, var.additional_env_vars, local.stack_configs)
  additional_env_vars_from_config_maps = var.additional_env_vars_from_config_maps
  additional_env_vars_from_secrets     = var.additional_env_vars_from_secrets
  additional_volumes_from_secrets      = var.additional_volumes_from_secrets
  additional_volumes_from_config_maps  = var.additional_volumes_from_config_maps

  tags = local.secret["tags"]

  regional_wafv2_arn = local.regional_waf_arn
}

module "tasks" {
  for_each              = local.task_definitions
  source                = "../happy-task-eks"
  task_name             = each.value.task_name
  image                 = each.value.image
  cpu                   = each.value.cpu
  memory                = each.value.memory
  cmd                   = each.value.cmd
  remote_dev_prefix     = var.stack_prefix
  deployment_stage      = var.deployment_stage
  k8s_namespace         = var.k8s_namespace
  stack_name            = var.stack_name
  platform_architecture = each.value.platform_architecture
}

