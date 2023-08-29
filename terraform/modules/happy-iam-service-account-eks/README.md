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

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_iam_policy.policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_role.role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy_attachment.attach](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [kubernetes_service_account.service_account](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service_account) | resource |
| [aws_iam_policy_document.assume-role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_iam_policies_json"></a> [aws\_iam\_policies\_json](#input\_aws\_iam\_policies\_json) | The additional AWS IAM policies to give to the pod. Backward compatibility with aws\_iam\_policy\_json | `list(string)` | `[]` | no |
| <a name="input_aws_iam_policy_json"></a> [aws\_iam\_policy\_json](#input\_aws\_iam\_policy\_json) | The AWS IAM policy to give to the pod. | `string` | n/a | yes |
| <a name="input_eks_cluster"></a> [eks\_cluster](#input\_eks\_cluster) | eks-cluster module output | <pre>object({<br>    cluster_id : string,<br>    cluster_arn : string,<br>    cluster_endpoint : string,<br>    cluster_ca : string,<br>    cluster_oidc_issuer_url : string,<br>    cluster_version : string,<br>    worker_iam_role_name : string,<br>    worker_security_group : string,<br>    oidc_provider_arn : string,<br>  })</pre> | n/a | yes |
| <a name="input_iam_path"></a> [iam\_path](#input\_iam\_path) | IAM path for the role. | `string` | `""` | no |
| <a name="input_k8s_namespace"></a> [k8s\_namespace](#input\_k8s\_namespace) | Kubernetes namespace that the service account is in | `string` | n/a | yes |
| <a name="input_max_session_duration"></a> [max\_session\_duration](#input\_max\_session\_duration) | Maximum CLI/API session duration in seconds between 3600 and 43200 | `number` | `3600` | no |
| <a name="input_role_permissions_boundary_arn"></a> [role\_permissions\_boundary\_arn](#input\_role\_permissions\_boundary\_arn) | Permissions boundary ARN to use for IAM role | `string` | `""` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | The happy conventional tags. | <pre>object({<br>    happy_env : string,<br>    happy_stack_name : string,<br>    happy_service_name : string,<br>    happy_region : string,<br>    happy_image_tag : string,<br>    happy_service_type : string,<br>    happy_last_applied : string,<br>  })</pre> | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_iam_role"></a> [iam\_role](#output\_iam\_role) | n/a |
| <a name="output_iam_role_arn"></a> [iam\_role\_arn](#output\_iam\_role\_arn) | n/a |
| <a name="output_iam_role_name_with_path"></a> [iam\_role\_name\_with\_path](#output\_iam\_role\_name\_with\_path) | n/a |
| <a name="output_service_account_name"></a> [service\_account\_name](#output\_service\_account\_name) | n/a |
<!-- END -->
