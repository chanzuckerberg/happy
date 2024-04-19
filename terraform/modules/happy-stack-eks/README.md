<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.14 |
| <a name="requirement_happy"></a> [happy](#requirement\_happy) | >= 0.108.0 |
| <a name="requirement_kubernetes"></a> [kubernetes](#requirement\_kubernetes) | >= 2.16 |
| <a name="requirement_random"></a> [random](#requirement\_random) | >= 3.4.3 |
| <a name="requirement_validation"></a> [validation](#requirement\_validation) | 1.0.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_happy"></a> [happy](#provider\_happy) | >= 0.108.0 |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | >= 2.16 |
| <a name="provider_random"></a> [random](#provider\_random) | >= 3.4.3 |
| <a name="provider_validation"></a> [validation](#provider\_validation) | 1.0.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_datadog_dashboard"></a> [datadog\_dashboard](#module\_datadog\_dashboard) | ../happy-datadog-dashboard | n/a |
| <a name="module_datadog_synthetic"></a> [datadog\_synthetic](#module\_datadog\_synthetic) | ../happy-datadog-synthetics | n/a |
| <a name="module_services"></a> [services](#module\_services) | ../happy-service-eks | n/a |
| <a name="module_tasks"></a> [tasks](#module\_tasks) | ../happy-task-eks | n/a |

## Resources

| Name | Type |
|------|------|
| [kubernetes_secret.oidc_config](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/secret) | resource |
| [random_pet.suffix](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [validation_error.mix_of_internal_and_external_services](https://registry.terraform.io/providers/tlkamp/validation/1.0.0/docs/resources/error) | resource |
| [happy_resolved_app_configs.configs](https://registry.terraform.io/providers/chanzuckerberg/happy/latest/docs/data-sources/resolved_app_configs) | data source |
| [kubernetes_secret.integration_secret](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/data-sources/secret) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_env_vars"></a> [additional\_env\_vars](#input\_additional\_env\_vars) | Additional environment variables to add to the container | `map(string)` | `{}` | no |
| <a name="input_additional_env_vars_from_config_maps"></a> [additional\_env\_vars\_from\_config\_maps](#input\_additional\_env\_vars\_from\_config\_maps) | Additional environment variables to add to the container from the following config maps | <pre>object({<br>    items : optional(list(string), []),<br>    prefix : optional(string, ""),<br>  })</pre> | <pre>{<br>  "items": [],<br>  "prefix": ""<br>}</pre> | no |
| <a name="input_additional_env_vars_from_secrets"></a> [additional\_env\_vars\_from\_secrets](#input\_additional\_env\_vars\_from\_secrets) | Additional environment variables to add to the container from the following secrets | <pre>object({<br>    items : optional(list(string), []),<br>    prefix : optional(string, ""),<br>  })</pre> | <pre>{<br>  "items": [],<br>  "prefix": ""<br>}</pre> | no |
| <a name="input_additional_hostnames"></a> [additional\_hostnames](#input\_additional\_hostnames) | The set of hostnames that will be allowed by the corresponding load balancers and ingress'. These hosts can be configured outside of happy, for instance through a CloudFront distribution. | `set(string)` | `[]` | no |
| <a name="input_additional_pod_labels"></a> [additional\_pod\_labels](#input\_additional\_pod\_labels) | Additional labels to add to the pods. | `map(string)` | `{}` | no |
| <a name="input_additional_volumes_from_config_maps"></a> [additional\_volumes\_from\_config\_maps](#input\_additional\_volumes\_from\_config\_maps) | Additional volumes to add to the container from the following config maps | <pre>object({<br>    items : optional(list(string), []),<br>  })</pre> | <pre>{<br>  "items": []<br>}</pre> | no |
| <a name="input_additional_volumes_from_secrets"></a> [additional\_volumes\_from\_secrets](#input\_additional\_volumes\_from\_secrets) | Additional volumes to add to the container from the following secrets | <pre>object({<br>    items : optional(list(string), []),<br>    base_dir : optional(string, "/var"),<br>  })</pre> | <pre>{<br>  "base_dir": "/var",<br>  "items": []<br>}</pre> | no |
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | The happy application name | `string` | `""` | no |
| <a name="input_create_dashboard"></a> [create\_dashboard](#input\_create\_dashboard) | Create a dashboard for this stack | `bool` | `false` | no |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | Deployment stage for the app | `string` | n/a | yes |
| <a name="input_enable_service_mesh"></a> [enable\_service\_mesh](#input\_enable\_service\_mesh) | Enable service mesh for this stack | `bool` | `false` | no |
| <a name="input_image_tag"></a> [image\_tag](#input\_image\_tag) | Please provide a default image tag | `string` | n/a | yes |
| <a name="input_image_tags"></a> [image\_tags](#input\_image\_tags) | Override image tag for each docker image | `map(string)` | `{}` | no |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this stack | `string` | n/a | yes |
| <a name="input_routing_method"></a> [routing\_method](#input\_routing\_method) | Traffic routing method for this stack. Valid options are 'DOMAIN', when every service gets a unique domain name, or a 'CONTEXT' when all services share the same domain name, and routing is done by request path. | `string` | `"DOMAIN"` | no |
| <a name="input_services"></a> [services](#input\_services) | The services you want to deploy as part of this stack. | <pre>map(object({<br>    name         = string,<br>    service_type = optional(string, "INTERNAL"),<br>    allow_mesh_services = optional(list(object({<br>      service              = optional(string, null),<br>      stack                = optional(string, null),<br>      service_account_name = optional(string, null)<br>    })), null),<br>    ingress_security_groups = optional(list(string), []), // Only used for VPC service_type<br>    alb = optional(object({<br>      name          = string,<br>      listener_port = number,<br>    }), null), // Only used for TARGET_GROUP_ONLY<br>    desired_count                    = optional(number, 2),<br>    max_count                        = optional(number, 5),<br>    max_unavailable_count            = optional(string, "1"),<br>    scaling_cpu_threshold_percentage = optional(number, 80),<br>    port                             = optional(number, 80),<br>    scheme                           = optional(string, "HTTP"),<br>    cmd                              = optional(list(string), []),<br>    args                             = optional(list(string), []),<br>    image_pull_policy                = optional(string, "IfNotPresent"), // Supported values= IfNotPresent, Always, Never<br>    tag_mutability                   = optional(bool, true),<br>    scan_on_push                     = optional(bool, false),<br>    service_port                     = optional(number, null),<br>    service_scheme                   = optional(string, "HTTP"),<br>    linkerd_additional_skip_ports    = optional(set(number), []),<br>    memory                           = optional(string, "500Mi"),<br>    memory_requests                  = optional(string, "200Mi"),<br>    cpu                              = optional(string, "1"),<br>    cpu_requests                     = optional(string, "500m"),<br>    gpu                              = optional(number, null), // Whole number of GPUs to request, 0 will schedule all available GPUs. Requires GPU-enabled nodes in the cluster, `k8s-device-plugin` installed, platform_architecture = "amd64", and additional_node_selectors = { "nvidia.com/gpu.present" = "true" } present.<br>    health_check_path                = optional(string, "/"),<br>    health_check_command             = optional(list(string), [])<br>    aws_iam = optional(object({<br>      policy_json          = optional(string, ""),<br>      service_account_name = optional(string, null),<br>    }), {}),<br>    path                      = optional(string, "/*"), // Only used for CONTEXT and TARGET_GROUP_ONLY routing<br>    priority                  = optional(number, 0),    // Only used for CONTEXT and TARGET_GROUP_ONLY routing<br>    success_codes             = optional(string, "200-499"),<br>    synthetics                = optional(bool, false),<br>    initial_delay_seconds     = optional(number, 30),<br>    alb_idle_timeout          = optional(number, 60) // in seconds<br>    period_seconds            = optional(number, 3),<br>    liveness_timeout_seconds  = optional(number, 30),<br>    readiness_timeout_seconds = optional(number, 30),<br>    progress_deadline_seconds = optional(number, 600),<br>    platform_architecture     = optional(string, "amd64"), // Supported values= amd64, arm64; GPU nodes are amd64 only.<br>    additional_node_selectors = optional(map(string), {}), // For GPU use= { "nvidia.com/gpu.present" = "true" }<br>    bypasses = optional(map(object({                       // Only used for INTERNAL service_type<br>      paths   = optional(set(string), [])<br>      methods = optional(set(string), [])<br>    })), {})<br>    sticky_sessions = optional(object({<br>      enabled          = optional(bool, false),<br>      duration_seconds = optional(number, 600),<br>      cookie_name      = optional(string, "happy_sticky_session"),<br>    }), {})<br>    sidecars = optional(map(object({<br>      image                     = string<br>      tag                       = string<br>      cmd                       = optional(list(string), [])<br>      args                      = optional(list(string), [])<br>      port                      = optional(number, 80)<br>      scheme                    = optional(string, "HTTP")<br>      memory                    = optional(string, "200Mi")<br>      cpu                       = optional(string, "500m")<br>      image_pull_policy         = optional(string, "IfNotPresent") // Supported values= IfNotPresent, Always, Never<br>      health_check_path         = optional(string, "/")<br>      initial_delay_seconds     = optional(number, 30)<br>      period_seconds            = optional(number, 3)<br>      liveness_timeout_seconds  = optional(number, 30)<br>      readiness_timeout_seconds = optional(number, 30)<br>    })), {})<br>    init_containers = optional(map(object({<br>      image = string<br>      tag   = string<br>      cmd   = optional(list(string), []),<br>    })), {}),<br>    additional_env_vars    = optional(map(string), {}),<br>    cache_volume_mount_dir = optional(string, "/var/shared/cache"),<br>    oidc_config = optional(object({<br>      issuer                = string<br>      authorizationEndpoint = string<br>      tokenEndpoint         = string<br>      userInfoEndpoint      = string<br>      secretName            = string<br>    }), null)<br>  }))</pre> | n/a | yes |
| <a name="input_skip_config_injection"></a> [skip\_config\_injection](#input\_skip\_config\_injection) | Skip injecting app configs into the services / tasks | `bool` | `false` | no |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Happy Path stack name | `string` | n/a | yes |
| <a name="input_stack_prefix"></a> [stack\_prefix](#input\_stack\_prefix) | Do bucket storage paths and db schemas need to be prefixed with the stack name? (Usually '/{stack\_name}' for dev stacks, and '' for staging/prod stacks) | `string` | `""` | no |
| <a name="input_tasks"></a> [tasks](#input\_tasks) | The deletion/migration tasks you want to run when a stack comes up and down. | <pre>map(object({<br>    image : string,<br>    memory : optional(string, "200Mi"),<br>    cpu : optional(string, "500m"),<br>    cmd : optional(list(string), []),<br>    args : optional(list(string), []),<br>    platform_architecture : optional(string, "amd64"), // Supported values: amd64, arm64<br>    is_cron_job : optional(bool, false),<br>    aws_iam : optional(object({<br>      policy_json : optional(string, ""),<br>      service_account_name : optional(string, null),<br>    }), {}),<br>    cron_schedule : optional(string, "0 0 1 1 *"),<br>    additional_env_vars : optional(map(string), {}),<br>  }))</pre> | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_dashboard"></a> [dashboard](#output\_dashboard) | n/a |
| <a name="output_service_ecrs"></a> [service\_ecrs](#output\_service\_ecrs) | n/a |
| <a name="output_service_endpoints"></a> [service\_endpoints](#output\_service\_endpoints) | The URL endpoints for services |
| <a name="output_task_arns"></a> [task\_arns](#output\_task\_arns) | ARNs for all the tasks |
<!-- END -->

