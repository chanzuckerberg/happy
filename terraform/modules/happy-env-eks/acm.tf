module "cert" {
  source = "github.com/chanzuckerberg/cztack//aws-acm-certificate?ref=v0.43.1"

  cert_domain_name = local.base_domain
  cert_subject_alternative_names = {
    "*.${local.base_domain}" = var.base_zone_id
  }

  aws_route53_zone_id = var.base_zone_id
  tags                = merge(var.tags, { "managedBy" : "terraform" })
}

module "cert_oauth" {
  source = "github.com/chanzuckerberg/cztack//aws-acm-certificate?ref=v0.43.1"
  count  = length(var.oauth_dns_prefix) == 0 ? 0 : 1

  cert_domain_name = local.env_domain
  cert_subject_alternative_names = {
    "*.${local.env_domain}" = var.base_zone_id
  }

  aws_route53_zone_id = var.base_zone_id
  tags                = merge(var.tags, { "managedBy" : "terraform" })
}
