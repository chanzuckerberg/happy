<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | 5.14.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 5.14.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_route53_record.dns_record_0](https://registry.terraform.io/providers/hashicorp/aws/5.14.0/docs/resources/route53_record) | resource |
| [aws_route53_zone.dns_record](https://registry.terraform.io/providers/hashicorp/aws/5.14.0/docs/data-sources/route53_zone) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_alb_dns"></a> [alb\_dns](#input\_alb\_dns) | DNS name for the shared ALB | `string` | n/a | yes |
| <a name="input_canonical_hosted_zone"></a> [canonical\_hosted\_zone](#input\_canonical\_hosted\_zone) | Route53 zone for the shared ALB | `string` | n/a | yes |
| <a name="input_dns_prefix"></a> [dns\_prefix](#input\_dns\_prefix) | Stack-specific prefix for DNS records | `string` | n/a | yes |
| <a name="input_zone"></a> [zone](#input\_zone) | Route53 zone name. Trailing . must be OMITTED! | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END -->