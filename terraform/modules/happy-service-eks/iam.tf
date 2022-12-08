module "iam_service_account" {
  for_each = var.aws_iam_policy_json == "" ? [] : [1]
  source   = "../happy-iam-service-account-eks"

  eks_cluster   = var.eks_cluster
  k8s_namespace = var.k8s_namespace
  tags = local.tags
}

resource "aws_iam_policy" "policy" {
  for_each    = module.iam_service_account

  name        = each.value.iam_role
  path        = "/"
  description = "Stack policy for ${each.value.iam_role}"
  policy      = var.aws_iam_policy_json
  tags        = local.tags
}

resource "aws_iam_policy_attachment" "attach" {
  for_each   = aws_iam_policy.policy

  name       = each.key.iam_role
  roles      = [each.key.iam_role]
  policy_arn = each.value.arn
}
