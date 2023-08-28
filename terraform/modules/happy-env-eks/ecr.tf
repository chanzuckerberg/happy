
module "ecrs" {
  for_each = var.ecr_repos
  source   = "git@github.com:chanzuckerberg/cztack//aws-ecr-repo?ref=v0.59.0"

  name       = each.value["name"]
  read_arns  = each.value["read_arns"]
  write_arns = each.value["write_arns"]

  tag_mutability = each.value["tag_mutability"] == null ? true : each.value["tag_mutability"]
  scan_on_push   = each.value["scan_on_push"] == null ? false : each.value["tagscan_on_push_mutability"]

  tags = var.tags
}
