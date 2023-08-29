The EKS specific permissions for happy to work with a Github Action role
<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.14 |
| <a name="requirement_kubernetes"></a> [kubernetes](#requirement\_kubernetes) | >= 2.19 |
| <a name="requirement_random"></a> [random](#requirement\_random) | >= 3.5.1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 5.14 |
| <a name="provider_random"></a> [random](#provider\_random) | >= 3.5.1 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_iam_role_policy.eks](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [random_pet.this](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [aws_iam_policy_document.eks](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_eks_cluster_arn"></a> [eks\_cluster\_arn](#input\_eks\_cluster\_arn) | The ARN of the EKS cluster that the role should have permissions to | `string` | n/a | yes |
| <a name="input_gh_actions_role_name"></a> [gh\_actions\_role\_name](#input\_gh\_actions\_role\_name) | The name of the role that was created for the Github Action. | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END -->