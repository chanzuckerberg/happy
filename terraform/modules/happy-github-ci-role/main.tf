locals {
  namespace = replace("${var.tags.project}_${var.tags.env}_${var.tags.service}", "-", "_")
  role_name = "gh_actions"
}

module "gh_actions_role" {
  source = "git@github.com:chanzuckerberg/cztack//aws-iam-role-github-action?ref=v0.54.0"

  role = {
    name = local.role_name
  }
  authorized_github_repos = {
    chanzuckerberg : [for repo in var.authorized_github_repos : repo.repo_name]
  }

  tags = var.tags
}
