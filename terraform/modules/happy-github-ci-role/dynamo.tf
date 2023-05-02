module "dynamodb_writer" {
  count     = var.dynamodb_table_arn != "" ? 1 : 0
  source    = "git@github.com:chanzuckerberg/cztack//aws-iam-policy-dynamodb-rw?ref=v0.55.0"
  table_arn = var.dynamodb_table_arn
  role_name = var.gh_actions_role_name
  tags      = var.tags
}