<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.45 |
| <a name="requirement_datadog"></a> [datadog](#requirement\_datadog) | >= 3.20.0 |
| <a name="requirement_happy"></a> [happy](#requirement\_happy) | >= 0.53.5 |
| <a name="requirement_kubernetes"></a> [kubernetes](#requirement\_kubernetes) | >= 2.16 |
| <a name="requirement_validation"></a> [validation](#requirement\_validation) | 1.0.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_datadog"></a> [datadog](#provider\_datadog) | >= 3.20.0 |
| <a name="provider_happy"></a> [happy](#provider\_happy) | >= 0.53.5 |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | >= 2.16 |
| <a name="provider_validation"></a> [validation](#provider\_validation) | 1.0.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_services"></a> [services](#module\_services) | ../happy-service-eks | n/a |
| <a name="module_tasks"></a> [tasks](#module\_tasks) | ../happy-task-eks | n/a |

## Resources

| Name | Type |
|------|------|
| [datadog_dashboard_json.stack_dashboard](https://registry.terraform.io/providers/datadog/datadog/latest/docs/resources/dashboard_json) | resource |
| [datadog_synthetics_test.test_api](https://registry.terraform.io/providers/datadog/datadog/latest/docs/resources/synthetics_test) | resource |
| [kubernetes_secret.oidc_config](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/secret) | resource |
| [validation_error.mix_of_internal_and_external_services](https://registry.terraform.io/providers/tlkamp/validation/1.0.0/docs/resources/error) | resource |
| [datadog_synthetics_locations.locations](https://registry.terraform.io/providers/datadog/datadog/latest/docs/data-sources/synthetics_locations) | data source |
| [happy_resolved_app_configs.configs](https://registry.terraform.io/providers/chanzuckerberg/happy/latest/docs/data-sources/resolved_app_configs) | data source |
| [kubernetes_secret.integration_secret](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/data-sources/secret) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_env_vars"></a> [additional\_env\_vars](#input\_additional\_env\_vars) | Additional environment variables to add to the container | `map(string)` | `{}` | no |
| <a name="input_additional_env_vars_from_config_maps"></a> [additional\_env\_vars\_from\_config\_maps](#input\_additional\_env\_vars\_from\_config\_maps) | Additional environment variables to add to the container from the following config maps | <pre>object({<br>    items : optional(list(string), []),<br>    prefix : optional(string, ""),<br>  })</pre> | <pre>{<br>  "items": [],<br>  "prefix": ""<br>}</pre> | no |
| <a name="input_additional_env_vars_from_secrets"></a> [additional\_env\_vars\_from\_secrets](#input\_additional\_env\_vars\_from\_secrets) | Additional environment variables to add to the container from the following secrets | <pre>object({<br>    items : optional(list(string), []),<br>    prefix : optional(string, ""),<br>  })</pre> | <pre>{<br>  "items": [],<br>  "prefix": ""<br>}</pre> | no |
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | The happy application name | `string` | `""` | no |
| <a name="input_create_dashboard"></a> [create\_dashboard](#input\_create\_dashboard) | Create a dashboard for this stack | `bool` | `false` | no |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | Deployment stage for the app | `string` | n/a | yes |
| <a name="input_image_tag"></a> [image\_tag](#input\_image\_tag) | Please provide a default image tag | `string` | n/a | yes |
| <a name="input_image_tags"></a> [image\_tags](#input\_image\_tags) | Override image tag for each docker image | `map(string)` | `{}` | no |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this stack | `string` | n/a | yes |
| <a name="input_routing_method"></a> [routing\_method](#input\_routing\_method) | Traffic routing method for this stack. Valid options are 'DOMAIN', when every service gets a unique domain name, or a 'CONTEXT' when all services share the same domain name, and routing is done by request path. | `string` | `"DOMAIN"` | no |
| <a name="input_services"></a> [services](#input\_services) | The services you want to deploy as part of this stack. | <pre>map(object({<br>    name : string,<br>    service_type : string, // oneof: EXTERNAL, INTERNAL, PRIVATE<br>    desired_count : optional(number, 2),<br>    max_count : optional(number, 2),<br>    scaling_cpu_threshold_percentage : optional(number, 80),<br>    port : number,<br>    memory : string,<br>    cpu : string,<br>    health_check_path : optional(string, "/"),<br>    aws_iam_policy_json : optional(string, ""),<br>    path : optional(string, "/*"),  // Only used for CONTEXT routing<br>    priority : optional(number, 0), // Only used for CONTEXT routing<br>    success_codes : optional(string, "200-499"),<br>    synthetics : optional(bool, false),<br>    initial_delay_seconds : optional(number, 30),<br>    period_seconds : optional(number, 3),<br>    bypasses : optional(map(object({ // Only used for INTERNAL service_type<br>      paths   = optional(set(string), [])<br>      methods = optional(set(string), [])<br>    })), {})<br>  }))</pre> | n/a | yes |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Happy Path stack name | `string` | n/a | yes |
| <a name="input_stack_prefix"></a> [stack\_prefix](#input\_stack\_prefix) | Do bucket storage paths and db schemas need to be prefixed with the stack name? (Usually '/{stack\_name}' for dev stacks, and '' for staging/prod stacks) | `string` | `""` | no |
| <a name="input_tasks"></a> [tasks](#input\_tasks) | The deletion/migration tasks you want to run when a stack comes up and down. | <pre>map(object({<br>    image : string,<br>    memory : string,<br>    cpu : string,<br>    cmd : set(string),<br>  }))</pre> | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_dashboard"></a> [dashboard](#output\_dashboard) | n/a |
| <a name="output_service_endpoints"></a> [service\_endpoints](#output\_service\_endpoints) | The URL endpoints for services |
| <a name="output_task_arns"></a> [task\_arns](#output\_task\_arns) | ARNs for all the tasks |
<!-- END -->

