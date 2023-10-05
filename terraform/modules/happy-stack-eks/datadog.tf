locals {
  cluster_id = local.secret["eks_cluster"].cluster_id
}
module "datadog_dashboard" {
  count     = var.create_dashboard ? 1 : 0
  source = "../happy-datadog-dashboard"
  cluster_id = local.cluster_id
  k8s_namespace = var.k8s_namespace
  stack_name = var.stack_name
}
