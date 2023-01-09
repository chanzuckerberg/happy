https://docs.google.com/drawings/d/1AsJts2qCmw7685A6WZPDb5ApkXyuPRc27Lg3zzWuPaA/edit

<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.45 |
| <a name="requirement_random"></a> [random](#requirement\_random) | >= 3.4 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 4.45 |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | n/a |
| <a name="provider_random"></a> [random](#provider\_random) | >= 3.4 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_cert"></a> [cert](#module\_cert) | github.com/chanzuckerberg/cztack//aws-acm-certificate | v0.43.1 |
| <a name="module_cert_oauth"></a> [cert\_oauth](#module\_cert\_oauth) | github.com/chanzuckerberg/cztack//aws-acm-certificate | v0.43.1 |
| <a name="module_dbs"></a> [dbs](#module\_dbs) | github.com/chanzuckerberg/cztack//aws-aurora-postgres | v0.49.0 |
| <a name="module_ecrs"></a> [ecrs](#module\_ecrs) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/ecr-repository | main |
| <a name="module_happy_github_ci_role"></a> [happy\_github\_ci\_role](#module\_happy\_github\_ci\_role) | ../happy-github-ci-role | n/a |
| <a name="module_proxy"></a> [proxy](#module\_proxy) | git@github.com:chanzuckerberg/shared-infra//terraform/modules/eks-multi-domain-oauth-proxy | eks-multi-domain-oauth-proxy-v1.3.0 |
| <a name="module_s3_buckets"></a> [s3\_buckets](#module\_s3\_buckets) | github.com/chanzuckerberg/cztack//aws-s3-private-bucket | v0.43.1 |

## Resources

| Name | Type |
|------|------|
| [aws_route53_record.happy_prefixed](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_route53_zone.happy_prefixed](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_zone) | resource |
| [kubernetes_namespace.happy](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/namespace) | resource |
| [kubernetes_secret.happy_env_secret](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/secret) | resource |
| [random_password.db_secret](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [aws_route53_zone.base_zone](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/route53_zone) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_secrets"></a> [additional\_secrets](#input\_additional\_secrets) | Any extra secret key/value pairs to make available to services | `any` | `{}` | no |
| <a name="input_authorized_github_repos"></a> [authorized\_github\_repos](#input\_authorized\_github\_repos) | Map of (arbitrary) identifier to Github repo and happy app name that are authorized to assume the created CI role | `map(object({ repo_name : string, app_name : string }))` | `{}` | no |
| <a name="input_base_zone_id"></a> [base\_zone\_id](#input\_base\_zone\_id) | The base zone all happy stacks and infrastructure will build on top of | `string` | n/a | yes |
| <a name="input_cloud-env"></a> [cloud-env](#input\_cloud-env) | n/a | <pre>object({<br>    public_subnets        = list(string)<br>    private_subnets       = list(string)<br>    database_subnets      = list(string)<br>    database_subnet_group = string<br>    vpc_id                = string<br>    vpc_cidr_block        = string<br>  })</pre> | n/a | yes |
| <a name="input_default_db_engine_version"></a> [default\_db\_engine\_version](#input\_default\_db\_engine\_version) | The default Aurora Postgres engine version if one is not specified in rds\_dbs | `string` | `"14.3"` | no |
| <a name="input_ecr_repos"></a> [ecr\_repos](#input\_ecr\_repos) | Map of ECR repositories to create. These should map exactly to the service names of your docker-compose | <pre>map(object({<br>    name       = string,<br>    read_arns  = list(string),<br>    write_arns = list(string),<br>  }))</pre> | `{}` | no |
| <a name="input_eks-cluster"></a> [eks-cluster](#input\_eks-cluster) | eks-cluster module output | <pre>object({<br>    cluster_id : string,<br>    cluster_arn : string,<br>    cluster_endpoint : string,<br>    cluster_ca : string,<br>    cluster_oidc_issuer_url : string,<br>    cluster_security_group : string,<br>    cluster_iam_role_name : string,<br>    cluster_version : string,<br>    worker_iam_role_name : string,<br>    kubeconfig : string,<br>    worker_security_group : string,<br>    oidc_provider_arn : string,<br>  })</pre> | n/a | yes |
| <a name="input_k8s-core"></a> [k8s-core](#input\_k8s-core) | K8s core. Typically the outputs of the remote state for the corresponding k8s-core component. | <pre>object({<br>    default_namespace : string,<br>    aws_ssm_iam_role_name : string,<br>  })</pre> | n/a | yes |
| <a name="input_oauth_bypass_paths"></a> [oauth\_bypass\_paths](#input\_oauth\_bypass\_paths) | Bypass these paths in the oauth proxy | `list(string)` | `[]` | no |
| <a name="input_oauth_dns_prefix"></a> [oauth\_dns\_prefix](#input\_oauth\_dns\_prefix) | DNS prefix for oauth-proxied stacks. Leave this empty if we don't need a prefix! | `string` | `""` | no |
| <a name="input_oidc_issuer_host"></a> [oidc\_issuer\_host](#input\_oidc\_issuer\_host) | The OIDC issuer host for the OIDC provider to use for happy authentication | `string` | `"czi-prod.okta.com"` | no |
| <a name="input_rds_dbs"></a> [rds\_dbs](#input\_rds\_dbs) | Map of DB's to create for your happy applications. If an engine\_version is not provided, the default\_db\_engine\_version is used | <pre>map(object({<br>    name           = string,<br>    username       = string,<br>    instance_class = string,<br>    engine_version = string,<br>  }))</pre> | `{}` | no |
| <a name="input_s3_buckets"></a> [s3\_buckets](#input\_s3\_buckets) | Map of S3 buckets to create for your happy applications | `map(object({ name = string }))` | `{}` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Standard tags. Typically generated by fogg | <pre>object({<br>    env : string,<br>    owner : string,<br>    project : string,<br>    service : string,<br>    managedBy : string,<br>  })</pre> | n/a | yes |

## Outputs

No outputs.
<!-- END -->
