resource "aws_iam_user" "tfe-happy" {
  name = "tfe-happy-${var.happy_app_name}"
  path = "/tfe/"

  tags = { "owner" : "infra-eng@chanzuckerberg.com" }
}

module "tfe-si-happy-roles" {
  source = "github.com/chanzuckerberg/cztack//aws-iam-group-assume-role?ref=v0.43.1"

  group_name = "tfe-si-happy-${var.happy_app_name}"
  iam_path   = "/tfe/"

  users           = [aws_iam_user.tfe-happy.name]
  target_accounts = var.aws_accounts_can_assume

  target_role = "tfe-si"
}
