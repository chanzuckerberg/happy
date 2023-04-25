locals {
  namespace = replace("${var.tags.project}_${var.tags.env}_${var.tags.service}_${var.happy_app_name}", "-", "_")
  role_name = "gh_actions_${local.namespace}"
}

module "gh_actions_role" {
  source = "git@github.com:chanzuckerberg/cztack//aws-iam-role-github-action?ref=v0.54.0"

  role = {
    name = local.role_name
  }
  authorized_github_repos = {
    chanzuckerberg : var.authorized_github_repos
  }

  tags = var.tags
}
