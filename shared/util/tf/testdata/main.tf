module "stack" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks?ref=happy-stack-eks-v4.2.1"

  image_tag        = var.image_tag
  image_tags       = jsondecode(var.image_tags)
  stack_name       = var.stack_name
  deployment_stage = "rdev"

  stack_prefix  = "/${var.stack_name}"
  k8s_namespace = var.k8s_namespace

  routing_method = "CONTEXT"

  services = {
    frontend = {
      name                             = "frontend"
      desired_count                    = 1
      max_count                        = 50
      scaling_cpu_threshold_percentage = 50
      port                             = 3000
      memory                           = "100Mi"
      cpu                              = "100m"
      health_check_path                = "/health"
      service_type                     = "INTERNAL"
      path                             = "/*"
      platform_architecture            = "amd64" // Has to match that in the docker-compose.yml
    },
    internal-api = {
      name                             = "internal-api"
      desired_count                    = 1
      max_count                        = 50
      scaling_cpu_threshold_percentage = 80
      port                             = 3000
      memory                           = "100Mi"
      cpu                              = "100m"
      health_check_path                = "/api/health"
      service_type                     = "INTERNAL"
      path                             = "/api/*"
      priority                         = 1
      platform_architecture            = "arm64" // Has to match that in the docker-compose.yml
      bypasses = {
        mybypass = {
          paths   = ["/api/health"]
          methods = ["GET"]
        }
        mybypass2 = {
          paths   = ["/*"]
          methods = ["PATCH"]
        }
      }
    }
  }
  tasks = {
  }
}
