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

output "panther_waf_configuration" {
  value = var.include_waf ? {
    panther_role = module.regional-waf[0].panther-role
    log_bucket = module.regional-waf[0].web_acl_log_bucket
  } : {}
  sensitive = false
  description = "Outputs that help Security Eng team configure Panther monitoring"
}
