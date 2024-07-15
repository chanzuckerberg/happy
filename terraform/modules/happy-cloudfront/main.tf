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

resource "aws_cloudfront_distribution" "this" {
  enabled     = true
  comment     = "Forward requests from alias to the origins"
  price_class = var.price_class
  aliases     = [var.frontend.domain_name]

  viewer_certificate {
    acm_certificate_arn      = module.cert.arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }

  dynamic "origin" {
    for_each = var.origins
    content {
      domain_name = origin.value.domain_name
      origin_id   = origin.value.domain_name
      custom_origin_config {
        http_port              = "80"
        https_port             = "443"
        origin_protocol_policy = "https-only"
        origin_ssl_protocols   = ["TLSv1.2"]
      }
    }
  }

  dynamic "ordered_cache_behavior" {
    for_each = var.origins
    content {
      viewer_protocol_policy   = "redirect-to-https"
      target_origin_id         = ordered_cache_behavior.value.domain_name
      path_pattern             = ordered_cache_behavior.value.path_pattern
      allowed_methods          = var.allowed_methods
      cached_methods           = var.cache_allowed_methods
      origin_request_policy_id = var.origin_request_policy_id
      cache_policy_id          = var.cache_policy_id

      min_ttl     = var.cache.min_ttl
      default_ttl = var.cache.default_ttl
      max_ttl     = var.cache.max_ttl
      compress    = var.cache.compress
    }
  }

  restrictions {
    geo_restriction {
      locations        = var.geo_restriction_locations
      restriction_type = "whitelist"
    }
  }

  default_cache_behavior {
    viewer_protocol_policy   = "redirect-to-https"
    target_origin_id         = var.origins[length(var.origins) - 1].domain_name
    allowed_methods          = var.allowed_methods
    cached_methods           = var.cache_allowed_methods
    origin_request_policy_id = var.origin_request_policy_id
    cache_policy_id          = var.cache_policy_id

    min_ttl     = var.cache.min_ttl
    default_ttl = var.cache.default_ttl
    max_ttl     = var.cache.max_ttl
    compress    = var.cache.compress
  }

  tags     = var.tags
  provider = aws.useast1
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
  provider = aws.useast1
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
  provider = aws.useast1
}
