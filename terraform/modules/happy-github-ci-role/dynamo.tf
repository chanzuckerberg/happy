module "dynamodb_writer" {
  count     = var.dynamodb_table_arn != "" ? 1 : 0
  source    = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-iam-policy-dynamodb-rw?ref=v0.290.3"
  table_arn = var.dynamodb_table_arn
  role_name = var.gh_actions_role_name
}