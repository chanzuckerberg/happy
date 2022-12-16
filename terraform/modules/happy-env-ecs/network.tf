locals {
  port_strings = toset([for num in var.app_ports : tostring(num)])
  lb_allow_ports = coalesce(merge(flatten([
    for service_name, _ in local.public_services : {
      for port in local.port_strings : "${service_name}-${port}" => {
        service_name = service_name
        port         = port
      }
    }
  ])...), {})
}
resource "aws_security_group" "happy_env_sg" {
  name   = "happy_${var.name}-sg"
  vpc_id = var.cloud-env.vpc_id
  tags   = var.tags
}

resource "aws_security_group_rule" "tasks" {
  for_each          = local.port_strings
  type              = "ingress"
  from_port         = tonumber(each.value)
  to_port           = tonumber(each.value)
  protocol          = "tcp"
  self              = true
  security_group_id = aws_security_group.happy_env_sg.id
}
resource "aws_security_group_rule" "oauth" {
  count                    = length(var.private_lb_services) > 0 ? 1 : 0
  type                     = "ingress"
  from_port                = 80
  to_port                  = 80
  protocol                 = "tcp"
  source_security_group_id = module.ecs-multi-domain-oauth-proxy[count.index].nginx_container_security_groups
  security_group_id        = aws_security_group.happy_env_sg.id
}
resource "aws_security_group_rule" "alb" {
  for_each                 = local.lb_allow_ports
  type                     = "ingress"
  from_port                = each.value.port
  to_port                  = each.value.port
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.public_alb_sg[each.value["service_name"]].id
  security_group_id        = aws_security_group.happy_env_sg.id
}
resource "aws_security_group_rule" "egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.happy_env_sg.id
}

