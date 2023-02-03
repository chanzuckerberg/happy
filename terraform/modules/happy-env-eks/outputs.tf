output "namespace" {
  value = kubernetes_namespace.happy.id
}


output "dashboard" {
  value = {
    id = datadog_dashboard_json.environment_dashboard.id
    url = datadog_dashboard_json.environment_dashboard.url
  }
}