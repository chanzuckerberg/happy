<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.45 |
| <a name="requirement_kubernetes"></a> [kubernetes](#requirement\_kubernetes) | >= 2.16 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 4.45 |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | >= 2.16 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_iam_service_account"></a> [iam\_service\_account](#module\_iam\_service\_account) | ../happy-iam-service-account-eks | n/a |
| <a name="module_ingress"></a> [ingress](#module\_ingress) | ../happy-ingress-eks | n/a |

## Resources

| Name | Type |
|------|------|
| [kubernetes_deployment.deployment](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/deployment) | resource |
| [kubernetes_service.service](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service) | resource |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_env_vars"></a> [additional\_env\_vars](#input\_additional\_env\_vars) | Additional environment variables to add to the task definition | `map(string)` | `{}` | no |
| <a name="input_aws_iam_policy_json"></a> [aws\_iam\_policy\_json](#input\_aws\_iam\_policy\_json) | The AWS IAM policy to give to the pod. | `string` | `""` | no |
| <a name="input_certificate_arn"></a> [certificate\_arn](#input\_certificate\_arn) | ACM certificate ARN to attach to the load balancer listener | `string` | n/a | yes |
| <a name="input_cloud_env"></a> [cloud\_env](#input\_cloud\_env) | Typically data.terraform\_remote\_state.cloud-env.outputs | <pre>object({<br>    public_subnets : list(string),<br>    private_subnets : list(string),<br>    database_subnets : list(string),<br>    database_subnet_group : string,<br>    vpc_id : string,<br>    vpc_cidr_block : string,<br>  })</pre> | n/a | yes |
| <a name="input_container_name"></a> [container\_name](#input\_container\_name) | The name of the container | `string` | n/a | yes |
| <a name="input_cpu"></a> [cpu](#input\_cpu) | CPU shares (1cpu=1000m) per pod | `string` | `"100m"` | no |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | The name of the deployment stage of the Application | `string` | `"dev"` | no |
| <a name="input_desired_count"></a> [desired\_count](#input\_desired\_count) | How many instances of this task should we run across our cluster? | `number` | `2` | no |
| <a name="input_eks_cluster"></a> [eks\_cluster](#input\_eks\_cluster) | eks-cluster module output | <pre>object({<br>    cluster_id : string,<br>    cluster_arn : string,<br>    cluster_endpoint : string,<br>    cluster_ca : string,<br>    cluster_oidc_issuer_url : string,<br>    cluster_security_group : string,<br>    cluster_iam_role_name : string,<br>    cluster_version : string,<br>    worker_iam_role_name : string,<br>    kubeconfig : string,<br>    worker_security_group : string,<br>    oidc_provider_arn : string,<br>  })</pre> | n/a | yes |
| <a name="input_health_check_path"></a> [health\_check\_path](#input\_health\_check\_path) | path to use for health checks | `string` | `"/"` | no |
| <a name="input_image"></a> [image](#input\_image) | Image name | `string` | n/a | yes |
| <a name="input_initial_delay_seconds"></a> [initial\_delay\_seconds](#input\_initial\_delay\_seconds) | The initial delay in seconds for the liveness and readiness probes. | `number` | `30` | no |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this service | `string` | n/a | yes |
| <a name="input_memory"></a> [memory](#input\_memory) | Memory in megabits per pod | `string` | `"100Mi"` | no |
| <a name="input_oauth_certificate_arn"></a> [oauth\_certificate\_arn](#input\_oauth\_certificate\_arn) | Oauth Proxy ACM certificate ARN to attach to the load balancer listener | `string` | n/a | yes |
| <a name="input_period_seconds"></a> [period\_seconds](#input\_period\_seconds) | The period in seconds used for the liveness and readiness probes. | `number` | `3` | no |
| <a name="input_routing"></a> [routing](#input\_routing) | Routing configuration for the ingress | <pre>object({<br>    method : optional(string, "CONTEXT")<br>    host_match : string<br>    group_name : string<br>    priority : number<br>    path : optional(string, "/*")<br>    service_name : string<br>    service_port : number<br>  })</pre> | n/a | yes |
| <a name="input_service_endpoints"></a> [service\_endpoints](#input\_service\_endpoints) | Service endpoints to be injected for service discovery | `map(string)` | `{}` | no |
| <a name="input_service_type"></a> [service\_type](#input\_service\_type) | The type of the service to deploy. Supported types include 'EXTERNAL', 'INTERNAL', and 'PRIVATE' | `string` | n/a | yes |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Happy Path stack name | `string` | n/a | yes |
| <a name="input_success_codes"></a> [success\_codes](#input\_success\_codes) | The range of success codes that are used by the ALB ingress controller. | `string` | `"200-499"` | no |
| <a name="input_wait_for_steady_state"></a> [wait\_for\_steady\_state](#input\_wait\_for\_steady\_state) | Whether Terraform should block until the service is in a steady state before exiting | `bool` | `true` | no |

## Outputs

No outputs.
<!-- END -->