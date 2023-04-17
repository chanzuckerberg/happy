The EKS specific permissions for happy to work with a Github Action role
<!-- START -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_iam_role_policy.eks](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [kubernetes_config_map_v1_data.aws-auth](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/config_map_v1_data) | resource |
| [aws_iam_policy_document.eks](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [kubernetes_config_map.aws-auth](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/data-sources/config_map) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_eks_cluster_arn"></a> [eks\_cluster\_arn](#input\_eks\_cluster\_arn) | The ARN of the EKS cluster that the role should have permissions to | `string` | n/a | yes |
| <a name="input_gh_actions_role"></a> [gh\_actions\_role](#input\_gh\_actions\_role) | The role that was created for the Github Action. | <pre>object({<br>    role : object({<br>      arn : string,<br>      name : string,<br>    }),<br>  })</pre> | n/a | yes |

## Outputs

No outputs.
<!-- END -->