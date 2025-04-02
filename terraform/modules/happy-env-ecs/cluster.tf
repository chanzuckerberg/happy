data "aws_region" "current" {}

module "ecs-cluster" {
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/ecs-cluster?ref=ecs-cluster-v2.4.0"

  ecs_cluster_name = "happy-${var.name}"

  region  = data.aws_region.current.name
  project = "happy"
  owner   = "happy"
  env     = var.name

  min_servers                        = var.min_servers
  max_servers                        = var.max_servers
  cluster_asg_rolling_interval_hours = var.roll_interval_hours

  instance_type       = var.instance_type
  vpc_id              = var.cloud-env.vpc_id
  ssh_key_name        = var.ssh_key_name
  subnets             = var.cloud-env.private_subnets
  allowed_cidr_blocks = [var.cloud-env.vpc_cidr_block]
  ssh_users           = var.ssh_users
  docker_storage_size = "214"

  datadog_api_key = var.datadog_api_key
}

resource "aws_cloudwatch_log_group" "ecs" {
  name = "ecs-logs-${var.name}"
}

resource "aws_autoscaling_policy" "scale-up" {
  name                   = "ecs-scale-up-${var.name}"
  scaling_adjustment     = 1
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 120
  autoscaling_group_name = module.ecs-cluster.asg_name
}

resource "aws_autoscaling_policy" "scale-down" {
  name                   = "ecs-scale-down-${var.name}"
  scaling_adjustment     = -1
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 300
  autoscaling_group_name = module.ecs-cluster.asg_name
}

resource "aws_cloudwatch_metric_alarm" "memory-res-high" {
  alarm_name  = "mem-res-high-ecs-${var.name}"
  namespace   = "AWS/ECS"
  metric_name = "MemoryReservation"

  dimensions = {
    ClusterName = module.ecs-cluster.cluster_name
  }

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  period              = "300"
  statistic           = "Average"
  threshold           = "80"

  alarm_actions = [
    aws_autoscaling_policy.scale-up.arn,
  ]
}

resource "aws_cloudwatch_metric_alarm" "memory-res-low" {
  alarm_name  = "mem-res-low-ecs-${var.name}"
  namespace   = "AWS/ECS"
  metric_name = "MemoryReservation"

  dimensions = {
    ClusterName = module.ecs-cluster.cluster_name
  }

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  period              = "300"
  statistic           = "Average"
  threshold           = "60"

  alarm_actions = [
    aws_autoscaling_policy.scale-down.arn,
  ]
}


data "aws_iam_policy_document" "ecs_execution_policy" {
  statement {
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "task_execution_role" {
  name                = "happy-${var.name}-executionrole"
  assume_role_policy  = data.aws_iam_policy_document.ecs_execution_policy.json
  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"]
}
