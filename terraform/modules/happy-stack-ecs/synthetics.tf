locals {
  opsgenie_owner = "${local.secret["tags"]["project"]}-${local.secret["tags"]["env"]}-${local.secret["tags"]["service"]}"
  url            = "https://${local.fqdn}"
}

data "datadog_synthetics_locations" "locations" {}

resource "datadog_synthetics_test" "test_api" {
  type    = "api"
  subtype = "http"
  request_definition {
    method = "GET"
    url    = local.url
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
  name    = "A website synthetic for the happy stack ${var.deployment_stage} ${var.stack_name} ${var.app_name} located at ${local.url}"
  message = "Notify @opsgenie-${local.opsgenie_owner}"
  status  = "live"
  tags    = values(local.secret["tags"])
}
