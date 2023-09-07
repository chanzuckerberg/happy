# ALB's for PUBLIC-FACING envs

locals {
  ssl_policy      = "ELBSecurityPolicy-TLS-1-2-Ext-2018-06"
  public_services = { for s in var.public_lb_services : s => var.services[s] }

  # If we have a regional wafv2 ARN, we keep track of that need in this local variable
  needs_public_waf_attachment = var.regional_wafv2_arn != null ? var.public_lb_services : []
}

module "cert-lb" {
  for_each = local.public_services
  source   = "github.com/chanzuckerberg/cztack//aws-acm-certificate?ref=v0.43.1"

  cert_domain_name    = "${each.key}.${local.base_domain}"
  aws_route53_zone_id = var.base_zone
  tags                = var.tags
}

resource "aws_route53_record" "services" {
  for_each = local.public_services
  zone_id  = var.base_zone
  name     = "${each.key}.${local.base_domain}"
  type     = "A"

  alias {
    zone_id                = aws_lb.lb-public[each.key].zone_id
    name                   = aws_lb.lb-public[each.key].dns_name
    evaluate_target_health = true
  }
}

resource "aws_security_group" "public_alb_sg" {
  for_each    = local.public_services
  name        = "happy-${var.name}-${each.key}-public"
  description = "Allow HTTPS inbound traffic"
  vpc_id      = var.cloud-env.vpc_id
  tags        = var.tags

  ingress {
    description = "HTTPS Inbound"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    description = "HTTPS Inbound"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    description = "all egress"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb" "lb-public" {
  for_each        = local.public_services
  name            = "happy-${var.name}-${each.key}"
  internal        = false
  security_groups = [aws_security_group.public_alb_sg[each.key].id]
  subnets         = var.cloud-env.public_subnets
  idle_timeout    = each.value.idle_timeout
  tags            = var.tags
}

resource "aws_lb_listener" "public-https" {
  for_each          = local.public_services
  load_balancer_arn = aws_lb.lb-public[each.key].arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = local.ssl_policy
  certificate_arn   = module.cert-lb[each.key].arn
  tags              = var.tags

  # NOTE: Happy will add listener rules to this listener.
  #       If we ever get to the default rule, we should error out - happy hasn't registered any targets.
  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "Not Found"
      status_code  = "404"
    }
  }
}

resource "aws_lb_listener" "public-http" {
  for_each          = local.public_services
  load_balancer_arn = aws_lb.lb-public[each.key].arn
  port              = 80
  protocol          = "HTTP"
  tags              = var.tags

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_wafv2_web_acl_association" "public" {
  count        = length(local.needs_public_waf_attachment)
  resource_arn = local.needs_public_waf_attachment[count.index]
  web_acl_arn  = var.regional_wafv2_arn
}
