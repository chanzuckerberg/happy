module "dynamodb_writer" {
  source    = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-iam-policy-dynamodb-rw?ref=v0.290.3"
  table_arn = var.dynamodb_table_arn
  role_name = module.gh_actions_role.role.name

  depends_on = [
    module.gh_actions_role
  ]
}