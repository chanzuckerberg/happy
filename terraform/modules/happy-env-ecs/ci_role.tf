data "aws_caller_identity" "current" {}
module "happy_github_ci_role" {
  for_each = toset([for role in var.github_actions_roles : role.name])
  source   = "../happy-github-ci-role"

  ecrs = module.ecr
  ecs = {
    arn = module.ecs-cluster.arn
    happy_app_name = module.ecs-cluster.cluster_name
  }
  gh_actions_role_name = each.value
  dynamodb_table_arn   = aws_dynamodb_table.locks.arn
  tags                 = var.tags
}

module "integration_secret_reader_policy" {
  for_each = var.authorized_github_repos

  source       = "git@github.com:chanzuckerberg/cztack//aws-iam-secrets-reader-policy?ref=v0.43.3"
  role_name    = module.happy_github_ci_role[each.key].role_name
  secrets_arns = [aws_secretsmanager_secret.happy_env_secret.arn]
  depends_on   = [module.happy_github_ci_role]
}

data "aws_iam_policy_document" "ecs_reader" {
  statement {
    sid = "GHActionsECSReader"
    actions = [
      "ecs:List*",
      "ecs:Describe*",
    ]
    resources = [
      "arn:aws:ecs:us-west-2:${data.aws_caller_identity.current.account_id}:cluster/${module.ecs-cluster.cluster_name}/*",
      "arn:aws:ecs:us-west-2:${data.aws_caller_identity.current.account_id}:service/${module.ecs-cluster.cluster_name}/*",
    ]
  }
}
resource "aws_iam_role_policy" "ecs_reader" {
  for_each = var.authorized_github_repos

  name       = "gh_actions_ecs_reader_${replace("${var.tags.project}_${var.tags.env}_${var.tags.service}_${each.value.app_name}", "-", "_")}"
  policy     = data.aws_iam_policy_document.ecs_reader.json
  role       = module.happy_github_ci_role[each.key].role_name
  depends_on = [module.happy_github_ci_role]
}

resource "aws_iam_role_policy_attachment" "locktable_policy_attachment" {
  for_each = var.authorized_github_repos

  role       = module.happy_github_ci_role[each.key].role_name
  policy_arn = aws_iam_policy.locktable_policy.arn
}
