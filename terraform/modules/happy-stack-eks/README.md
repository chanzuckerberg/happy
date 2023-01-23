<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
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
| <a name="module_services"></a> [services](#module\_services) | ../happy-service-eks | n/a |
| <a name="module_stack_ingress"></a> [stack\_ingress](#module\_stack\_ingress) | ../happy-ingress-eks | n/a |
| <a name="module_tasks"></a> [tasks](#module\_tasks) | ../happy-task-eks | n/a |

## Resources

| Name | Type |
|------|------|
| [kubernetes_ingress_v1.ingress](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/ingress_v1) | resource |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |
| [kubernetes_secret.integration_secret](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/data-sources/secret) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_account_id"></a> [aws\_account\_id](#input\_aws\_account\_id) | AWS account ID to apply changes to | `string` | `""` | no |
| <a name="input_backend_url"></a> [backend\_url](#input\_backend\_url) | For non-proxied stacks, send in the canonical front/backend URL's | `string` | `""` | no |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | Deployment stage for the app | `string` | n/a | yes |
| <a name="input_frontend_url"></a> [frontend\_url](#input\_frontend\_url) | For non-proxied stacks, send in the canonical front/backend URL's | `string` | `""` | no |
| <a name="input_happy_config_secret"></a> [happy\_config\_secret](#input\_happy\_config\_secret) | Happy Path configuration secret name | `string` | n/a | yes |
| <a name="input_happymeta_"></a> [happymeta\_](#input\_happymeta\_) | Happy Path metadata. Ignored by actual terraform. | `string` | n/a | yes |
| <a name="input_image_tag"></a> [image\_tag](#input\_image\_tag) | Please provide a default image tag | `string` | n/a | yes |
| <a name="input_image_tags"></a> [image\_tags](#input\_image\_tags) | Override image tag for each docker image | `map(string)` | `{}` | no |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this stack | `string` | n/a | yes |
| <a name="input_services"></a> [services](#input\_services) | The services you want to deploy as part of this stack. | <pre>map(object({<br>    name : string,<br>    service_type : string,<br>    desired_count : number,<br>    port : number,<br>    memory : string,<br>    cpu : string,<br>    health_check_path : optional(string, "/"),<br>    aws_iam_policy_json : optional(string, ""),<br>    create_ingress : bool,<br>  }))</pre> | n/a | yes |
| <a name="input_stack_ingress"></a> [stack\_ingress](#input\_stack\_ingress) | The combined ingress object to deploy as part of this stack. | <pre>object({<br>    create_ingress : bool,<br>    service_type : optional(string, "EXTERNAL"),<br>    health_check_path : optional(string, "/*"),<br>    backends : optional(list(object({<br>      path : string,<br>      service_name : string,<br>      service_port : number,<br>    })), [])<br>  })</pre> | n/a | yes |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Happy Path stack name | `string` | n/a | yes |
| <a name="input_stack_prefix"></a> [stack\_prefix](#input\_stack\_prefix) | Do bucket storage paths and db schemas need to be prefixed with the stack name? (Usually '/{stack\_name}' for dev stacks, and '' for staging/prod stacks) | `string` | `""` | no |
| <a name="input_tasks"></a> [tasks](#input\_tasks) | The deletion/migration tasks you want to run when a stack comes up and down. | <pre>map(object({<br>    image : string,<br>    memory : string,<br>    cpu : string,<br>    cmd : set(string),<br>  }))</pre> | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_service_endpoints"></a> [service\_endpoints](#output\_service\_endpoints) | The URL endpoints for services |
| <a name="output_task_arns"></a> [task\_arns](#output\_task\_arns) | ARNs for all the tasks |
<!-- END -->