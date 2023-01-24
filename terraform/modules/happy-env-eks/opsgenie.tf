locals {
  opsgenie_tags = [
    var.tags.env,
    var.tags.service,
    var.tags.project,
    var.tags.owner
  ]
  opsgenie_dd_service = "${var.tags.project}-${var.tags.env}-${var.tags.service} Infrastructure Monitoring"
}

module "ops-genie" {
  source               = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/ops-genie-service?ref=main"
  service_name         = local.opsgenie_dd_service
  api_integration_name = "${var.tags.project}-${var.tags.env}-${var.tags.service}"
  api_integration_type = "Datadog"
  datadog_service_name = "${var.tags.project}-${var.tags.env}-${var.tags.service}"
  owner_team           = var.ops_genie_owner_team
  opsgenie_tags        = var.tags
}
