
data "aws_route53_zone" "base_zone" {
  zone_id      = var.base_zone_id
  private_zone = true
}

module "cert" {
  source = "github.com/chanzuckerberg/cztack//aws-acm-certificate?ref=v0.43.1"

  cert_domain_name = data.aws_route53_zone.base_zone.name
  cert_subject_alternative_names = {
    "*.${data.aws_route53_zone.base_zone.name}" = var.base_zone_id
  }

  aws_route53_zone_id = var.base_zone_id
  tags                = merge(var.tags, { "managedBy" : "terraform" })
}
