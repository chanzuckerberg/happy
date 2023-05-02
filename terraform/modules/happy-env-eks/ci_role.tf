module "happy_github_ci_role" {
  for_each = toset([for role in var.github_actions_roles: role.name])
  source   = "../happy-github-ci-role"

  ecrs                 = module.ecrs
  gh_actions_role_name = each.value
  eks_cluster_arn      = var.eks-cluster.cluster_arn
  tags                 = var.tags
}
