<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.45 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 4.45 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_dns"></a> [dns](#module\_dns) | ../happy-dns-ecs | n/a |
| <a name="module_service"></a> [service](#module\_service) | ../happy-service-ecs | n/a |

## Resources

| Name | Type |
|------|------|
| [aws_secretsmanager_secret_version.config](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/secretsmanager_secret_version) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | Please provide the ECS service name | `string` | n/a | yes |
| <a name="input_chamber_service"></a> [chamber\_service](#input\_chamber\_service) | The name of the chamber service from which to load env vars | `string` | `""` | no |
| <a name="input_cpu"></a> [cpu](#input\_cpu) | CPU shares (1cpu=1024) per task | `number` | `256` | no |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | Deployment stage for the app | `string` | n/a | yes |
| <a name="input_desired_count"></a> [desired\_count](#input\_desired\_count) | How many instances of this task should we run across our cluster? | `number` | `2` | no |
| <a name="input_happy_config_secret"></a> [happy\_config\_secret](#input\_happy\_config\_secret) | Happy Path configuration secret name | `string` | n/a | yes |
| <a name="input_image_tag"></a> [image\_tag](#input\_image\_tag) | Please provide a default image tag | `string` | n/a | yes |
| <a name="input_image_tags"></a> [image\_tags](#input\_image\_tags) | Override image tag for each docker image | `map(string)` | `{}` | no |
| <a name="input_launch_type"></a> [launch\_type](#input\_launch\_type) | Launch type on which to run your service. The valid values are EC2, FARGATE, and EXTERNAL | `string` | `"FARGATE"` | no |
| <a name="input_memory"></a> [memory](#input\_memory) | Memory in megabytes per task | `number` | `1024` | no |
| <a name="input_priority"></a> [priority](#input\_priority) | Listener rule priority number within the given listener | `number` | n/a | yes |
| <a name="input_require_okta"></a> [require\_okta](#input\_require\_okta) | Whether the ALB's should be on private subnets | `bool` | `true` | no |
| <a name="input_service_port"></a> [service\_port](#input\_service\_port) | What ports does this service run on? | `number` | `80` | no |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Happy Path stack name | `string` | n/a | yes |
| <a name="input_stack_prefix"></a> [stack\_prefix](#input\_stack\_prefix) | Do bucket storage paths and db schemas need to be prefixed with the stack name? (Usually '/{stack\_name}' for dev stacks, and '' for staging/prod stacks) | `string` | `""` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | The happy conventional tags. | <pre>object({<br>    happy_env : string,<br>    happy_stack_name : string,<br>    happy_service_name : string,<br>    happy_region : string,<br>    happy_image : string,<br>    happy_service_type : string,<br>    happy_last_applied : string,<br>  })</pre> | n/a | yes |
| <a name="input_url"></a> [url](#input\_url) | For non-proxied stacks, send in the canonical front/backend URL's | `string` | `""` | no |
| <a name="input_wait_for_steady_state"></a> [wait\_for\_steady\_state](#input\_wait\_for\_steady\_state) | Should terraform block until ECS services reach a steady state? | `bool` | `false` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_url"></a> [url](#output\_url) | The URL endpoint for the website service |
<!-- END -->