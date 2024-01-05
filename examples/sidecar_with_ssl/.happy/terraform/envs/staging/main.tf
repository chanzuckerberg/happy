module "stack" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks?ref=main"

  image_tag        = var.image_tag
  image_tags       = jsondecode(var.image_tags)
  app_name         = var.app
  stack_name       = var.stack_name
  deployment_stage = "staging"
  stack_prefix     = "/${var.stack_name}"
  k8s_namespace    = var.k8s_namespace

  services = {
    frontend = {
      name                  = "frontend"
      cpu                   = "100m"
      memory                = "100Mi"
      port                  = 3000
      service_type          = "INTERNAL"
      platform_architecture = "arm64"
    },
    backend = {
      name                  = "backend"
      cpu                   = "100m"
      memory                = "100Mi"
      port                  = 3000
      service_type          = "INTERNAL"
      platform_architecture = "arm64"
    }
  }
}
