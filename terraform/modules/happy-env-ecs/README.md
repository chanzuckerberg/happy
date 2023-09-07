Default happy path environment module that supports creating S3 buckets, RDS databases, ECS clusters, Batch environments, authenticated remote-dev environments, and more!
// bump
<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.14 |
| <a name="requirement_random"></a> [random](#requirement\_random) | >= 3.4 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 5.14 |
| <a name="provider_random"></a> [random](#provider\_random) | >= 3.4 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_batch"></a> [batch](#module\_batch) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-batch-env | v0.292.0 |
| <a name="module_batch-swipe"></a> [batch-swipe](#module\_batch-swipe) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-batch-env-swipe | v0.227.0 |
| <a name="module_cert-lb"></a> [cert-lb](#module\_cert-lb) | github.com/chanzuckerberg/cztack//aws-acm-certificate | v0.43.1 |
| <a name="module_db"></a> [db](#module\_db) | github.com/chanzuckerberg/cztack//aws-aurora-postgres | v0.49.0 |
| <a name="module_ecrs"></a> [ecrs](#module\_ecrs) | git@github.com:chanzuckerberg/cztack//aws-ecr-repo | v0.59.0 |
| <a name="module_ecs-cluster"></a> [ecs-cluster](#module\_ecs-cluster) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/ecs-cluster | ecs-cluster-v1.0.1 |
| <a name="module_ecs-multi-domain-oauth-proxy"></a> [ecs-multi-domain-oauth-proxy](#module\_ecs-multi-domain-oauth-proxy) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/ecs-multi-domain-oauth-proxy | ecs-multi-domain-oauth-proxy-v1.3.4 |
| <a name="module_happy_github_ci_role"></a> [happy\_github\_ci\_role](#module\_happy\_github\_ci\_role) | ../happy-github-ci-role | n/a |
| <a name="module_happy_service_account"></a> [happy\_service\_account](#module\_happy\_service\_account) | ../happy-tfe-okta-service-account | n/a |
| <a name="module_instance-cloud-init-script"></a> [instance-cloud-init-script](#module\_instance-cloud-init-script) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/instance-cloud-init-script | v0.227.0 |
| <a name="module_s3_bucket"></a> [s3\_bucket](#module\_s3\_bucket) | github.com/chanzuckerberg/cztack//aws-s3-private-bucket | v0.56.2 |
| <a name="module_swipe"></a> [swipe](#module\_swipe) | git@github.com:chanzuckerberg/swipe | v1.2.1 |

## Resources

| Name | Type |
|------|------|
| [aws_autoscaling_policy.scale-down](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/autoscaling_policy) | resource |
| [aws_autoscaling_policy.scale-up](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/autoscaling_policy) | resource |
| [aws_cloudwatch_log_group.ecs](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group) | resource |
| [aws_cloudwatch_metric_alarm.memory-res-high](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_metric_alarm) | resource |
| [aws_cloudwatch_metric_alarm.memory-res-low](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_metric_alarm) | resource |
| [aws_dynamodb_table.locks](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table) | resource |
| [aws_iam_policy.locktable_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_role.proxy_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role.task_execution_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy.ecs_reader](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_lb.lb-private](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb) | resource |
| [aws_lb.lb-public](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb) | resource |
| [aws_lb_listener.private-lb-listener](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener) | resource |
| [aws_lb_listener.public-http](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener) | resource |
| [aws_lb_listener.public-https](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener) | resource |
| [aws_route53_record.happy-NS](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_route53_record.services](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_route53_zone.happy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_zone) | resource |
| [aws_secretsmanager_secret.happy_env_secret](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/secretsmanager_secret) | resource |
| [aws_secretsmanager_secret_version.happy_env_secret_version](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/secretsmanager_secret_version) | resource |
| [aws_security_group.happy_env_sg](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group) | resource |
| [aws_security_group.public_alb_sg](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group) | resource |
| [aws_security_group_rule.alb](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group_rule) | resource |
| [aws_security_group_rule.egress](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group_rule) | resource |
| [aws_security_group_rule.oauth](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group_rule) | resource |
| [aws_security_group_rule.tasks](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group_rule) | resource |
| [aws_wafv2_web_acl_association.private](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl_association) | resource |
| [aws_wafv2_web_acl_association.public](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl_association) | resource |
| [random_password.db-secret](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [aws_caller_identity.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity) | data source |
| [aws_iam_policy_document.ecs_execution_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_iam_policy_document.ecs_reader](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_iam_policy_document.locktable_policy_document](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_iam_policy_document.proxy_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |
| [aws_route53_zone.base_zone](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/route53_zone) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_secrets"></a> [additional\_secrets](#input\_additional\_secrets) | Any extra secret key/value pairs to make available to services | `any` | `{}` | no |
| <a name="input_app_ports"></a> [app\_ports](#input\_app\_ports) | What ports do tasks need to be able to reach each other on? | `set(number)` | <pre>[<br>  80,<br>  8080,<br>  8000,<br>  5000,<br>  9000<br>]</pre> | no |
| <a name="input_base_zone"></a> [base\_zone](#input\_base\_zone) | base route53 zone | `string` | n/a | yes |
| <a name="input_batch_envs"></a> [batch\_envs](#input\_batch\_envs) | set of batch envs to create | <pre>map(object({<br>    version         = string,<br>    name            = string,<br>    job_policy_arns = list(string),<br>    min_vcpus       = number,<br>    max_vcpus       = number,<br>    desired_vcpus   = number,<br>    instance_type   = list(string),<br>    init_script     = optional(string),<br>  volume_size = number }))</pre> | `{}` | no |
| <a name="input_cloud-env"></a> [cloud-env](#input\_cloud-env) | n/a | <pre>object({<br>    public_subnets        = list(string)<br>    private_subnets       = list(string)<br>    database_subnets      = list(string)<br>    database_subnet_group = string<br>    vpc_id                = string<br>    vpc_cidr_block        = string<br>  })</pre> | n/a | yes |
| <a name="input_datadog_api_key"></a> [datadog\_api\_key](#input\_datadog\_api\_key) | A datadog api key to enable the datadog agent on the instance | `string` | `""` | no |
| <a name="input_db_engine_version"></a> [db\_engine\_version](#input\_db\_engine\_version) | The Aurora Postgres engine version | `string` | `"14.3"` | no |
| <a name="input_ecr_repos"></a> [ecr\_repos](#input\_ecr\_repos) | Map of ECR repositories to create. These should map exactly to the service names of your docker-compose | <pre>map(object({<br>    name           = string,<br>    read_arns      = optional(list(string), []),<br>    write_arns     = optional(list(string), []),<br>    tag_mutability = optional(bool, true),<br>    scan_on_push   = optional(bool, false),<br>  }))</pre> | `{}` | no |
| <a name="input_extra_proxy_args"></a> [extra\_proxy\_args](#input\_extra\_proxy\_args) | Add to the proxy's default arguments. | `set(string)` | `[]` | no |
| <a name="input_extra_security_groups"></a> [extra\_security\_groups](#input\_extra\_security\_groups) | Security groups that need access to RDS DB's | `list(string)` | `[]` | no |
| <a name="input_github_actions_roles"></a> [github\_actions\_roles](#input\_github\_actions\_roles) | Roles to be used by Github Actions to perform Happy CI. | <pre>set(object({<br>    name = string<br>    arn  = string<br>  }))</pre> | `[]` | no |
| <a name="input_hapi_base_url"></a> [hapi\_base\_url](#input\_hapi\_base\_url) | The base URL for HAPI | `string` | `"https://hapi.hapi.prod.si.czi.technology"` | no |
| <a name="input_instance_type"></a> [instance\_type](#input\_instance\_type) | n/a | `string` | `"m5.large"` | no |
| <a name="input_max_servers"></a> [max\_servers](#input\_max\_servers) | Maximum number of instances for the cluster. Must be at least var.min\_servers + 1. | `number` | `5` | no |
| <a name="input_min_servers"></a> [min\_servers](#input\_min\_servers) | Minimum number of instances for the cluster | `number` | `2` | no |
| <a name="input_name"></a> [name](#input\_name) | n/a | `string` | n/a | yes |
| <a name="input_oauth_bypass_paths"></a> [oauth\_bypass\_paths](#input\_oauth\_bypass\_paths) | Bypass these paths in the oauth proxy | `set(string)` | `[]` | no |
| <a name="input_oauth_dns_prefix"></a> [oauth\_dns\_prefix](#input\_oauth\_dns\_prefix) | DNS prefix for oauth-proxied stacks. Leave this empty if we don't need a prefix! | `string` | `""` | no |
| <a name="input_private_lb_services"></a> [private\_lb\_services](#input\_private\_lb\_services) | Create a private load balancers for a set of services sitting behind Okta/OAuth proxy | `set(string)` | `[]` | no |
| <a name="input_public_lb_services"></a> [public\_lb\_services](#input\_public\_lb\_services) | Create a public-facing ALB for these services | `set(string)` | `[]` | no |
| <a name="input_rds_dbs"></a> [rds\_dbs](#input\_rds\_dbs) | set of DB's to create | `map(object({ name = string, username = string, instance_class = string }))` | `{}` | no |
| <a name="input_regional_wafv2_arn"></a> [regional\_wafv2\_arn](#input\_regional\_wafv2\_arn) | A WAF to protect the happy env if needed | `string` | `null` | no |
| <a name="input_roll_interval_hours"></a> [roll\_interval\_hours](#input\_roll\_interval\_hours) | how often to roll hosts | `number` | `8` | no |
| <a name="input_s3_buckets"></a> [s3\_buckets](#input\_s3\_buckets) | S3 buckets to create | `map(object({ name = string }))` | `{}` | no |
| <a name="input_services"></a> [services](#input\_services) | set of services to prebuild LB's for | `map(object({ idle_timeout = number }))` | `{}` | no |
| <a name="input_ssh_key_name"></a> [ssh\_key\_name](#input\_ssh\_key\_name) | Deprecated | `string` | `"happy_key"` | no |
| <a name="input_ssh_users"></a> [ssh\_users](#input\_ssh\_users) | Okta groups that should have SSH access to ECS nodes | `list(object({ username : string, sudo_enabled : bool }))` | `[]` | no |
| <a name="input_swipe_envs"></a> [swipe\_envs](#input\_swipe\_envs) | set of swipe-batch envs to create | `map(object({ version = string, name = string, job_policy_arns = list(string), min_vcpus = number, max_vcpus = number, spot_desired_vcpus = number, ec2_desired_vcpus = number, instance_type = list(string), ami_id = string }))` | `{}` | no |
| <a name="input_swipe_module_envs"></a> [swipe\_module\_envs](#input\_swipe\_module\_envs) | set of swipe-batch envs to create | <pre>map(object({<br>    name                   = string,<br>    ami_id                 = string,<br>    mock                   = bool,<br>    job_policy_arns        = list(string),<br>    spot_min_vcpus         = number,<br>    spot_max_vcpus         = number,<br>    on_demand_min_vcpus    = number,<br>    on_demand_max_vcpus    = number,<br>    instance_types         = list(string),<br>    extra_env_vars         = map(string),<br>    sqs_queues             = map(map(string)),<br>    workspace_s3_prefix    = string,<br>    wdl_workflow_s3_prefix = string,<br>  }))</pre> | `{}` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | tags to associate with env resources | `map(string)` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_ecs"></a> [ecs](#output\_ecs) | n/a |
| <a name="output_github_ci_role_arns"></a> [github\_ci\_role\_arns](#output\_github\_ci\_role\_arns) | n/a |
| <a name="output_integration_secret"></a> [integration\_secret](#output\_integration\_secret) | Output the same interface we provide to the CLI. |
| <a name="output_public-lbs"></a> [public-lbs](#output\_public-lbs) | n/a |
| <a name="output_security_groups"></a> [security\_groups](#output\_security\_groups) | n/a |
<!-- END -->
