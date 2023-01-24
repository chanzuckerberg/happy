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

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_cloudwatch_log_group.cloud_watch_logs_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group) | resource |
| [aws_ecs_service.service](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service) | resource |
| [aws_ecs_task_definition.task_definition](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_task_definition) | resource |
| [aws_lb_listener_rule.listener_rule](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener_rule) | resource |
| [aws_lb_target_group.target_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group) | resource |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_env_vars"></a> [additional\_env\_vars](#input\_additional\_env\_vars) | Additional environment variables to add to the task definition | `list(object({ name : string, value : string }))` | `[]` | no |
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | Please provide the ECS service name | `string` | n/a | yes |
| <a name="input_cloud_env"></a> [cloud\_env](#input\_cloud\_env) | Typically data.terraform\_remote\_state.cloud-env.outputs | <pre>object({<br>    public_subnets : list(string),<br>    private_subnets : list(string),<br>    database_subnets : list(string),<br>    database_subnet_group : string,<br>    vpc_id : string,<br>    vpc_cidr_block : string,<br>  })</pre> | n/a | yes |
| <a name="input_cluster"></a> [cluster](#input\_cluster) | Please provide the ECS Cluster ID that this service should run on | `string` | n/a | yes |
| <a name="input_cpu"></a> [cpu](#input\_cpu) | CPU shares (1cpu=1024) per task | `number` | `256` | no |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | The name of the deployment stage of the Application | `string` | `"dev"` | no |
| <a name="input_desired_count"></a> [desired\_count](#input\_desired\_count) | How many instances of this task should we run across our cluster? | `number` | `2` | no |
| <a name="input_execution_role"></a> [execution\_role](#input\_execution\_role) | Execution role to use for fargate tasks - required for fargate services! | `string` | `""` | no |
| <a name="input_health_check_path"></a> [health\_check\_path](#input\_health\_check\_path) | path to use for health checks | `string` | `"/"` | no |
| <a name="input_host_match"></a> [host\_match](#input\_host\_match) | Host header to match for target rule. Leave empty to match all requests | `string` | n/a | yes |
| <a name="input_image"></a> [image](#input\_image) | Image name | `string` | n/a | yes |
| <a name="input_launch_type"></a> [launch\_type](#input\_launch\_type) | Launch type on which to run your service. The valid values are EC2, FARGATE, and EXTERNAL | `string` | `"FARGATE"` | no |
| <a name="input_listener"></a> [listener](#input\_listener) | The Application Load Balancer listener to register with | `string` | n/a | yes |
| <a name="input_memory"></a> [memory](#input\_memory) | Memory in megabytes per task | `number` | `1024` | no |
| <a name="input_priority"></a> [priority](#input\_priority) | Listener rule priority number within the given listener | `number` | n/a | yes |
| <a name="input_security_groups"></a> [security\_groups](#input\_security\_groups) | Security groups for ECS tasks | `list(string)` | n/a | yes |
| <a name="input_service_port"></a> [service\_port](#input\_service\_port) | What ports does this service run on? | `number` | `80` | no |
| <a name="input_service_type"></a> [service\_type](#input\_service\_type) | The type of the service to deploy. Supported types include 'EXTERNAL', 'INTERNAL', and 'PRIVATE' | `string` | n/a | yes |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Please provide the stack name | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | Standard tags to attach to all happy services | <pre>object({<br>    env : string,<br>    owner : string,<br>    project : string,<br>    service : string,<br>    managedBy : string,<br>  })</pre> | <pre>{<br>  "env": "ADDTAGS",<br>  "managedBy": "ADDTAGS",<br>  "owner": "ADDTAGS",<br>  "project": "ADDTAGS",<br>  "service": "ADDTAGS"<br>}</pre> | no |
| <a name="input_task_role"></a> [task\_role](#input\_task\_role) | ARN and name for the role assumed by tasks | `object({ arn = string, name = string })` | n/a | yes |
| <a name="input_wait_for_steady_state"></a> [wait\_for\_steady\_state](#input\_wait\_for\_steady\_state) | Whether Terraform should block until the service is in a steady state before exiting | `bool` | `false` | no |

## Outputs

No outputs.
<!-- END -->
