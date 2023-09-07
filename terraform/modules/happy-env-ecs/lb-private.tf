# ALB's for PRIVATE (dev) envs

locals {
  private_services = { for s in var.private_lb_services : s => var.services[s] }
  # If we have a regional wafv2 ARN, we keep track of that need in this local variable
  needs_private_waf_attachment = (var.regional_wafv2_arn != null && len(private_services) > 0 ) ? true : false

}

resource "aws_lb" "lb-private" {
  for_each        = local.private_services
  name            = "happy-${var.name}-${each.key}"
  internal        = true
  security_groups = [aws_security_group.happy_env_sg.id]
  subnets         = var.cloud-env.private_subnets
  idle_timeout    = each.value.idle_timeout
  tags            = var.tags
}

resource "aws_lb_listener" "private-lb-listener" {
  for_each          = local.private_services
  load_balancer_arn = aws_lb.lb-private[each.key].arn
  port              = 80
  protocol          = "HTTP"

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

resource "aws_wafv2_web_acl_association" "private" {
  for_each     = (local.needs_private_waf_attachment) ? local.private_services : []
  resource_arn = aws_lb.lb-private[each.key].arn
  web_acl_arn  = var.regional_wafv2_arn
}
