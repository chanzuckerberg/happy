module "iam_service_account" {
  source = "../happy-iam-service-account-eks"

  eks_cluster         = var.eks_cluster
  k8s_namespace       = var.k8s_namespace
  aws_iam_policy_json = var.aws_iam.policy_json
  tags                = local.tags
}
