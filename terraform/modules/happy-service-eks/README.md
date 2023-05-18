<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.45 |
| <a name="requirement_kubernetes"></a> [kubernetes](#requirement\_kubernetes) | >= 2.16 |
| <a name="requirement_random"></a> [random](#requirement\_random) | >= 3.4.3 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 4.45 |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | >= 2.16 |
| <a name="provider_random"></a> [random](#provider\_random) | >= 3.4.3 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_ecr"></a> [ecr](#module\_ecr) | git@github.com:chanzuckerberg/cztack//aws-ecr-repo | v0.53.1 |
| <a name="module_iam_service_account"></a> [iam\_service\_account](#module\_iam\_service\_account) | ../happy-iam-service-account-eks | n/a |
| <a name="module_ingress"></a> [ingress](#module\_ingress) | ../happy-ingress-eks | n/a |

## Resources

| Name | Type |
|------|------|
| [aws_lb_listener_rule.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener_rule) | resource |
| [aws_lb_target_group.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group) | resource |
| [kubernetes_deployment_v1.deployment](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/deployment_v1) | resource |
| [kubernetes_horizontal_pod_autoscaler_v1.hpa](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/horizontal_pod_autoscaler_v1) | resource |
| [kubernetes_manifest.this](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/manifest) | resource |
| [kubernetes_service_v1.service](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service_v1) | resource |
| [random_pet.this](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [aws_lb.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/lb) | data source |
| [aws_lb_listener.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/lb_listener) | data source |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_env_vars"></a> [additional\_env\_vars](#input\_additional\_env\_vars) | Additional environment variables to add to the task definition | `map(string)` | `{}` | no |
| <a name="input_additional_env_vars_from_config_maps"></a> [additional\_env\_vars\_from\_config\_maps](#input\_additional\_env\_vars\_from\_config\_maps) | Additional environment variables to add to the container from the following config maps | <pre>object({<br>    items : optional(list(string), []),<br>    prefix : optional(string, ""),<br>  })</pre> | <pre>{<br>  "items": [],<br>  "prefix": ""<br>}</pre> | no |
| <a name="input_additional_env_vars_from_secrets"></a> [additional\_env\_vars\_from\_secrets](#input\_additional\_env\_vars\_from\_secrets) | Additional environment variables to add to the container from the following secrets | <pre>object({<br>    items : optional(list(string), []),<br>    prefix : optional(string, ""),<br>  })</pre> | <pre>{<br>  "items": [],<br>  "prefix": ""<br>}</pre> | no |
| <a name="input_additional_volumes_from_config_maps"></a> [additional\_volumes\_from\_config\_maps](#input\_additional\_volumes\_from\_config\_maps) | Additional volumes to add to the container from the following config maps | <pre>object({<br>    items : optional(list(string), []),<br>  })</pre> | <pre>{<br>  "items": []<br>}</pre> | no |
| <a name="input_additional_volumes_from_secrets"></a> [additional\_volumes\_from\_secrets](#input\_additional\_volumes\_from\_secrets) | Additional volumes to add to the container from the following secrets | <pre>object({<br>    items : optional(list(string), []),<br>  })</pre> | <pre>{<br>  "items": []<br>}</pre> | no |
| <a name="input_aws_iam_policy_json"></a> [aws\_iam\_policy\_json](#input\_aws\_iam\_policy\_json) | The AWS IAM policy to give to the pod. | `string` | `""` | no |
| <a name="input_certificate_arn"></a> [certificate\_arn](#input\_certificate\_arn) | ACM certificate ARN to attach to the load balancer listener | `string` | n/a | yes |
| <a name="input_cloud_env"></a> [cloud\_env](#input\_cloud\_env) | Typically data.terraform\_remote\_state.cloud-env.outputs | <pre>object({<br>    public_subnets : list(string),<br>    private_subnets : list(string),<br>    database_subnets : list(string),<br>    database_subnet_group : string,<br>    vpc_id : string,<br>    vpc_cidr_block : string,<br>  })</pre> | n/a | yes |
| <a name="input_container_name"></a> [container\_name](#input\_container\_name) | The name of the container | `string` | n/a | yes |
| <a name="input_cpu"></a> [cpu](#input\_cpu) | CPU shares (1cpu=1000m) per pod | `string` | `"100m"` | no |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | The name of the deployment stage of the Application | `string` | `"dev"` | no |
| <a name="input_desired_count"></a> [desired\_count](#input\_desired\_count) | How many instances of this task should we run across our cluster? | `number` | `2` | no |
| <a name="input_eks_cluster"></a> [eks\_cluster](#input\_eks\_cluster) | eks-cluster module output | <pre>object({<br>    cluster_id : string,<br>    cluster_arn : string,<br>    cluster_endpoint : string,<br>    cluster_ca : string,<br>    cluster_oidc_issuer_url : string,<br>    cluster_version : string,<br>    worker_iam_role_name : string,<br>    worker_security_group : string,<br>    oidc_provider_arn : string,<br>  })</pre> | n/a | yes |
| <a name="input_health_check_path"></a> [health\_check\_path](#input\_health\_check\_path) | path to use for health checks | `string` | `"/"` | no |
| <a name="input_image_tag"></a> [image\_tag](#input\_image\_tag) | The image tag to deploy | `string` | n/a | yes |
| <a name="input_initial_delay_seconds"></a> [initial\_delay\_seconds](#input\_initial\_delay\_seconds) | The initial delay in seconds for the liveness and readiness probes. | `number` | `30` | no |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this service | `string` | n/a | yes |
| <a name="input_max_count"></a> [max\_count](#input\_max\_count) | The maximum number of instances of this task that should be running across our cluster | `number` | `2` | no |
| <a name="input_memory"></a> [memory](#input\_memory) | Memory in megabits per pod | `string` | `"100Mi"` | no |
| <a name="input_period_seconds"></a> [period\_seconds](#input\_period\_seconds) | The period in seconds used for the liveness and readiness probes. | `number` | `3` | no |
| <a name="input_platform_architecture"></a> [platform\_architecture](#input\_platform\_architecture) | The platform to deploy to (valid values: `amd64`, `arm64`). Defaults to `amd64`. | `string` | `"amd64"` | no |
| <a name="input_regional_wafv2_arn"></a> [regional\_wafv2\_arn](#input\_regional\_wafv2\_arn) | A WAF to protect the EKS Ingress if needed | `string` | `null` | no |
| <a name="input_routing"></a> [routing](#input\_routing) | Routing configuration for the ingress | <pre>object({<br>    method : optional(string, "DOMAIN")<br>    host_match : string<br>    group_name : string<br>    alb : optional(object({<br>      name : string,<br>      listener_port : number,<br>    }), null)<br>    priority : number<br>    path : optional(string, "/*")<br>    service_name : string<br>    port : number<br>    service_port : number<br>    service_scheme: optional(string, "HTTP")<br>    scheme : optional(string, "HTTP")<br>    success_codes : optional(string, "200-499")<br>    service_type : string<br>    oidc_config : optional(object({<br>      issuer : string<br>      authorizationEndpoint : string<br>      tokenEndpoint : string<br>      userInfoEndpoint : string<br>      secretName : string<br>      }), {<br>      issuer                = ""<br>      authorizationEndpoint = ""<br>      tokenEndpoint         = ""<br>      userInfoEndpoint      = ""<br>      secretName            = ""<br>    })<br>    bypasses : optional(map(object({<br>      paths   = optional(set(string), [])<br>      methods = optional(set(string), [])<br>    })))<br>  })</pre> | n/a | yes |
| <a name="input_scaling_cpu_threshold_percentage"></a> [scaling\_cpu\_threshold\_percentage](#input\_scaling\_cpu\_threshold\_percentage) | The CPU threshold percentage at which we should scale up | `number` | `80` | no |
| <a name="input_service_endpoints"></a> [service\_endpoints](#input\_service\_endpoints) | Service endpoints to be injected for service discovery | `map(string)` | `{}` | no |
| <a name="input_sidecars"></a> [sidecars](#input\_sidecars) | Map of sidecar containers to be deployed alongside the service | <pre>map(object({<br>    image : string<br>    tag : string<br>    port : optional(number, 80)<br>    scheme : optional(string, "HTTP")<br>    memory : optional(string, "100Mi")<br>    cpu : optional(string, "100m")<br>    image_pull_policy : optional(string, "IfNotPresent")<br>    health_check_path : optional(string, "/")<br>    initial_delay_seconds : optional(number, 30),<br>    period_seconds : optional(number, 3),<br>  }))</pre> | `{}` | no |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Happy Path stack name | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | Standard tags to attach to all happy services | <pre>object({<br>    env : string,<br>    owner : string,<br>    project : string,<br>    service : string,<br>    managedBy : string,<br>  })</pre> | <pre>{<br>  "env": "ADDTAGS",<br>  "managedBy": "ADDTAGS",<br>  "owner": "ADDTAGS",<br>  "project": "ADDTAGS",<br>  "service": "ADDTAGS"<br>}</pre> | no |
| <a name="input_wait_for_steady_state"></a> [wait\_for\_steady\_state](#input\_wait\_for\_steady\_state) | Whether Terraform should block until the service is in a steady state before exiting | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_ecr"></a> [ecr](#output\_ecr) | n/a |
<!-- END -->