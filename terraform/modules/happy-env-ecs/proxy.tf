locals {
  base_domain = data.aws_route53_zone.base_zone.name
  env_domain  = length(var.oauth_dns_prefix) == 0 ? local.base_domain : "${var.oauth_dns_prefix}.${local.base_domain}"
}

data "aws_route53_zone" "base_zone" {
  zone_id      = var.base_zone
  private_zone = true
}

resource "aws_route53_zone" "happy" {
  count = length(var.oauth_dns_prefix) == 0 ? 0 : 1
  name  = local.env_domain
  tags  = var.tags
}

resource "aws_route53_record" "happy-NS" {
  count   = length(var.oauth_dns_prefix) == 0 ? 0 : 1
  zone_id = var.base_zone
  name    = local.env_domain
  type    = "NS"
  ttl     = 300
  records = aws_route53_zone.happy[count.index].name_servers
}

data "aws_iam_policy_document" "proxy_role" {
  statement {
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "proxy_role" {
  count              = length(var.private_lb_services) > 0 ? 1 : 0
  name               = "happy-${var.name}-proxy"
  assume_role_policy = data.aws_iam_policy_document.proxy_role.json
}

module "ecs-multi-domain-oauth-proxy" {
  count  = length(var.private_lb_services) > 0 ? 1 : 0
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/ecs-multi-domain-oauth-proxy?ref=ecs-multi-domain-oauth-proxy-v2.1.0"
  cloud-env = {
    public_subnets  = var.cloud-env.public_subnets,
    private_subnets = var.cloud-env.private_subnets,
    vpc_id          = var.cloud-env.vpc_id
  }
  ecs                         = module.ecs-cluster.ecs
  route53_base_zone_id        = length(aws_route53_zone.happy) == 0 ? data.aws_route53_zone.base_zone.zone_id : aws_route53_zone.happy[count.index].zone_id
  tags                        = var.tags
  target_port                 = "80"
  task_role_arn               = aws_iam_role.proxy_role[count.index].arn
  bypass_paths                = var.oauth_bypass_paths
  extra_proxy_args            = var.extra_proxy_args
  oauth2_proxy_registry_image = var.oauth2_proxy_registry_image
  oauth2_proxy_image_version  = var.oauth2_proxy_image_version
}
