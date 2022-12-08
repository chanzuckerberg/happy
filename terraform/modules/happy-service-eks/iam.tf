module "iam_service_account" {
  count = var.aws_iam_policy_json == "" ? 0 : 1
  source   = "../happy-iam-service-account-eks"

  eks_cluster   = var.eks_cluster
  k8s_namespace = var.k8s_namespace
  tags = local.tags
}

resource "aws_iam_policy" "policy" {
    count = var.aws_iam_policy_json == "" ? 0 : 1

  name        = module.iam_service_account[0].iam_role
  path        = "/"
  description = "Stack policy for ${module.iam_service_account[0].iam_role}"
  policy      = var.aws_iam_policy_json
  tags        = local.tags
}

resource "aws_iam_policy_attachment" "attach" {
    count = var.aws_iam_policy_json == "" ? 0 : 1

  name       = module.iam_service_account[0].iam_role
  roles      = [module.iam_service_account[0].iam_role]
  policy_arn = aws_iam_policy.policy.arn
}
