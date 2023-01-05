locals {
  base_domain = data.aws_route53_zone.base_zone.name
  env_domain  = length(var.oauth_dns_prefix) == 0 ? local.base_domain : "${var.oauth_dns_prefix}.${local.base_domain}"
}

data "aws_route53_zone" "base_zone" {
  zone_id      = var.base_zone_id
  private_zone = true
}

resource "aws_route53_zone" "happy_prefixed" {
  count = length(var.oauth_dns_prefix) == 0 ? 0 : 1
  name  = local.env_domain
  tags  = var.tags
}

resource "aws_route53_record" "happy_prefixed" {
  count   = length(var.oauth_dns_prefix) == 0 ? 0 : 1
  zone_id = var.base_zone_id
  name    = local.env_domain
  type    = "NS"
  ttl     = 300
  records = aws_route53_zone.happy_prefixed[count.index].name_servers
}

module "proxy" {
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/eks-multi-domain-oauth-proxy?ref=v0.237.0"

  tags      = var.tags
  eks       = var.eks-cluster
  k8s-core  = var.k8s-core
  namespace = kubernetes_namespace.happy.metadata[0].name
  cloud-env = {
    public_subnets  = var.cloud-env.public_subnets,
    private_subnets = var.cloud-env.private_subnets,
    vpc_id          = var.cloud-env.vpc_id
  }
  route53_base_zone_id = length(aws_route53_zone.happy_prefixed) == 0 ? data.aws_route53_zone.base_zone.zone_id : aws_route53_zone.happy_prefixed[0].zone_id
  oidc_issuer_host     = var.oidc_issuer_host
}
