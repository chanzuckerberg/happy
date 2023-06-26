module "eks_cluster_permissions" {
  count                = var.eks_cluster_arn != "" ? 1 : 0
  source               = "../happy-github-ci-role-eks"
  eks_cluster_arn      = var.eks_cluster_arn
  gh_actions_role_name = var.gh_actions_role_name
}