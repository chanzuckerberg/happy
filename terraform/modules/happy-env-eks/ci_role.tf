module "happy_github_ci_role" {
  for_each = var.authorized_github_repos
  source   = "../happy-github-ci-role"

  ecr_repo_arns           = flatten([for ecr in module.ecrs : ecr.repository_arn])
  authorized_github_repos = [each.value.repo_name]
  happy_app_name          = each.value.app_name
  eks_cluster_arn = var.eks-cluster.cluster_arn
  tags = var.tags
}
