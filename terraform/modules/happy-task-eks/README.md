<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.14 |
| <a name="requirement_kubernetes"></a> [kubernetes](#requirement\_kubernetes) | >= 2.16 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 5.14 |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | >= 2.16 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_iam_service_account"></a> [iam\_service\_account](#module\_iam\_service\_account) | ../happy-iam-service-account-eks | n/a |

## Resources

| Name | Type |
|------|------|
| [kubernetes_cron_job_v1.task_definition](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/cron_job_v1) | resource |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_env_vars"></a> [additional\_env\_vars](#input\_additional\_env\_vars) | Additional environment variables to add to the task definition | `map(string)` | `{}` | no |
| <a name="input_additional_env_vars_from_config_maps"></a> [additional\_env\_vars\_from\_config\_maps](#input\_additional\_env\_vars\_from\_config\_maps) | Additional environment variables to add to the container from the following config maps | <pre>object({<br>    items : optional(list(string), []),<br>    prefix : optional(string, ""),<br>  })</pre> | <pre>{<br>  "items": [],<br>  "prefix": ""<br>}</pre> | no |
| <a name="input_additional_env_vars_from_secrets"></a> [additional\_env\_vars\_from\_secrets](#input\_additional\_env\_vars\_from\_secrets) | Additional environment variables to add to the container from the following secrets | <pre>object({<br>    items : optional(list(string), []),<br>    prefix : optional(string, ""),<br>  })</pre> | <pre>{<br>  "items": [],<br>  "prefix": ""<br>}</pre> | no |
| <a name="input_additional_volumes_from_config_maps"></a> [additional\_volumes\_from\_config\_maps](#input\_additional\_volumes\_from\_config\_maps) | Additional volumes to add to the container from the following config maps | <pre>object({<br>    items : optional(list(string), []),<br>  })</pre> | <pre>{<br>  "items": []<br>}</pre> | no |
| <a name="input_additional_volumes_from_secrets"></a> [additional\_volumes\_from\_secrets](#input\_additional\_volumes\_from\_secrets) | Additional volumes to add to the container from the following secrets | <pre>object({<br>    items : optional(list(string), []),<br>    base_dir : optional(string, "/var"),<br>  })</pre> | <pre>{<br>  "base_dir": "/var",<br>  "items": []<br>}</pre> | no |
| <a name="input_args"></a> [args](#input\_args) | Args to pass to the command | `list(string)` | `[]` | no |
| <a name="input_aws_iam"></a> [aws\_iam](#input\_aws\_iam) | The AWS IAM service account or policy JSON to give to the pod. Only one of these should be set. | <pre>object({<br>    service_account_name : optional(string, null),<br>    policy_json : optional(string, ""),<br>  })</pre> | `{}` | no |
| <a name="input_backoff_limit"></a> [backoff\_limit](#input\_backoff\_limit) | kubernetes\_cron\_job backoff\_limit | `number` | `2` | no |
| <a name="input_cmd"></a> [cmd](#input\_cmd) | Command to run | `list(string)` | `[]` | no |
| <a name="input_cpu"></a> [cpu](#input\_cpu) | CPU shares (1cpu=1000m) per pod | `string` | `"100m"` | no |
| <a name="input_cpu_requests"></a> [cpu\_requests](#input\_cpu\_requests) | CPU shares (1cpu=1000m) requested per pod | `string` | `"10m"` | no |
| <a name="input_cron_schedule"></a> [cron\_schedule](#input\_cron\_schedule) | Cron schedule for this job | `string` | `"0 0 1 1 *"` | no |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | The name of the deployment stage of the Application | `string` | n/a | yes |
| <a name="input_eks_cluster"></a> [eks\_cluster](#input\_eks\_cluster) | eks-cluster module output | <pre>object({<br>    cluster_id : string,<br>    cluster_arn : string,<br>    cluster_endpoint : string,<br>    cluster_ca : string,<br>    cluster_oidc_issuer_url : string,<br>    cluster_version : string,<br>    worker_iam_role_name : string,<br>    worker_security_group : string,<br>    oidc_provider_arn : string,<br>  })</pre> | n/a | yes |
| <a name="input_failed_jobs_history_limit"></a> [failed\_jobs\_history\_limit](#input\_failed\_jobs\_history\_limit) | kubernetes\_cron\_job failed jobs history limit | `number` | `5` | no |
| <a name="input_image"></a> [image](#input\_image) | Image name | `string` | n/a | yes |
| <a name="input_is_cron_job"></a> [is\_cron\_job](#input\_is\_cron\_job) | Indicates if this job should be run on a schedule or one-off. If true, set cron\_schedule as well | `bool` | `false` | no |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this task | `string` | n/a | yes |
| <a name="input_memory"></a> [memory](#input\_memory) | Memory in megabits per pod | `string` | `"100Mi"` | no |
| <a name="input_memory_requests"></a> [memory\_requests](#input\_memory\_requests) | Memory requests per pod | `string` | `"10Mi"` | no |
| <a name="input_platform_architecture"></a> [platform\_architecture](#input\_platform\_architecture) | Platform architecture | `string` | `"amd64"` | no |
| <a name="input_remote_dev_prefix"></a> [remote\_dev\_prefix](#input\_remote\_dev\_prefix) | S3 storage path / db schema prefix | `string` | `""` | no |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Happy Path stack name | `string` | n/a | yes |
| <a name="input_starting_deadline_seconds"></a> [starting\_deadline\_seconds](#input\_starting\_deadline\_seconds) | kubernetes\_cron\_job starting\_deadline\_seconds | `number` | `30` | no |
| <a name="input_successful_jobs_history_limit"></a> [successful\_jobs\_history\_limit](#input\_successful\_jobs\_history\_limit) | kubernetes\_cron\_job successful\_jobs\_history\_limit | `number` | `5` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Standard tags to attach to all happy services | <pre>object({<br>    env : string,<br>    owner : string,<br>    project : string,<br>    service : string,<br>    managedBy : string,<br>  })</pre> | <pre>{<br>  "env": "ADDTAGS",<br>  "managedBy": "ADDTAGS",<br>  "owner": "ADDTAGS",<br>  "project": "ADDTAGS",<br>  "service": "ADDTAGS"<br>}</pre> | no |
| <a name="input_task_name"></a> [task\_name](#input\_task\_name) | Happy Path task name | `string` | n/a | yes |
| <a name="input_ttl_seconds_after_finished"></a> [ttl\_seconds\_after\_finished](#input\_ttl\_seconds\_after\_finished) | kubernetes\_cron\_job ttl\_seconds\_after\_finished | `number` | `10` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_task_definition_arn"></a> [task\_definition\_arn](#output\_task\_definition\_arn) | Task definition name |
<!-- END -->
<!-- BEGIN_TF_DOCS -->
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

No modules.

## Resources

| Name | Type |
|------|------|
| [kubernetes_cron_job_v1.task_definition](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/cron_job_v1) | resource |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_env_vars"></a> [additional\_env\_vars](#input\_additional\_env\_vars) | Additional environment variables to add to the task definition | `map(string)` | `{}` | no |
| <a name="input_backoff_limit"></a> [backoff\_limit](#input\_backoff\_limit) | kubernetes\_cron\_job backoff\_limit | `number` | `2` | no |
| <a name="input_cmd"></a> [cmd](#input\_cmd) | Command to run | `list(string)` | `[]` | no |
| <a name="input_cpu"></a> [cpu](#input\_cpu) | CPU shares (1cpu=1000m) per pod | `string` | `"100m"` | no |
| <a name="input_cpu_requests"></a> [cpu\_requests](#input\_cpu\_requests) | CPU shares (1cpu=1000m) requested per pod | `string` | `"10m"` | no |
| <a name="input_cron_schedule"></a> [cron\_schedule](#input\_cron\_schedule) | Cron schedule for this job | `string` | `"0 0 1 1 *"` | no |
| <a name="input_deployment_stage"></a> [deployment\_stage](#input\_deployment\_stage) | The name of the deployment stage of the Application | `string` | n/a | yes |
| <a name="input_failed_jobs_history_limit"></a> [failed\_jobs\_history\_limit](#input\_failed\_jobs\_history\_limit) | kubernetes\_cron\_job failed jobs history limit | `number` | `5` | no |
| <a name="input_image"></a> [image](#input\_image) | Image name | `string` | n/a | yes |
| <a name="input_is_cron_job"></a> [is\_cron\_job](#input\_is\_cron\_job) | Indicates if this job should be run on a schedule or one-off. If true, set cron\_schedule as well | `bool` | `false` | no |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | K8S namespace for this task | `string` | n/a | yes |
| <a name="input_memory"></a> [memory](#input\_memory) | Memory in megabits per pod | `string` | `"100Mi"` | no |
| <a name="input_memory_requests"></a> [memory\_requests](#input\_memory\_requests) | Memory requests per pod | `string` | `"10Mi"` | no |
| <a name="input_platform_architecture"></a> [platform\_architecture](#input\_platform\_architecture) | Platform architecture | `string` | `"amd64"` | no |
| <a name="input_remote_dev_prefix"></a> [remote\_dev\_prefix](#input\_remote\_dev\_prefix) | S3 storage path / db schema prefix | `string` | `""` | no |
| <a name="input_stack_name"></a> [stack\_name](#input\_stack\_name) | Happy Path stack name | `string` | n/a | yes |
| <a name="input_starting_deadline_seconds"></a> [starting\_deadline\_seconds](#input\_starting\_deadline\_seconds) | kubernetes\_cron\_job starting\_deadline\_seconds | `number` | `30` | no |
| <a name="input_successful_jobs_history_limit"></a> [successful\_jobs\_history\_limit](#input\_successful\_jobs\_history\_limit) | kubernetes\_cron\_job successful\_jobs\_history\_limit | `number` | `5` | no |
| <a name="input_task_name"></a> [task\_name](#input\_task\_name) | Happy Path task name | `string` | n/a | yes |
| <a name="input_ttl_seconds_after_finished"></a> [ttl\_seconds\_after\_finished](#input\_ttl\_seconds\_after\_finished) | kubernetes\_cron\_job ttl\_seconds\_after\_finished | `number` | `10` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_task_definition_arn"></a> [task\_definition\_arn](#output\_task\_definition\_arn) | Task definition name |
<!-- END_TF_DOCS -->