module "cert" {
  source = "github.com/chanzuckerberg/cztack//aws-acm-certificate?ref=v0.59.0"

  cert_domain_name    = var.frontend.domain_name
  aws_route53_zone_id = var.frontend.zone_id
  tags                = var.tags

  # NOTE: certificates need to be in us-east-1 for cloudfront
  providers = {
    aws = aws.useast1
  }
}

resource "random_pet" "this" {
  keepers = {
    origin_domain_name = var.origin.domain_name
  }
}

locals {
  origin_id                                = "happy_cloudfront_${random_pet.this.keepers.origin_domain_name}"
  managed_caching_disabled_policy_id       = "4135ea2d-6df8-44a3-9df3-4b5a84be39ad"
  managed_all_viewer_except_host_policy_id = "b689b0a8-53d0-40ab-baf2-68738e2966ac"
}

resource "aws_cloudfront_distribution" "this" {
  enabled = true
  comment = "Forward requests from ${var.frontend.domain_name} to ${var.origin.domain_name}."

  aliases = [var.frontend.domain_name]
  viewer_certificate {
    acm_certificate_arn      = module.cert.arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }

  origin {
    domain_name = var.origin.domain_name
    origin_id   = local.origin_id
    custom_origin_config {
      https_port             = "443"
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  default_cache_behavior {
    viewer_protocol_policy   = "https"
    target_origin_id         = local.origin_id
    allowed_methods          = ["GET", "HEAD"]
    cached_methods           = ["GET", "HEAD"]
    origin_request_policy_id = local.managed_all_viewer_except_host_policy_id
    cache_policy_id          = local.managed_caching_disabled_policy

    min_ttl     = var.cache.min_ttl
    default_ttl = var.cache.default_ttl
    max_ttl     = var.cache.max_ttl
    compress    = var.cache.compress
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  tags = var.tags
}

resource "aws_route53_record" "alias_ipv4" {
  zone_id = var.frontend.zone_id
  name    = var.frontend.domain_name
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.this.domain_name
    zone_id                = aws_cloudfront_distribution.this.hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "alias_ipv6" {
  zone_id = var.frontend.zone_id
  name    = var.frontend.domain_name
  type    = "AAAA"

  alias {
    name                   = aws_cloudfront_distribution.this.domain_name
    zone_id                = aws_cloudfront_distribution.this.hosted_zone_id
    evaluate_target_health = false
  }
}