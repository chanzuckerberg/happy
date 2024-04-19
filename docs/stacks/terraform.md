---
parent: Stacks
layout: default
has_toc: true
---

# Terraform
{: .no_toc }

<details open markdown="block">
  <summary>
    Table of contents
  </summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

## Overview

Happy uses terraform to do all deployments. Whenever you do a `happy create`, `happy update`, or `happy delete` happy is applying terraform plans within the stack's terraform workspace.
Each stack's terraform is executed in its own terraform workspace and is not shared with other stacks. The terraform that is applied is specified in the [happy config.json file](../config/config_json.md).
Conventially, it is usually located in `.happy/terraform/envs/<env>/*.tf`. The module we use to build stacks is [`git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks`](https://github.com/chanzuckerberg/happy/tree/main/terraform/modules/happy-stack-eks).

## Stack Configuration

Here is an example configuration of a happy stack terraform module:

~~~terraform
module "stack" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks?ref=main"

  image_tag        = var.image_tag
  image_tags       = jsondecode(var.image_tags)
  app_name         = var.app
  stack_name       = var.stack_name
  deployment_stage = var.env

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
      platform_architecture            = "arm64"
      // for internal services protected by OIDC, you might want to bypass authentication
      // for certain HTTP methods or certain paths (such as health checks)
      // these bypasses bypass authentication so use them sparingly
      bypasses = {
        mybypass = {
          paths   = ["/api/health"]
          methods = ["GET"]
        }
        mybypass2 = {
          paths   = ["/api/*"]
          methods = ["PATCH"]
        }
      }
    }
  }

  // tasks can be utilized to run post-deployment tasks such as database migrations or deletions
  tasks = {
  }
}
~~~

This module configuration specifies which services you'd like to deploy using happy. Services are set in the services variables
and should match the same service names that show up in your docker-compose. [The rest of the options](https://github.com/chanzuckerberg/happy/blob/main/terraform/modules/happy-stack-eks/variables.tf#L45)
are just additional configuration that you might provide to your deployed container. All of these options have sane defaults so
they can be easily omitted if you aren't sure what to provide. A couple of important ones to always double check:

* `name` - the name of the service, and it should match the service name in docker-compose
* `port` - the port your application is listening on; if this is wrong your deployment might fail to pass its healthcheck
* `health_check_path` - a valid path that your application will respond with an HTTP 200
* `platform_architecture` - (either amd64 or arm64) make sure that if you containerized application matches this architecture by setting a "platform" attribute in your docker-compose.yml; a mismatch will cause the notorious `exec format error`
* `service_type` - this represents the category of service you are deploying. For example, EXTERNAL means an application exposed to the public internet. Valid values are INTERNAL,EXTERNAL,PRIVATE,IMAGE_TEMPLATE,TARGET_GROUP_ONLY
