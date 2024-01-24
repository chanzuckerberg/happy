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

data "datadog_synthetics_locations" "locations" {}

resource "datadog_synthetics_test" "test_api" {
  for_each = local.synthetics
  type     = "api"
  subtype  = "http"

  request_definition {
    method = "GET"
    url    = each.value
  }
  assertion {
    type     = "statusCode"
    operator = "is"
    target   = "200"
  }
  locations = keys(data.datadog_synthetics_locations.locations.locations)
  options_list {
    tick_every = 900

    retry {
      count    = 2
      interval = 300
    }

    monitor_options {
      renotify_interval = 120
    }
  }
  name    = "A website synthetic for the happy stack ${var.deployment_stage} ${var.stack_name} ${each.key} located at ${each.value}"
  message = "Notify @opsgenie-${local.opsgenie_owner}"
  status  = "live"
  tags    = values(local.secret["tags"])
}
