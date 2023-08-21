output "frontend_url" {
  value       = nonsensitive(module.stack.url)
  description = "The URL endpoint for the service"
}
