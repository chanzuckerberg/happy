locals {
  # this is making an assumption that the health check path is accessible to the internet
  # if this is a non-prod OIDC protected stack, make sure to allow the health check endpoint
  # through the OIDC proxy.
  synthetics = { for k, v in local.service_definitions : v.service_name =>
    var.routing_method == "DOMAIN" ? "https://${v.service_name}.${local.external_dns}${v.health_check_path}" : "https://${var.stack_name}.${local.external_dns}${v.health_check_path}"
    if v.synthetics && (v.service_type == "EXTERNAL" || v.service_type == "INTERNAL")
  }
  opsgenie_owner = "${local.secret["tags"]["project"]}-${local.secret["tags"]["env"]}-${local.secret["tags"]["service"]}"
}

module "datadog_synthetic" {
  for_each         = local.synthetics
  source           = "../happy-datadog-synthetics"
  service_name     = each.key
  synthetic_url    = each.value
  stack_name       = var.stack_name
  deployment_stage = var.deployment_stage
  opsgenie_owner   = local.opsgenie_owner
  tags             = values(local.secret["tags"])
}
