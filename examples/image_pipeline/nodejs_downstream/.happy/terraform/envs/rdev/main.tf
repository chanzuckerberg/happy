module "stack" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks?ref=main"

  image_tag        = var.image_tag
  image_tags       = jsondecode(var.image_tags)
  stack_name       = var.stack_name
  deployment_stage = "rdev"
  stack_prefix     = "/${var.stack_name}"
  k8s_namespace    = var.k8s_namespace

  services = {
    frontend = {
      name         = "frontend"
      cpu          = "100m"
      memory       = "100Mi"
      port         = "80"
      service_type = "INTERNAL"
    },
    backend = {
      name         = "backend"
      cpu          = "100m"
      memory       = "100Mi"
      port         = "80"
      service_type = "INTERNAL"
    }
  }
  tasks = {
  }
}
