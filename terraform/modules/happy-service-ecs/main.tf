data "aws_region" "current" {}

resource "aws_ecs_service" "service" {
  cluster         = var.cluster
  desired_count   = var.desired_count
  task_definition = aws_ecs_task_definition.task_definition.id
  name            = "${var.custom_stack_name}-${var.app_name}"
  launch_type     = var.launch_type
  load_balancer {
    container_name   = var.app_name
    container_port   = var.service_port
    target_group_arn = aws_lb_target_group.target_group.id
  }
  network_configuration {
    security_groups  = var.security_groups
    subnets          = var.subnets
    assign_public_ip = false
  }

  enable_execute_command = true
  wait_for_steady_state  = var.wait_for_steady_state
  tags                   = var.tags
}

locals {
  task_definition = [
    {
      name      = "datadog-agent"
      essential = true
      image     = "${var.datadog_agent.registry}:${var.datadog_agent.tag}"
      cpu       = var.datadog_agent.cpu
      memory    = var.datadog_agent.memory

      environment = concat(
        [
          {
            name  = "DD_API_KEY"
            value = var.datadog_api_key
          },
          {
            name  = "DD_SITE"
            value = "datadoghq.com"
          },
          {
            name  = "DD_SERVICE"
            value = var.app_name
          },
          {
            name  = "DD_ENV"
            value = var.deployment_stage
          },
          {
            name  = "ECS_FARGATE"
            value = "true"
          },
          {
            name  = "DD_APM_ENABLED"
            value = "false"
          },
          {
            name  = "DD_DOGSTATSD_NON_LOCAL_TRAFFIC"
            value = "true"
          },
          {
            name  = "DD_APM_NON_LOCAL_TRAFFIC"
            value = "true"
          },
          {
            name  = "DD_PROCESS_AGENT_ENABLED"
            value = "true"
          },
          {
            name  = "DD_RUNTIME_METRICS_ENABLED"
            value = "true"
          },
          {
            name  = "DD_SYSTEM_PROBE_ENABLED"
            value = "false"
          },
          {
            name  = "DD_GEVENT_PATCH_ALL"
            value = "true"
          },
          {
            name  = "DD_APM_FILTER_TAGS_REJECT"
            value = "http.useragent:ELB-HealthChecker/2.0"
          },
          {
            name  = "DD_TRACE_DEBUG"
            value = "true"
          },
          {
            name  = "DD_LOG_LEVEL"
            value = "debug"
          },
          {
            name  = "DD_EXPVAR_PORT"
            value = "6000"
          },
          {
            name  = "DD_CMD_PORT"
            value = "6001"
          },
          {
            name  = "DD_GUI_PORT"
            value = "6002"
          }
      ])

      "port_mappings" = [
        {
          containerPort = 8126
          hostPort      = 8126
          protocol      = "tcp"
        },
        {
          containerPort = 8125
          hostPort      = 8125
          protocol      = "udp"
      }]

      "logConfiguration" = {
        logDriver = "awslogs"

        options = {
          awslogs-stream-prefix = var.app_name,
          awslogs-group         = aws_cloudwatch_log_group.cloud_watch_datadog_agent_logs_group.id,
          awslogs-region        = data.aws_region.current.name
        }
      }
    },
    {
      name              = var.app_name
      essential         = true
      image             = var.image
      cpu               = var.cpu
      memoryReservation = var.memory
      essential         = true
      portMappings = [{
        containerPort = var.service_port
      }]
      environment = concat(
        [
          {
            name  = "REMOTE_DEV_PREFIX"
            value = var.remote_dev_prefix
          },
          {
            name  = "DEPLOYMENT_STAGE"
            value = var.deployment_stage
          },
          {
            name  = "AWS_REGION"
            value = data.aws_region.current.name
          },
          {
            name  = "AWS_DEFAULT_REGION"
            value = data.aws_region.current.name
          },
          {
            name  = "CHAMBER_SERVICE"
            value = var.chamber_service
          },
        ],
        var.additional_env_vars
      )
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          awslogs-stream-prefix = var.app_name,
          awslogs-group         = aws_cloudwatch_log_group.cloud_watch_logs_group.id,
          awslogs-region        = data.aws_region.current.name
        }
      }
      dockerLabels = {
        # TODO: Add happy_owner and happy_project tags
        "com.datadoghq.ad.tags" = jsonencode([
          "happy_stack:${var.tags.happy_stack_name}",
          "happy_service:${var.tags.happy_service_name}",
          "deployment_stage:${var.deployment_stage}",
          "env:${var.tags.happy_env}",
          "service:${var.tags.happy_service_name}"
        ])
      },
    }
  ]
}

resource "aws_ecs_task_definition" "task_definition" {
  family                   = "${var.stack_resource_prefix}-${var.deployment_stage}-${var.custom_stack_name}-${var.app_name}"
  memory                   = var.memory
  cpu                      = var.cpu
  network_mode             = "awsvpc"
  task_role_arn            = var.task_role.arn
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = var.execution_role
  container_definitions    = jsonencode(local.task_definition)
  tags                     = var.tags
}

resource "aws_cloudwatch_log_group" "cloud_watch_logs_group" {
  retention_in_days = 365
  name              = "/${var.stack_resource_prefix}/${var.deployment_stage}/${var.custom_stack_name}/${var.app_name}"
  tags              = var.tags
}

resource "aws_cloudwatch_log_group" "cloud_watch_datadog_agent_logs_group" {
  retention_in_days = 365
  name              = "/${var.stack_resource_prefix}/${var.deployment_stage}/${var.custom_stack_name}/${var.app_name}/datadog-agent"
  tags              = var.tags
}

resource "aws_lb_target_group" "target_group" {
  vpc_id               = var.vpc
  port                 = var.service_port
  protocol             = "HTTP"
  target_type          = "ip"
  deregistration_delay = 10
  health_check {
    interval            = 15
    path                = var.health_check_path
    protocol            = "HTTP"
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 10
    matcher             = "200-299"
  }
  tags = var.tags
}

resource "aws_lb_listener_rule" "listener_rule" {
  listener_arn = var.listener
  priority     = var.priority
  dynamic "condition" {
    for_each = length(var.host_match) == 0 ? [] : [var.host_match]
    content {
      host_header {
        values = [
          condition.value
        ]
      }
    }
  }
  dynamic "condition" {
    for_each = length(var.host_match) == 0 ? ["/*"] : []
    content {
      path_pattern {
        values = [condition.value]
      }
    }
  }
  action {
    target_group_arn = aws_lb_target_group.target_group.id
    type             = "forward"
  }
  tags = var.tags
}
