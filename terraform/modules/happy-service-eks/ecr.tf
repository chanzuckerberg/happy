module "ecr" {
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/ecr-repository?ref=main"

  name = "${var.stack_name}-${var.container_name}"

  env     = local.tags.env
  owner   = local.tags.owner
  project = local.tags.project
  service = local.tags.service
}
