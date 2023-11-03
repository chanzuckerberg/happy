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
  name = var.routing.alb.name
}

data "aws_lb_listener" "this" {
  load_balancer_arn = data.aws_lb.this.arn
  port              = var.routing.alb.listener_port
}

resource "aws_lb_target_group" "this" {
  name     = local.target_group_name
  port     = var.routing.service_port
  protocol = "HTTP"
  vpc_id   = var.cloud_env.vpc_id
  health_check {
    path = var.health_check_path
  }
}

resource "aws_lb_listener_rule" "this" {
  listener_arn = data.aws_lb_listener.this.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }

  condition {
    path_pattern {
      values = [var.routing.path]
    }
  }
}