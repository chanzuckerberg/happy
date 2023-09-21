<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.14 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | ~> 5.14 |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_cert"></a> [cert](#module\_cert) | github.com/chanzuckerberg/cztack//aws-acm-certificate | v0.59.0 |

## Resources

| Name | Type |
|------|------|
| [aws_cloudfront_distribution.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudfront_distribution) | resource |
| [aws_route53_record.alias_ipv4](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_route53_record.alias_ipv6](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [random_pet.this](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cache"></a> [cache](#input\_cache) | The cache settings for the CloudFront distribution. | <pre>object({<br>    min_ttl     = optional(number, 0)<br>    default_ttl = optional(number, 300)<br>    max_ttl     = optional(number, 300)<br>    compress    = optional(bool, true)<br>  })</pre> | `{}` | no |
| <a name="input_frontend"></a> [frontend](#input\_frontend) | The domain name and zone ID the user will see. | <pre>object({<br>    domain_name = string<br>    zone_id     = string<br>  })</pre> | n/a | yes |
| <a name="input_origin"></a> [origin](#input\_origin) | The domain name of the origin. | <pre>object({<br>    domain_name = string<br>  })</pre> | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | Tags to associate with env resources | `map(string)` | n/a | yes |

## Outputs

No outputs.
<!-- END -->