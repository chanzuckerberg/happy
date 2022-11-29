locals {
  # all resources associated this is stack should have these tags
  tags = {
    happy-env           = var.deployment_stage,
    happy-stack-name    = var.custom_stack_name,
    happy-service-name  = var.app_name,
    happy-region        = data.aws_region.current.name,
    happy-image         = var.image,
    happy-service-type  = var.service_type
    happy-last-applied  = timestamp(),
    happy-last-git-hash = data.external.git_sha.result.sha
    happy-repo          = data.external.git_repo.result
  }
}

data "external" "git_sha" {
  program = [
    "git",
    "log",
    "--pretty=format:{ \"sha\": \"%H\" }",
    "-1",
    "HEAD"
  ]
}

data "external" "git_repo" {
  program = [
    "git",
    "remote",
    "get-url",
    "origin"
  ]
}
