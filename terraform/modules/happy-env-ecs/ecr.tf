module "ecrs" {
  for_each = var.ecr_repos
  source   = "git@github.com:chanzuckerberg/cztack//aws-ecr-repo?ref=v0.59.0"

  name       = each.value.name
  read_arns  = each.value.read_arns
  write_arns = each.value.write_arns

  allow_lambda_pull = true
  tag_mutability    = each.value.tag_mutability
  scan_on_push      = each.value.scan_on_push
  tags              = var.tags
}

move_resource "ecr" {
  from = module.ecr
  to   = module.ecrs
}