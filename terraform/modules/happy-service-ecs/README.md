<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.14 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 5.14 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_cloudwatch_log_group.cloud_watch_datadog_agent_logs_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group) | resource |
| [aws_cloudwatch_log_group.cloud_watch_logs_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group) | resource |
| [aws_ecs_service.service](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service) | resource |
| [aws_ecs_task_definition.task_definition](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_task_definition) | resource |
| [aws_lb_listener_rule.listener_rule](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener_rule) | resource |
| [aws_lb_target_group.target_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group) | resource |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_env_vars"></a> [additional\_env\_vars](#input\_additional\_env\_vars) | Additional environment variables to add to the task definition | `map(string)` | `{}` | no |
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | Please provide the ECS service name | `string` | n/a | yes |
| <a name="input_chamber_service"></a> [chamber\_service](#input\_chamber\_service) | The name of the chamber service from which to load env vars | `string` | `""` | no |
| <a name="input_cluster"></a> [cluster](#input\_cluster) | Please provide the ECS Cluster ID that this service should run on | `string` | n/a | yes |
| <a name="input_cpu"></a> [cpu](#input\_cpu) | CPU shares (1cpu=1024) per task | `number` | `256` | no |
| <a name="input_custom_stack_name"></a> [custom\_stack\_name](#input\_custom\_stack\_name) | Please provide the stack name | `string` | n/a | yes |
| <a name="input_datadog_agent"></a> [datadog\_agent](#input\_datadog\_agent) | DataDog agent image to use | <pre>object({<br>    registry : optional(string, "public.ecr.aws/datadog/agent"),<br>    tag : optional(string, "latest"),<br>    memory : optional(number, 512),<br>    cpu : optional(number, 256),<br>    enabled : optional(bool, false),<br>  })</pre> | <pre>{<br>  "cpu": 256,<br>  "enabled": false,<br>  "memory": 512,<br>  "registry": "public.ecr.aws/datadog/agent",<br>  "tag": "latest"<br>}</pre> | no |
| <a name="input_datadog_api_key"></a> [datadog\_api\_key](#input\_datadog\_api\_key) | DataDog API Key | `string` | `""` | no |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | The name of the deployment stage of the Application | `string` | `"dev"` | no |
| <a name="input_desired_count"></a> [desired\_count](#input\_desired\_count) | How many instances of this task should we run across our cluster? | `number` | `2` | no |
| <a name="input_execution_role"></a> [execution\_role](#input\_execution\_role) | Execution role to use for fargate tasks - required for fargate services! | `string` | `""` | no |
| <a name="input_fail_fast"></a> [fail\_fast](#input\_fail\_fast) | Should containers fail fast if any errors are encountered? | `bool` | `false` | no |
| <a name="input_health_check_path"></a> [health\_check\_path](#input\_health\_check\_path) | path to use for health checks | `string` | `"/"` | no |
| <a name="input_host_match"></a> [host\_match](#input\_host\_match) | Host header to match for target rule. Leave empty to match all requests | `string` | n/a | yes |
| <a name="input_image"></a> [image](#input\_image) | Image name | `string` | n/a | yes |
| <a name="input_launch_type"></a> [launch\_type](#input\_launch\_type) | Launch type on which to run your service. The valid values are EC2, FARGATE, and EXTERNAL | `string` | `"FARGATE"` | no |
| <a name="input_listener"></a> [listener](#input\_listener) | The Application Load Balancer listener to register with | `string` | n/a | yes |
| <a name="input_memory"></a> [memory](#input\_memory) | Memory in megabytes per task | `number` | `1024` | no |
| <a name="input_priority"></a> [priority](#input\_priority) | Listener rule priority number within the given listener | `number` | n/a | yes |
| <a name="input_remote_dev_prefix"></a> [remote\_dev\_prefix](#input\_remote\_dev\_prefix) | S3 storage path / db schema prefix | `string` | `""` | no |
| <a name="input_security_groups"></a> [security\_groups](#input\_security\_groups) | Security groups for ECS tasks | `list(string)` | n/a | yes |
| <a name="input_service_port"></a> [service\_port](#input\_service\_port) | What ports does this service run on? | `number` | `80` | no |
| <a name="input_stack_resource_prefix"></a> [stack\_resource\_prefix](#input\_stack\_resource\_prefix) | Prefix for account-level resources | `string` | n/a | yes |
| <a name="input_subnets"></a> [subnets](#input\_subnets) | Subnets for ecs tasks | `list(string)` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | The happy conventional tags. | <pre>object({<br>    happy_env : string,<br>    happy_stack_name : string,<br>    happy_service_name : string,<br>    happy_region : string,<br>    happy_image : string,<br>    happy_service_type : string,<br>    happy_last_applied : string,<br>  })</pre> | n/a | yes |
| <a name="input_task_role"></a> [task\_role](#input\_task\_role) | ARN and name for the role assumed by tasks | `object({ arn = string, name = string })` | n/a | yes |
| <a name="input_vpc"></a> [vpc](#input\_vpc) | The VPC that the ECS cluster is deployed to | `string` | n/a | yes |
| <a name="input_wait_for_steady_state"></a> [wait\_for\_steady\_state](#input\_wait\_for\_steady\_state) | Whether Terraform should block until the service is in a steady state before exiting | `bool` | `false` | no |

## Outputs

No outputs.
<!-- END -->
