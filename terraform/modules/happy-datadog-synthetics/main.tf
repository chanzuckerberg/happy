resource "datadog_synthetics_test" "test_api" {
  type    = "api"
  subtype = "http"
  request_definition {
    method = "GET"
    url    = var.synthetic_url
  }
  assertion {
    type     = "statusCode"
    operator = "is"
    target   = "200"
  }
  locations = ["aws:us-east-1", "aws:us-west-2"]
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
  name    = "A website synthetic for the happy stack ${var.deployment_stage} ${var.stack_name} ${var.service_name} located at ${var.synthetic_url}"
  message = "Notify @opsgenie-${var.opsgenie_owner}"
  status  = "live"
  tags    = var.tags
}
