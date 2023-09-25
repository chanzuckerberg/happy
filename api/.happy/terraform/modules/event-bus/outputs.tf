output "k8s_secrets_name" {
  sensitive = false
  value     = kubernetes_secret_v1.event_bus_secrets.metadata[0].name
}
