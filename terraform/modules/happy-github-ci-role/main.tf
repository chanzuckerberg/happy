data "aws_caller_identity" "current" {}
locals {
  namespace = replace("${var.tags.project}_${var.tags.env}_${var.tags.service}_${var.happy_app_name}", "-", "_")
  role_name = "gh_actions_${local.namespace}"
}

module "gh_actions_role" {
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-iam-role-github-action?ref=v0.125.0"

  role = {
    name = local.role_name
  }
  authorized_github_repos = {
    chanzuckerberg : var.authorized_github_repos
  }

  tags = var.tags
}

module "ecr_writer_policy" {
  count               = length(var.ecr_repo_arns) > 0 ? 1 : 0
  source              = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-iam-policy-ecr-writer?ref=v0.125.0"
  role_name           = local.role_name
  ecr_repository_arns = var.ecr_repo_arns
  policy_name         = "gh_actions_ecr_push_${local.namespace}"
  depends_on          = [module.gh_actions_role]

  project = var.tags.project
  env     = var.tags.env
  service = var.tags.service
  owner   = var.tags.owner
}


data "aws_iam_policy_document" "ssm_reader_writer" {
  statement {
    sid = "GhActionsSSMReaderWriter"
    actions = [
      "ssm:Get*",
      "ssm:Put*"
    ]
    resources = [
      # this is the legacy location of SSM parameters
      "arn:aws:ssm:us-west-2:${data.aws_caller_identity.current.account_id}:parameter/happy/${var.tags.env}/*",
      # this is the new location of SSM parameters (namespaced on the happy app)
      "arn:aws:ssm:us-west-2:${data.aws_caller_identity.current.account_id}:parameter/happy/${var.happy_app_name}/${var.tags.env}/*"
    ]
  }
}
resource "aws_iam_role_policy" "ssm_reader_writer" {
  name       = "gh_actions_ssm_reader_writer_${local.namespace}"
  policy     = data.aws_iam_policy_document.ssm_reader_writer.json
  role       = local.role_name
  depends_on = [module.gh_actions_role]

}
