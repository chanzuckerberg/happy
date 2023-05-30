module "ecr" {
  source = "git@github.com:chanzuckerberg/cztack//aws-ecr-repo?ref=v0.56.2"

  name         = "${var.stack_name}/${local.tags.env}/${var.container_name}"
  force_delete = true
  tags         = var.tags
}
// TODO: enable ECR scanning
