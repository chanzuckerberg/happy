locals {
  base_zone_name   = "${var.env}.${var.base_domain}"
  hosted_zone_name = "${var.subdomain}.${local.base_zone_name}"
}

resource "aws_route53_zone" "happy_route53_zone" {
  name = local.hosted_zone_name
  tags = var.tags
}

data "aws_route53_zone" "base" {
  name         = local.base_zone_name
  private_zone = false
  provider     = aws.czi-si
}

resource "aws_route53_record" "happy" {
  zone_id  = data.aws_route53_zone.base.zone_id
  name     = local.hosted_zone_name
  type     = "NS"
  ttl      = 300
  records  = aws_route53_zone.happy_route53_zone.name_servers
  provider = aws.czi-si
}
