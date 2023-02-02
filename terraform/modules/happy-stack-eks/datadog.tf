locals {
  cluster_id = local.secret["eks_cluster"]
}
resource "datadog_dashboard_json" "stack_dashboard" {
  count     = var.create_dashboard ? 1 : 0
  dashboard = <<EOF
  {
	"title": "[HAPPY] ${local.cluster_id} / ${var.stack_name} stack Dashboard"
}
  EOF
}
