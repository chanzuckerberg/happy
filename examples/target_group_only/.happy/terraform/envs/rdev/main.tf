module "stack" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks?ref=heathj/target-group-only"

  image_tag        = var.image_tag
  image_tags       = jsondecode(var.image_tags)
  stack_name       = var.stack_name
  deployment_stage = "rdev"
  stack_prefix     = "/${var.stack_name}"
  k8s_namespace    = var.k8s_namespace

  services = {
    tgonly = {
      name              = "tgonly"
      service_type      = "TARGET_GROUP_ONLY"
      health_check_path = "/mypath/health"
      path              = "/mypath"
      port              = local.port
      alb_name          = aws_lb.this.name
    }
  }
  tasks = {
  }
}

// !!!pretend the below code (ALB, listeners, and security group) was created elsewhere!!!
// !!!such as in a legacy piece of infra that you'd like to!!!
// !!!attach to this happy service!!!
# ---------------------------------------------------------------------------
data "kubernetes_secret" "integration_secret" {
  metadata {
    name      = "integration-secret"
    namespace = var.k8s_namespace
  }
}
locals {
  secret = jsondecode(nonsensitive(data.kubernetes_secret.integration_secret.data.integration_secret))
  port   = 8080
}
// a security group for myalb
resource "aws_security_group" "this" {
  name        = "allow_tls_${var.stack_name}"
  description = "Allow TLS inbound traffic"
  vpc_id      = local.secret["cloud_env"]["vpc_id"]

  ingress {
    description = "TLS from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [local.secret["cloud_env"]["vpc_cidr_block"]]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "allow_tls"
  }
}

resource "aws_lb" "this" {
  name               = "myalb-${var.stack_name}"
  internal           = true
  load_balancer_type = "application"
  subnets            = local.secret["cloud_env"]["private_subnets"]
  security_groups    = [aws_security_group.this.id]
}

resource "aws_lb_listener" "this" {
  load_balancer_arn = aws_lb.this.arn
  port              = local.port
  protocol          = "HTTP"
  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "hello world"
      status_code  = "200"
    }
  }
}
