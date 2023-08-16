module "ecrs" {
  for_each = var.ecr_repos
  source   = "git@github.com:chanzuckerberg/cztack//aws-ecr-repo?ref=v0.56.2"

  name       = each.value["name"]
  read_arns  = each.value["read_arns"]
  write_arns = each.value["write_arns"]

  allow_lambda_pull = true
  tags              = var.tags
}
