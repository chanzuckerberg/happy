output "namespace" {
  value = local.k8s_namespace
}


output "dashboard" {
  value = {
    id  = ""
    url = ""
  }
}

output "integration_secret" {
  value     = local.secret_string
  sensitive = true
}

output "databases" {
  value = { for k, v in module.dbs : k => {
    database_host     = v.endpoint
    database_name     = v.database_name
    database_username = v.master_username
    database_password = v.master_password
  } }
  sensitive = true
}
