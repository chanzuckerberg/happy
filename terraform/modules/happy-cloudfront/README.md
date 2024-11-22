<!-- START -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.14 |
| <a name="requirement_random"></a> [random](#requirement\_random) | ~> 3.5 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws.useast1"></a> [aws.useast1](#provider\_aws.useast1) | ~> 5.14 |

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

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_allowed_methods"></a> [allowed\_methods](#input\_allowed\_methods) | The allowed HTTP methods for the CloudFront distribution. | `set(string)` | <pre>[<br>  "DELETE",<br>  "GET",<br>  "HEAD",<br>  "OPTIONS",<br>  "PATCH",<br>  "POST",<br>  "PUT"<br>]</pre> | no |
| <a name="input_cache"></a> [cache](#input\_cache) | The cache settings for the CloudFront distribution. | <pre>object({<br>    min_ttl     = optional(number, 0)<br>    default_ttl = optional(number, 300)<br>    max_ttl     = optional(number, 300)<br>    compress    = optional(bool, true)<br>  })</pre> | `{}` | no |
| <a name="input_cache_allowed_methods"></a> [cache\_allowed\_methods](#input\_cache\_allowed\_methods) | The allowed cache methods for the CloudFront distribution. | `set(string)` | <pre>[<br>  "GET",<br>  "HEAD"<br>]</pre> | no |
| <a name="input_cache_policy_id"></a> [cache\_policy\_id](#input\_cache\_policy\_id) | The cache policy ID for the CloudFront distribution. | `string` | `"4135ea2d-6df8-44a3-9df3-4b5a84be39ad"` | no |
| <a name="input_frontend"></a> [frontend](#input\_frontend) | The domain name and zone ID the user will see. | <pre>object({<br>    domain_name = string<br>    zone_id     = string<br>  })</pre> | n/a | yes |
| <a name="input_geo_restriction_locations"></a> [geo\_restriction\_locations](#input\_geo\_restriction\_locations) | The countries to whitelist for the CloudFront distribution. | `set(string)` | <pre>[<br>  "US"<br>]</pre> | no |
| <a name="input_origin_request_policy_id"></a> [origin\_request\_policy\_id](#input\_origin\_request\_policy\_id) | The origin request policy ID for the CloudFront distribution. | `string` | `"b689b0a8-53d0-40ab-baf2-68738e2966ac"` | no |
| <a name="input_origins"></a> [origins](#input\_origins) | The domain names and the path used for the origin. | <pre>list(object({<br>    domain_name  = string<br>    path_pattern = string<br>    s3_origin_config = optional(object({ origin_access_identity = string }))<br>  }))</pre> | n/a | yes |
| <a name="input_price_class"></a> [price\_class](#input\_price\_class) | The price class for the CloudFront distribution. | `string` | `"PriceClass_100"` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Tags to associate with env resources | `map(string)` | n/a | yes |

## Outputs

No outputs.
<!-- END -->