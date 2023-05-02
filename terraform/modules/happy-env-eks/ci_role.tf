module "happy_github_ci_role" {
  for_each = var.github_actions_roles
  source   = "../happy-github-ci-role"

  ecrs                 = module.ecrs
  gh_actions_role_name = each.value.name
  eks_cluster_arn      = var.eks-cluster.cluster_arn
  tags                 = var.tags
}
