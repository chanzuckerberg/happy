https://docs.google.com/drawings/d/1AsJts2qCmw7685A6WZPDb5ApkXyuPRc27Lg3zzWuPaA/edit
<!-- bump3  -->
<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.14 |
| <a name="requirement_datadog"></a> [datadog](#requirement\_datadog) | >= 3.20.0 |
| <a name="requirement_kubernetes"></a> [kubernetes](#requirement\_kubernetes) | >= 2.16 |
| <a name="requirement_random"></a> [random](#requirement\_random) | >= 3.4 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 5.14 |
| <a name="provider_datadog"></a> [datadog](#provider\_datadog) | >= 3.20.0 |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | >= 2.16 |
| <a name="provider_random"></a> [random](#provider\_random) | >= 3.4 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_cert"></a> [cert](#module\_cert) | github.com/chanzuckerberg/cztack//aws-acm-certificate | v0.43.1 |
| <a name="module_dbs"></a> [dbs](#module\_dbs) | github.com/chanzuckerberg/cztack//aws-aurora-postgres | v0.49.0 |
| <a name="module_ecrs"></a> [ecrs](#module\_ecrs) | git@github.com:chanzuckerberg/cztack//aws-ecr-repo | v0.59.0 |
| <a name="module_happy_github_ci_role"></a> [happy\_github\_ci\_role](#module\_happy\_github\_ci\_role) | ../happy-github-ci-role | n/a |
| <a name="module_happy_okta_app"></a> [happy\_okta\_app](#module\_happy\_okta\_app) | ../happy-tfe-okta-app | n/a |
| <a name="module_happy_service_account"></a> [happy\_service\_account](#module\_happy\_service\_account) | ../happy-tfe-okta-service-account | n/a |
| <a name="module_ops-genie"></a> [ops-genie](#module\_ops-genie) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/ops-genie-service | main |
| <a name="module_s3_buckets"></a> [s3\_buckets](#module\_s3\_buckets) | github.com/chanzuckerberg/cztack//aws-s3-private-bucket | v0.56.2 |

## Resources

| Name | Type |
|------|------|
| [aws_dynamodb_table.locks](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table) | resource |
| [aws_iam_policy.locktable_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [datadog_dashboard_json.environment_dashboard](https://registry.terraform.io/providers/datadog/datadog/latest/docs/resources/dashboard_json) | resource |
| [kubernetes_namespace.happy](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/namespace) | resource |
| [kubernetes_secret.happy_env_secret](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/secret) | resource |
| [random_password.db_secret](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [aws_caller_identity.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity) | data source |
| [aws_iam_policy_document.locktable_policy_document](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_route53_zone.base_zone](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/route53_zone) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_secrets"></a> [additional\_secrets](#input\_additional\_secrets) | Any extra secret key/value pairs to make available to services | `any` | `{}` | no |
| <a name="input_base_zone_id"></a> [base\_zone\_id](#input\_base\_zone\_id) | The base zone all happy stacks and infrastructure will build on top of | `string` | n/a | yes |
| <a name="input_cloud-env"></a> [cloud-env](#input\_cloud-env) | n/a | <pre>object({<br>    public_subnets        = list(string)<br>    private_subnets       = list(string)<br>    database_subnets      = list(string)<br>    database_subnet_group = string<br>    vpc_id                = string<br>    vpc_cidr_block        = string<br>  })</pre> | n/a | yes |
| <a name="input_default_db_engine_version"></a> [default\_db\_engine\_version](#input\_default\_db\_engine\_version) | The default Aurora Postgres engine version if one is not specified in rds\_dbs | `string` | `"14.3"` | no |
| <a name="input_ecr_repos"></a> [ecr\_repos](#input\_ecr\_repos) | Map of ECR repositories to create. These should map exactly to the service names of your docker-compose | <pre>map(object({<br>    name           = string,<br>    read_arns      = optional(list(string), []),<br>    write_arns     = optional(list(string), []),<br>    tag_mutability = optional(bool, true),<br>    scan_on_push   = optional(bool, false),<br>  }))</pre> | `{}` | no |
| <a name="input_eks-cluster"></a> [eks-cluster](#input\_eks-cluster) | eks-cluster module output | <pre>object({<br>    cluster_id : string,<br>    cluster_arn : string,<br>    cluster_endpoint : string,<br>    cluster_ca : string,<br>    cluster_oidc_issuer_url : string,<br>    cluster_version : string,<br>    worker_iam_role_name : string,<br>    worker_security_group : string,<br>    oidc_provider_arn : string,<br>  })</pre> | n/a | yes |
| <a name="input_github_actions_roles"></a> [github\_actions\_roles](#input\_github\_actions\_roles) | Roles to be used by Github Actions to perform Happy CI. | <pre>set(object({<br>    name = string<br>    arn  = string<br>  }))</pre> | `[]` | no |
| <a name="input_hapi_base_url"></a> [hapi\_base\_url](#input\_hapi\_base\_url) | The base URL for HAPI | `string` | `"https://hapi.hapi.prod.si.czi.technology"` | no |
| <a name="input_oidc_config"></a> [oidc\_config](#input\_oidc\_config) | OIDC configuration for the happy stacks in this environment. | <pre>object({<br>    login_uri                  = optional(string, ""),<br>    grant_types                = optional(set(string), ["authorization_code", "refresh_token"])<br>    redirect_uris              = optional(set(string), []),<br>    teams                      = optional(set(string), []),<br>    app_type                   = optional(string, "web"),<br>    token_endpoint_auth_method = optional(string, "client_secret_basic"),<br>  })</pre> | `{}` | no |
| <a name="input_okta_teams"></a> [okta\_teams](#input\_okta\_teams) | The set of Okta teams to give access to the Okta app | `set(string)` | `null` | no |
| <a name="input_ops_genie_owner_team"></a> [ops\_genie\_owner\_team](#input\_ops\_genie\_owner\_team) | The name of the Opsgenie team that will own the alerts for this happy environment | `string` | `"Core Infra Eng"` | no |
| <a name="input_rds_dbs"></a> [rds\_dbs](#input\_rds\_dbs) | Map of DB's to create for your happy applications. If an engine\_version is not provided, the default\_db\_engine\_version is used | <pre>map(object({<br>    engine_version : string,<br>    instance_class : string,<br>    username : string,<br>    name : string,<br>    rds_cluster_parameters : optional(list(<br>      map(any)), []<br>    ),<br>  }))</pre> | `{}` | no |
| <a name="input_s3_buckets"></a> [s3\_buckets](#input\_s3\_buckets) | Map of S3 buckets to create for your happy applications | <pre>map(object(<br>    {<br>      name   = string<br>      policy = optional(string, "")<br>  }))</pre> | `{}` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Standard tags. Typically generated by fogg | <pre>object({<br>    env : string,<br>    owner : string,<br>    project : string,<br>    service : string,<br>    managedBy : string,<br>  })</pre> | n/a | yes |
| <a name="input_waf_arn"></a> [waf\_arn](#input\_waf\_arn) | A regional WAF ARN to attach to the happy ingress. | `string` | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_dashboard"></a> [dashboard](#output\_dashboard) | n/a |
| <a name="output_databases"></a> [databases](#output\_databases) | n/a |
| <a name="output_integration_secret"></a> [integration\_secret](#output\_integration\_secret) | n/a |
| <a name="output_namespace"></a> [namespace](#output\_namespace) | n/a |
<!-- END -->
//

