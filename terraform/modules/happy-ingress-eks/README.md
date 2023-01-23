<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_kubernetes"></a> [kubernetes](#requirement\_kubernetes) | >= 2.16 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | >= 2.16 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [kubernetes_ingress_v1.ingress](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/ingress_v1) | resource |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_backends"></a> [backends](#input\_backends) | The backends to register with the Application Load Balancer | <pre>list(object(<br>    {<br>      service_name : string<br>      service_port : number<br>      path : string<br>    }<br>  ))</pre> | `[]` | no |
| <a name="input_certificate_arn"></a> [certificate\_arn](#input\_certificate\_arn) | ACM certificate ARN to attach to the load balancer listener | `string` | n/a | yes |
| <a name="input_cloud_env"></a> [cloud\_env](#input\_cloud\_env) | Typically data.terraform\_remote\_state.cloud-env.outputs | <pre>object({<br>    public_subnets : list(string),<br>    private_subnets : list(string),<br>    database_subnets : list(string),<br>    database_subnet_group : string,<br>    vpc_id : string,<br>    vpc_cidr_block : string,<br>  })</pre> | n/a | yes |
| <a name="input_health_check_path"></a> [health\_check\_path](#input\_health\_check\_path) | path to use for health checks | `string` | `"/"` | no |
| <a name="input_host_match"></a> [host\_match](#input\_host\_match) | Host header to match for target rule. Leave empty to match all requests | `string` | n/a | yes |
| <a name="input_ingress_name"></a> [ingress\_name](#input\_ingress\_name) | Name of the ingress resource | `string` | n/a | yes |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this service | `string` | n/a | yes |
| <a name="input_service_type"></a> [service\_type](#input\_service\_type) | The type of the service to deploy. Supported types include 'EXTERNAL', 'INTERNAL', and 'PRIVATE' | `string` | n/a | yes |
| <a name="input_success_codes"></a> [success\_codes](#input\_success\_codes) | The range of success codes that are used by the ALB ingress controller. | `string` | `"200-499"` | no |
| <a name="input_tags_string"></a> [tags\_string](#input\_tags\_string) | Tags to apply to ingress resource, comma delimited key=value pairs | `string` | `""` | no |

## Outputs

No outputs.
<!-- END -->