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
| [aws_security_group.alb_sg](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group) | resource |
| [kubernetes_ingress_v1.ingress](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/ingress_v1) | resource |
| [kubernetes_ingress_v1.ingress_bypasses](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/ingress_v1) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_alb_healthcheck_interval_seconds"></a> [aws\_alb\_healthcheck\_interval\_seconds](#input\_aws\_alb\_healthcheck\_interval\_seconds) | The time in seconds to ping the target group for a health check; defaults to a high numbers since k8s also has a healthcheck | `string` | `"300"` | no |
| <a name="input_certificate_arn"></a> [certificate\_arn](#input\_certificate\_arn) | ACM certificate ARN to attach to the load balancer listener | `string` | n/a | yes |
| <a name="input_cloud_env"></a> [cloud\_env](#input\_cloud\_env) | Typically data.terraform\_remote\_state.cloud-env.outputs | <pre>object({<br>    public_subnets : list(string),<br>    private_subnets : list(string),<br>    database_subnets : list(string),<br>    database_subnet_group : string,<br>    vpc_id : string,<br>    vpc_cidr_block : string,<br>  })</pre> | n/a | yes |
| <a name="input_health_check_path"></a> [health\_check\_path](#input\_health\_check\_path) | path to use for health checks | `string` | `"/"` | no |
| <a name="input_ingress_name"></a> [ingress\_name](#input\_ingress\_name) | Name of the ingress resource | `string` | n/a | yes |
| <a name="input_ingress_security_groups"></a> [ingress\_security\_groups](#input\_ingress\_security\_groups) | A list of security groups that should be allowed to communicate with this ingress. | `list(string)` | `[]` | no |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this service | `string` | n/a | yes |
| <a name="input_labels"></a> [labels](#input\_labels) | Labels to apply to ingress resource | `map(string)` | n/a | yes |
| <a name="input_regional_wafv2_arn"></a> [regional\_wafv2\_arn](#input\_regional\_wafv2\_arn) | A WAF to protect the EKS Ingress if needed | `string` | `null` | no |
| <a name="input_routing"></a> [routing](#input\_routing) | Routing configuration for the ingress | <pre>object({<br>    method : optional(string, "DOMAIN")<br>    host_match : string<br>    group_name : string<br>    priority : number<br>    path : optional(string, "/*")<br>    service_name : string<br>    service_port : number<br>    service_scheme : string<br>    service_type : string<br>    idle_timeout : optional(number, 60)<br>    oidc_config : optional(object({<br>      issuer : string<br>      authorizationEndpoint : string<br>      tokenEndpoint : string<br>      userInfoEndpoint : string<br>      secretName : string<br>      }), {<br>      issuer                = ""<br>      authorizationEndpoint = ""<br>      tokenEndpoint         = ""<br>      userInfoEndpoint      = ""<br>      secretName            = ""<br>    })<br>    bypasses : optional(map(object({<br>      paths   = optional(set(string), [])<br>      methods = optional(set(string), [])<br>    })))<br>    success_codes : optional(string, "200-499")<br>  })</pre> | n/a | yes |
| <a name="input_tags_string"></a> [tags\_string](#input\_tags\_string) | Tags to apply to ingress resource, comma delimited key=value pairs | `string` | `""` | no |
| <a name="input_target_service_name"></a> [target\_service\_name](#input\_target\_service\_name) | Name of destination service that the ingress should route to | `string` | n/a | yes |
| <a name="input_target_service_port"></a> [target\_service\_port](#input\_target\_service\_port) | Port of destination service that the ingress should route to | `number` | n/a | yes |
| <a name="input_target_service_scheme"></a> [target\_service\_scheme](#input\_target\_service\_scheme) | Scheme of destination service that the ingress should route to | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END -->