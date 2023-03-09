output "namespace" {
  value = local.k8s_namespace
}


output "dashboard" {
  value = {
    id  = datadog_dashboard_json.environment_dashboard.id
    url = datadog_dashboard_json.environment_dashboard.url
  }
}

output "integration_secret" {
  value     = local.secret_string
  sensitive = true
}
