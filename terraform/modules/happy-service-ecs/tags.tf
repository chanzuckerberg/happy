locals {
  # all resources associated this is stack should have these tags
  tags = {
    happy_env           = var.deployment_stage,
    happy_stack_name    = var.custom_stack_name,
    happy_service_name  = var.app_name,
    happy_region        = data.aws_region.current.name,
    happy_image         = var.image,
    happy_service_type  = "TODO", # implement this: PRIVATE, PUBLIC, INTERNAL,
    happy_last_applied  = timestamp(),
    #happy-last-git-hash = data.external.git_sha.result.sha
    #happy-repo          = data.external.git_repo.result
  }
}
/*
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
}*/
