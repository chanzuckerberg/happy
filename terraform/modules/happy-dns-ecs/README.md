<!-- START -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_route53_record.dns_record_0](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_route53_zone.dns_record](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/route53_zone) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_alb_dns"></a> [alb\_dns](#input\_alb\_dns) | DNS name for the shared ALB | `string` | n/a | yes |
| <a name="input_app_name"></a> [app\_name](#input\_app\_name) | Please provide the ECS service name | `string` | n/a | yes |
| <a name="input_canonical_hosted_zone"></a> [canonical\_hosted\_zone](#input\_canonical\_hosted\_zone) | Route53 zone for the shared ALB | `string` | n/a | yes |
| <a name="input_custom_stack_name"></a> [custom\_stack\_name](#input\_custom\_stack\_name) | Please provide the stack name | `string` | n/a | yes |
| <a name="input_zone"></a> [zone](#input\_zone) | Route53 zone name. Trailing . must be OMITTED! | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_dns_prefix"></a> [dns\_prefix](#output\_dns\_prefix) | User-facing URL for this service. |
<!-- END -->