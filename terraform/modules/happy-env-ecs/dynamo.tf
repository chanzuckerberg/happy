# create a dynamodb resource here
resource "aws_dynamodb_table" "locks" {
  name           = "${local.project}-${local.env}-${local.service}-locks"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "key"

  attribute {
    name = "key"
    type = "S"
  }

  server_side_encryption {
    enabled = true
  }

  tags = var.tags
}

resource "aws_iam_policy" "locktable_policy" {
  name   = "${local.project}-${local.env}-${local.service}-locks-access"
  path   = "/"
  policy = data.aws_iam_policy_document.locktable_policy_document.json
}

data "aws_iam_policy_document" "locktable_policy_document" {
  statement {
    sid    = "DynamoDBLocktableAccess"
    effect = "Allow"

    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem"
    ]
    resources = [aws_dynamodb_table.locks.arn]
  }
}
