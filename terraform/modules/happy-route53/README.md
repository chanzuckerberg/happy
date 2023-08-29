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
| <a name="provider_aws.czi-si"></a> [aws.czi-si](#provider\_aws.czi-si) | >= 5.14 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_route53_record.happy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_route53_zone.happy_route53_zone](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_zone) | resource |
| [aws_route53_zone.base](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/route53_zone) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_env"></a> [env](#input\_env) | The deployment environment for the app | `string` | n/a | yes |
| <a name="input_subdomain"></a> [subdomain](#input\_subdomain) | The hosted zone subdomain to create | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | Tags to associate with env resources | `map(string)` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_happy_route53_zone_name"></a> [happy\_route53\_zone\_name](#output\_happy\_route53\_zone\_name) | n/a |
| <a name="output_happy_route53_zone_zone_id"></a> [happy\_route53\_zone\_zone\_id](#output\_happy\_route53\_zone\_zone\_id) | n/a |
<!-- END -->