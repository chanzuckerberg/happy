module "ecr" {
  source = "git@github.com:chanzuckerberg/cztack//aws-ecr-repo?ref=v0.59.0"

  name           = "${var.stack_name}/${local.tags.env}/${var.container_name}"
  force_delete   = true
  tag_mutability = var.tag_mutability
  scan_on_push   = var.scan_on_push
  tags           = var.tags
}

