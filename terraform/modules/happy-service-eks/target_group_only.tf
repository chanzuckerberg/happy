resource "random_pet" "this" {
  keepers = {
    target_group_name = var.routing.service_name
  }
}

locals {
  # only hyphens and a max of 32 characters
  target_group_name = replace(substr(random_pet.this.id, 0, 32), "_", "-")
}

data "aws_lb" "this" {
  count = var.routing.service_type == "TARGET_GROUP_ONLY" ? 1 : 0

  name = var.routing.alb.name
}

data "aws_lb_listener" "this" {
  count = var.routing.service_type == "TARGET_GROUP_ONLY" ? 1 : 0

  load_balancer_arn = data.aws_lb.this[0].arn
  port              = var.routing.alb.listener_port
}

resource "aws_lb_target_group" "this" {
  count = var.routing.service_type == "TARGET_GROUP_ONLY" ? 1 : 0

  name     = local.target_group_name
  port     = var.routing.service_port
  protocol = "HTTP"
  vpc_id   = var.cloud_env.vpc_id
  health_check {
    path = var.health_check_path
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lb_listener_rule" "this" {
  count = var.routing.service_type == "TARGET_GROUP_ONLY" ? 1 : 0

  listener_arn = data.aws_lb_listener.this[0].arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this[0].arn
  }

  condition {
    path_pattern {
      values = [var.routing.path]
    }
  }
}

resource "kubernetes_manifest" "this" {
  count = var.routing.service_type == "TARGET_GROUP_ONLY" ? 1 : 0

  manifest = {
    apiVersion = "elbv2.k8s.aws/v1beta1"
    kind       = "TargetGroupBinding"

    metadata = {
      name      = local.target_group_name
      namespace = var.k8s_namespace

    }

    spec = {
      serviceRef = {
        name = var.routing.service_name
        port = var.routing.service_port
      }
      targetGroupARN = aws_lb_target_group.this[0].arn
      networking = {
        ingress = [{
          from = [for sg_id in data.aws_lb.this[0].security_groups : { securityGroup = { groupID = sg_id } }]
          ports = [{
            protocol = "TCP"
          }]
        }]
      }
    }
  }
}
