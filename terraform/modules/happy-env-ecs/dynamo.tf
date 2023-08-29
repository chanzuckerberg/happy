resource "aws_dynamodb_table" "locks" {
  name         = "${var.tags.project}-${var.tags.env}-${var.tags.service}-stacklist"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "key"

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
  name   = "${var.tags.project}-${var.tags.env}-${var.tags.service}-stacklist"
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
