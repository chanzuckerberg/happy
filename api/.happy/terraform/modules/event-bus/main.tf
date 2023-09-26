data "kubernetes_namespace_v1" "namespace" {
  metadata {
    name = var.k8s_namespace
  }
}

data "kubernetes_secret_v1" "integration_secret" {
  metadata {
    name      = "integration-secret"
    namespace = data.kubernetes_namespace_v1.namespace.metadata[0].name
  }
}

data "aws_caller_identity" "current" {}

locals {
  secret = jsondecode(nonsensitive(data.kubernetes_secret_v1.integration_secret.data.integration_secret))
  tags   = local.secret["tags"]
}

resource "kubernetes_secret_v1" "event_bus_secrets" {
  metadata {
    name      = "event-consumer-${var.stack_name}"
    namespace = data.kubernetes_namespace_v1.namespace.metadata[0].name
  }

  data = {
    "EVENT_CONSUMER_QUEUE_URL" = aws_sqs_queue.events_queue.url
    "EVENTS_TOPIC_ARN" = aws_sns_topic.events_topic.arn
  }
}

resource "aws_sns_topic" "events_topic" {
  name = "hapi-events-${local.tags.env}-${var.stack_name}"
  tags = local.tags
}

resource "aws_sqs_queue" "events_queue" {
  name                      = "hapi-events-${local.tags.env}-${var.stack_name}"
  receive_wait_time_seconds = 20

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.events_queue_deadletter.arn
    maxReceiveCount = 5
  })

  tags = local.tags
}

resource "aws_sqs_queue" "events_queue_deadletter" {
  name                      = "hapi-events-${local.tags.env}-${var.stack_name}-deadletter"
  receive_wait_time_seconds = 20
  tags                      = local.tags
}


data "aws_iam_policy_document" "events_queue_policy_document" {
  statement {
    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
    }

    actions   = ["sqs:*"]
    resources = [aws_sqs_queue.events_queue.arn]
  }

  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["sns.amazonaws.com"]
    }

    actions   = ["sqs:SendMessage"]
    resources = [aws_sqs_queue.events_queue.arn]

    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"
      values   = [aws_sns_topic.events_topic.arn]
    }
  }
}

resource "aws_sqs_queue_policy" "events_queue_policy" {
  queue_url = aws_sqs_queue.events_queue.id
  policy    = data.aws_iam_policy_document.events_queue_policy_document.json
}

resource "aws_sns_topic_subscription" "user_updates_sqs_target" {
  topic_arn = aws_sns_topic.events_topic.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.events_queue.arn
}
