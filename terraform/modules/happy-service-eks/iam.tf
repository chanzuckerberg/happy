module "iam_service_account" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-iam-service-account-eks?ref=c57e66efef8190638c3df678c318d9d9c316e4a6"

  eks_cluster   = var.eks_cluster
  k8s_namespace = var.k8s_namespace

  tags = var.tags
}

resource "aws_iam_policy" "policy" {
  name        = module.iam_service_account.iam_role
  path        = "/"
  description = "Stack policy for ${module.iam_service_account.iam_role}"
  policy      = var.aws_iam_policy_json

  tags        = var.tags
}

resource "aws_iam_policy_attachment" "attach-logs" {
  name       = module.iam_service_account.iam_role
  roles      = [module.iam_service_account.iam_role]
  policy_arn = aws_iam_policy.policy.arn

  tags       = var.tags
}
