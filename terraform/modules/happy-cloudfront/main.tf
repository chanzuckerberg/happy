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
  origin_id = "happy_cloudfront_${random_pet.this.keepers.origin_domain_name}"
}

resource "aws_cloudfront_distribution" "this" {
  enabled     = true
  comment     = "Forward requests from ${var.frontend.domain_name} to ${var.origin.domain_name}."
  price_class = var.price_class
  aliases     = [var.frontend.domain_name]

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
    allowed_methods          = var.allowed_methods
    cached_methods           = var.allowed_methods
    origin_request_policy_id = var.origin_request_policy_id
    cache_policy_id          = var.cache_policy_id

    min_ttl     = var.cache.min_ttl
    default_ttl = var.cache.default_ttl
    max_ttl     = var.cache.max_ttl
    compress    = var.cache.compress
  }

  restrictions {
    geo_restriction {
      locations        = var.geo_restriction_locations
      restriction_type = "whitelist"
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