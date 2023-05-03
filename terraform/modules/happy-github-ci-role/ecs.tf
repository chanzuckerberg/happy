module "ecs_cluster_permissions" {
  count  = length(var.ecs.arn) == 0 ? 0 : 1
  source = "../happy-github-ci-role-ecs"
  #  ecs_cluster_arn      = var.ecs.cluster_arn TODO: not needed yet
  env                  = var.tags.env
  happy_app_name       = var.ecs.happy_app_name
  gh_actions_role_name = var.gh_actions_role_name
}