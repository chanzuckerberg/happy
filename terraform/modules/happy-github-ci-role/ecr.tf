module "ecr_writer_policy" {
  count               = length(var.ecr_repo_arns) > 0 ? 1 : 0
  source              = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-iam-policy-ecr-writer?ref=v0.125.0"
  role_name           = local.role_name
  ecr_repository_arns = var.ecr_repo_arns
  policy_name         = "gh_actions_ecr_push_${local.namespace}"

  project = var.tags.project
  env     = var.tags.env
  service = var.tags.service
  owner   = var.tags.owner

  depends_on = [module.gh_actions_role]
}

// used for the dynamic autocreated ECRs
module "autocreated_ecr_writer_policy" {
  source    = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-iam-policy-ecr-writer?ref=v0.125.0"
  role_name = local.role_name
  // TODO: not a super fan of this. Would be ideal to have the role only have access to the stacks created by this happy project
  ecr_repository_arns = ["arn:aws:ecr:us-west-2:${data.aws_caller_identity.current.account_id}:repository/*/${var.tags.env}/*"]
  policy_name         = "gh_actions_ecr_push_${local.namespace}"

  project = var.tags.project
  env     = var.tags.env
  service = var.tags.service
  owner   = var.tags.owner

  depends_on = [module.gh_actions_role]
}