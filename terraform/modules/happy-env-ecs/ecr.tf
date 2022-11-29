
module "ecr" {
  for_each = var.ecr_repos
  source   = "../ecr-repository"
  name     = each.value["name"]

  env     = local.env
  owner   = local.owner
  project = local.project
  service = local.component

  # ARN's that can read/write this repo.
  read_arns  = try(each.value["read_arns"], [])
  write_arns = try(each.value["write_arns"], [])

  allow_lambda_pull = true
}
