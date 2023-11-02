data "aws_region" "current" {}
resource "random_pet" "suffix" {}

data "kubernetes_secret_v1" "integration_secret" {
  metadata {
    name      = "integration-secret"
    namespace = var.k8s_namespace
  }
}

locals {
  suffix = random_pet.suffix.id

  secret          = jsondecode(nonsensitive(data.kubernetes_secret_v1.integration_secret.data.integration_secret))
  external_dns    = local.secret["external_zone_name"]
  certificate_arn = local.secret["certificate_arn"]

  s = { for k, v in var.services : k => merge(v, {
    external_stack_host_match   = try(join(".", [var.stack_name, local.external_dns]), "")
    external_service_host_match = try(join(".", ["${var.stack_name}-${k}", local.external_dns]), "")

    stack_host_match   = try(join(".", [var.stack_name, local.external_dns]), "")
    service_host_match = try(join(".", ["${var.stack_name}-${k}", local.external_dns]), "")
  }) }

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

  # Enriched Services
  service_endpoints = merge(local.flat_external_endpoints, local.flat_private_endpoints)

  # Enriched Tasks
  task_definitions = {
    for k, v in var.tasks : k => merge(v, {
      task_name = "${var.stack_name}-${k}"
      // substitute {service} references in task image with the appropriate ECR repo urls
      image = format(
        replace(v.image, "/{(${join("|", keys(local.service_ecrs))})}/", "%s"),
        [
          for repo in flatten(regexall("{(${join("|", keys(local.service_ecrs))})}", v.image)) :
          lookup(local.service_ecrs, repo)
        ]...
      )
    })
  }

  service_ecrs = { for k, v in module.ecr : k => v.repository_url }
  tags         = local.secret["tags"]
  cloud_env    = local.secret["cloud_env"]

  # WAF
  waf_config       = lookup(local.secret, "waf_config", {})
  regional_waf_arn = lookup(local.waf_config, "arn", null)

  # OIDC
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
