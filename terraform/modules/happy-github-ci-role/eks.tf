module "eks_cluster_permissions" {
  count           = var.eks_cluster_arn != "" ? 1 : 0
  source          = "../happy-github-ci-role-eks"
  eks_cluster_arn = module.eks_cluster.arn
  gh_actions_role = module.gh_actions_role

  depends_on = [
    module.gh_actions_role
  ]
}