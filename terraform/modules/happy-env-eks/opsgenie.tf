module "ops-genie" {
  # tflint-ignore: terraform_module_pinned_source
  source               = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/ops-genie-service?ref=main"
  service_name         = "${var.tags.project}-${var.tags.env}-${var.tags.service} Infrastructure Monitoring"
  api_integration_name = "${var.tags.project}-${var.tags.env}-${var.tags.service}"
  api_integration_type = "Datadog"
  datadog_service_name = "${var.tags.project}-${var.tags.env}-${var.tags.service}"
  owner_team           = var.ops_genie_owner_team
  opsgenie_tags        = sort(values(var.tags))
}
