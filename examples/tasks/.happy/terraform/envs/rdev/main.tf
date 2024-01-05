module "stack" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks?ref=main"

  image_tag        = var.image_tag
  image_tags       = jsondecode(var.image_tags)
  app_name         = var.app
  stack_name       = var.stack_name
  deployment_stage = "rdev"

  stack_prefix  = "/${var.stack_name}"
  k8s_namespace = var.k8s_namespace

  // this allow these services under the same domain
  // each service is reachable via their path configured below
  routing_method = "CONTEXT"

  services = {
    frontend = {
      name          = "frontend"
      desired_count = 1
      // the maximum number of copies of this service it can autoscale to
      max_count = 50
      // the signal used to trigger autoscaling (ie. 50% of CPU means scale up)
      scaling_cpu_threshold_percentage = 50
      // the port the service is running on
      port   = 3000
      memory = "100Mi"
      cpu    = "100m"
      // an endpoint that returns a 200. Your service will not start if this endpoint is not healthy
      health_check_path = "/health"
      // oneof: INTERNAL, EXTERNAL, PRIVATE, TARGET_GROUP_ONLY, IMAGE_TEMPLATE
      // INTERNAL: OIDC protected endpoints
      // EXTERNAL: internet accessible
      // PRIVATE: only accessible within the cluster
      // TARGET_GROUP_ONLY: attach to an existing ALB rather than making a new one
      // IMAGE_TEMPLATE: don't deploy any services, just use to create and push images
      service_type = "INTERNAL"
      // the path to reach this search
      path = "/*"
      // the platform architecture of the container. this should match what is in
      // the platform attribute of your docker-compose.yml file for your service.
      // oneof: amd64, arm64.
      // Try to always select arm since it comes with a lot of cost savings and performance
      // benefits and has little to no impact on developers.
      platform_architecture = "amd64"
    }
  }

  // tasks can be utilized to run post-deployment tasks such as database migrations or deletions
  tasks = {
    delete = {
      image  = "{frontend}:${var.image_tag}"
      memory = "100Mi"
      cpu    = "100m"
      cmd    = ["/app/delete.sh"]
    }

    migrate = {
      image  = "{frontend}:${var.image_tag}"
      memory = "100Mi"
      cpu    = "100m"
      cmd    = ["/app/migrate.sh"]
    }
  }
}
