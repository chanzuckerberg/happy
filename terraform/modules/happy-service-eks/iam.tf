module "iam_service_account" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-iam-service-account-eks?ref=heathj/fix-bugs"

  eks_cluster   = var.eks_cluster
  k8s_namespace = var.k8s_namespace

  tags = var.tags
}

resource "aws_iam_policy" "policy" {
  name        = module.iam_service_account.iam_role
  path        = "/"
  description = "Stack policy for ${module.iam_service_account.iam_role}"
  policy      = var.aws_iam_policy_json

  tags        = local.tags
}

resource "aws_iam_policy_attachment" "attach" {
  name       = module.iam_service_account.iam_role
  roles      = [module.iam_service_account.iam_role]
  policy_arn = aws_iam_policy.policy.arn
}
