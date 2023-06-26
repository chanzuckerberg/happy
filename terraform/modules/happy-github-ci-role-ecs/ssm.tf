data "aws_caller_identity" "current" {}

resource "random_pet" "this" {
  keepers = {
    role_name = var.gh_actions_role_name
  }
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
      "arn:aws:ssm:us-west-2:${data.aws_caller_identity.current.account_id}:parameter/happy/${var.env}/*",
      # this is the new location of SSM parameters (namespaced on the happy app)
      "arn:aws:ssm:us-west-2:${data.aws_caller_identity.current.account_id}:parameter/happy/${var.happy_app_name}/${var.env}/*"
    ]
  }
}
resource "aws_iam_role_policy" "ssm_reader_writer" {
  name   = "gh_actions_ssm_reader_writer_${random_pet.this.id}"
  policy = data.aws_iam_policy_document.ssm_reader_writer.json
  role   = var.gh_actions_role_name
}
