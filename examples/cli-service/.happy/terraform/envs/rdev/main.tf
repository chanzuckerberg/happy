# Auto-generated by 'happy infra'. Do not edit
# Make improvements in happy, so that everyone can benefit.
module "stack" {
  source           = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks?ref=main"
  image_tag        = var.image_tag
  app_name         = var.app
  stack_name       = var.stack_name
  k8s_namespace    = var.k8s_namespace
  image_tags       = jsondecode(var.image_tags)
  stack_prefix     = "/${var.stack_name}"
  deployment_stage = "rdev"
  services = {
    offlinejob = {
      cpu                              = "200m"
      cpu_requests                     = "100m"
      desired_count                    = 1
      initial_delay_seconds            = 30
      max_count                        = 1
      memory                           = "256Mi"
      memory_requests                  = "128Mi"
      name                             = "offlinejob"
      path                             = "/*"
      period_seconds                   = 3
      platform_architecture            = "amd64"
      priority                         = 0
      scaling_cpu_threshold_percentage = 80
      service_type                     = "CLI"
      synthetics                       = false
      health_check_command              = ["/bin/true"]
    }
  }
  create_dashboard = false
  routing_method   = "CONTEXT"
}
