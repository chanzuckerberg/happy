locals {
  # this is making an assumption that the health check path is accessible to the internet
  # if this is a non-prod OIDC protected stack, make sure to allow the health check endpoint
  # through the OIDC proxy.
  base_synthetics = { for k, v in local.service_definitions : v.service_name =>
    var.routing_method == "DOMAIN" ? "https://${v.service_name}.${local.external_dns}${v.health_check_path}" : "https://${var.stack_name}.${local.external_dns}${v.health_check_path}"
    if v.synthetics && (v.service_type == "EXTERNAL" || v.service_type == "INTERNAL")
  }
  # If you are using a custom domain, you can add additional_hostnames. This will make only 1 synthetic per additonal_hostname per service
  # and no synthetics to the internal URLs created for the services.
  additional_hosts_synthetics = merge([for k, v in local.service_definitions :
    { for domain in var.additional_hostnames : v.service_name => "https://${domain}${v.health_check_path}" }
  ]...)


  synthetics     = length(var.additional_hostnames) == 0 ? local.base_synthetics : local.additional_hosts_synthetics
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
