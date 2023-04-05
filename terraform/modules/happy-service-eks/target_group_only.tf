resource "random_pet" "this" {
  keepers = {
    target_group_name = var.routing.service_name
  }
}

resource "aws_lb_target_group" "this" {
  count = var.routing.service_type == "TARGET_GROUP_ONLY" ? 1 : 0

  name     = random_pet.this.keepers.target_group_name
  port     = 80
  protocol = "HTTP"
  vpc_id   = var.cloud_env.vpc_id
  health_check {
    path = var.health_check_path
  }
}

resource "kubernetes_manifest" "this" {
  for_each = toset(aws_lb_target_group.this)

  manifest = {
    apiVersion = "elbv2.k8s.aws/v1beta1"
    kind       = "TargetGroupBinding"

    metadata = {
      name      = random_pet.this.keepers.target_group_name
      namespace = var.k8s_namespace

    }

    spec = {
      serviceRef = {
        name = var.routing.service_name
        port = var.routing.service_port
      }
      targetGroupARN = each.value.arn
    }
  }
}
