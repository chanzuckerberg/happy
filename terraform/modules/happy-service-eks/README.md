<!-- START -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |
| <a name="provider_external"></a> [external](#provider\_external) | n/a |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [kubernetes_deployment.deployment](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/deployment) | resource |
| [kubernetes_ingress_v1.ingress](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/ingress_v1) | resource |
| [kubernetes_service.service](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service) | resource |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |
| [external_external.git_repo](https://registry.terraform.io/providers/hashicorp/external/latest/docs/data-sources/external) | data source |
| [external_external.git_sha](https://registry.terraform.io/providers/hashicorp/external/latest/docs/data-sources/external) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_certificate_arn"></a> [certificate\_arn](#input\_certificate\_arn) | ACM certificate ARN to attach to the load balancer listener | `string` | n/a | yes |
| <a name="input_cloud_env"></a> [cloud\_env](#input\_cloud\_env) | Typically data.terraform\_remote\_state.cloud-env.outputs | <pre>object({<br>    public_subnets : list(string),<br>    private_subnets : list(string),<br>    database_subnets : list(string),<br>    database_subnet_group : string,<br>    vpc_id : string,<br>    vpc_cidr_block : string,<br>  })</pre> | n/a | yes |
| <a name="input_container_name"></a> [container\_name](#input\_container\_name) | The name of the container | `string` | n/a | yes |
| <a name="input_cpu"></a> [cpu](#input\_cpu) | CPU shares (1cpu=1000m) per pod | `string` | `"100m"` | no |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | The name of the deployment stage of the Application | `string` | `"dev"` | no |
| <a name="input_desired_count"></a> [desired\_count](#input\_desired\_count) | How many instances of this task should we run across our cluster? | `number` | `2` | no |
| <a name="input_health_check_path"></a> [health\_check\_path](#input\_health\_check\_path) | path to use for health checks | `string` | `"/"` | no |
| <a name="input_host_match"></a> [host\_match](#input\_host\_match) | Host header to match for target rule. Leave empty to match all requests | `string` | n/a | yes |
| <a name="input_image"></a> [image](#input\_image) | Image name | `string` | n/a | yes |
| <a name="input_initial_delay_seconds"></a> [initial\_delay\_seconds](#input\_initial\_delay\_seconds) | The initial delay in seconds for the liveness and readiness probes. | `number` | `30` | no |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this service | `string` | n/a | yes |
| <a name="input_memory"></a> [memory](#input\_memory) | Memory in megabits per pod | `string` | `"100Mi"` | no |
| <a name="input_oauth_certificate_arn"></a> [oauth\_certificate\_arn](#input\_oauth\_certificate\_arn) | Oauth Proxy ACM certificate ARN to attach to the load balancer listener | `string` | n/a | yes |
| <a name="input_path"></a> [path](#input\_path) | The path to register with the Application Load Balancer | `string` | `"/*"` | no |
| <a name="input_period_seconds"></a> [period\_seconds](#input\_period\_seconds) | The period in seconds used for the liveness and readiness probes. | `number` | `3` | no |
| <a name="input_service_endpoints"></a> [service\_endpoints](#input\_service\_endpoints) | Service endpoints to be injected for service discovery | `map(string)` | `{}` | no |
| <a name="input_service_name"></a> [service\_name](#input\_service\_name) | Service name to be deployed | `string` | n/a | yes |
| <a name="input_service_port"></a> [service\_port](#input\_service\_port) | What ports does this service run on? | `number` | `80` | no |
| <a name="input_service_type"></a> [service\_type](#input\_service\_type) | The type of the service to deploy. Supported types include 'EXTERNAL', 'INTERNAL', and 'PRIVATE' | `string` | n/a | yes |
| <a name="input_success_codes"></a> [success\_codes](#input\_success\_codes) | The range of success codes that are used by the ALB ingress controller. | `string` | `"200-499"` | no |
| <a name="input_wait_for_steady_state"></a> [wait\_for\_steady\_state](#input\_wait\_for\_steady\_state) | Whether Terraform should block until the service is in a steady state before exiting | `bool` | `true` | no |

## Outputs

No outputs.
<!-- END -->