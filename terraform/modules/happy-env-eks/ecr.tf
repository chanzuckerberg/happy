
module "ecrs" {
  for_each = var.ecr_repos
  # tflint-ignore: terraform_module_pinned_source
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/ecr-repository?ref=main"

  name              = each.value["name"]
  read_arns         = try(each.value["read_arns"], [])
  write_arns        = try(each.value["write_arns"], [])
  allow_lambda_pull = true

  env     = var.tags.env
  owner   = var.tags.owner
  project = var.tags.project
  service = var.tags.service
}
