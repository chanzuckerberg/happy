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

| Name | Source | Version |
|------|--------|---------|
| <a name="module_tfe-si-happy-roles"></a> [tfe-si-happy-roles](#module\_tfe-si-happy-roles) | github.com/chanzuckerberg/cztack//aws-iam-group-assume-role | v0.43.1 |

## Resources

| Name | Type |
|------|------|
| [aws_iam_user.tfe-happy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_user) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_accounts_can_assume"></a> [aws\_accounts\_can\_assume](#input\_aws\_accounts\_can\_assume) | The set of AWS account names the TFE user should be allowed to assume into | `set(string)` | n/a | yes |
| <a name="input_happy_app_name"></a> [happy\_app\_name](#input\_happy\_app\_name) | The name of the happy application | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END -->