module "ecr" {
  source   = "git@github.com:chanzuckerberg/cztack//aws-ecr-repo?ref=v0.59.0"
  for_each = local.service_definitions

  name           = "${var.stack_name}/${local.tags.env}/${each.value.name}"
  force_delete   = true
  tag_mutability = each.value.tag_mutability
  scan_on_push   = each.value.scan_on_push
  tags           = local.secret["tags"]
}