<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.45 |
| <a name="requirement_datadog"></a> [datadog](#requirement\_datadog) | >= 3.20.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 4.45 |
| <a name="provider_datadog"></a> [datadog](#provider\_datadog) | >= 3.20.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_dns"></a> [dns](#module\_dns) | ../happy-dns-ecs | n/a |
| <a name="module_service"></a> [service](#module\_service) | ../happy-service-ecs | n/a |

## Resources

| Name | Type |
|------|------|
| [datadog_synthetics_test.test_api](https://registry.terraform.io/providers/datadog/datadog/latest/docs/resources/synthetics_test) | resource |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |
| [aws_secretsmanager_secret_version.config](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/secretsmanager_secret_version) | data source |
| [datadog_synthetics_locations.locations](https://registry.terraform.io/providers/datadog/datadog/latest/docs/data-sources/synthetics_locations) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_env_vars"></a> [additional\_env\_vars](#input\_additional\_env\_vars) | Additional environment variables to add to the container | `map(string)` | `{}` | no |
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | Please provide the ECS service name | `string` | n/a | yes |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | Deployment stage for the app | `string` | n/a | yes |
| <a name="input_happy_config_secret"></a> [happy\_config\_secret](#input\_happy\_config\_secret) | Happy Path configuration secret name | `string` | n/a | yes |
| <a name="input_image_tag"></a> [image\_tag](#input\_image\_tag) | Please provide a default image tag | `string` | n/a | yes |
| <a name="input_image_tags"></a> [image\_tags](#input\_image\_tags) | Override image tag for each docker image | `map(string)` | `{}` | no |
| <a name="input_launch_type"></a> [launch\_type](#input\_launch\_type) | Launch type on which to run your service. The valid values are EC2 or FARGATE. We strongly suggest Fargate | `string` | `"FARGATE"` | no |
| <a name="input_priority"></a> [priority](#input\_priority) | Listener rule priority number within the given listener | `number` | n/a | yes |
| <a name="input_service_name"></a> [service\_name](#input\_service\_name) | Service name to be deployed | `string` | n/a | yes |
| <a name="input_services"></a> [services](#input\_services) | The services you want to deploy as part of this stack. | <pre>map(object({<br>    name : string,<br>    service_type : string,<br>    desired_count : number,<br>    port : number,<br>    memory : string,<br>    cpu : string,<br>    health_check_path : optional(string, "/"),<br>    #TODO: match the EKS interface aws_iam_policy_json : optional(string, ""),<br>    synthetics : optional(bool, false)<br>  }))</pre> | n/a | yes |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Happy Path stack name | `string` | n/a | yes |
| <a name="input_stack_prefix"></a> [stack\_prefix](#input\_stack\_prefix) | Do bucket storage paths and db schemas need to be prefixed with the stack name? (Usually '/{stack\_name}' for dev stacks, and '' for staging/prod stacks) | `string` | `""` | no |
| <a name="input_tasks"></a> [tasks](#input\_tasks) | The deletion/migration tasks you want to run when a stack comes up and down. | <pre>map(object({<br>    image : string,<br>    memory : string,<br>    cpu : string,<br>    cmd : set(string),<br>  }))</pre> | n/a | yes |
| <a name="input_url"></a> [url](#input\_url) | For non-proxied stacks, send in the canonical front/backend URL's | `string` | `""` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_url"></a> [url](#output\_url) | The URL endpoint for the website service |
<!-- END -->