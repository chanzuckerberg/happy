locals {
<<<<<<< HEAD
  base_zone_name = "${var.env}.czi.si.technology"
  hosted_zone_name = "${var.subdomain}.${local.base_zone_name}"
=======
  base_zone_name = "<%= envName %>.czi.si.technology"
  hosted_zone_name = "<%= hostedZoneSubdomain %>.${local.base_zone_name}"
>>>>>>> eb4c92ac1d6268eef2fc5b55eb207728c5dad28d
}

resource "aws_route53_zone" "happy_route53_zone" {
  name = local.hosted_zone_name
  tags = var.tags
}

data "aws_route53_zone" "base" {
  name         = local.base_zone_name
  private_zone = false
  providers = {
    aws = aws.czi-si
  }
}

resource "aws_route53_record" "happy" {
  zone_id = data.aws_route53_zone.base.zone_id
  name    = local.hosted_zone_name
  type    = "NS"
  ttl     = 300
  records = aws_route53_zone.happy_route53_zone.name_servers
}
