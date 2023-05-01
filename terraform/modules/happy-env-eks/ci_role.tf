module "happy_github_ci_role" {
  source = "../happy-github-ci-role"

  ecrs                    = module.ecrs
  authorized_github_repos = var.authorized_github_repos
  eks_cluster_arn         = var.eks-cluster.cluster_arn
  tags                    = var.tags
}
