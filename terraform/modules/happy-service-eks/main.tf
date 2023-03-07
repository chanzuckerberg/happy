data "aws_region" "current" {}

locals {
  tags_string  = join(",", [for key, val in local.routing_tags : "${key}=${val}"])
  service_type = var.routing.service_type == "PRIVATE" ? "ClusterIP" : "NodePort"
  labels = {
    app                            = var.routing.service_name
    "app.kubernetes.io/name"       = var.stack_name
    "app.kubernetes.io/component"  = var.routing.service_name
    "app.kubernetes.io/part-of"    = var.stack_name
    "app.kubernetes.io/managed-by" = "happy"
  }
}

resource "kubernetes_deployment_v1" "deployment" {
  metadata {
    name      = var.routing.service_name
    namespace = var.k8s_namespace
    labels    = local.labels
    annotations = {
      "ad.datadoghq.com/tags" = jsonencode({
        "happy_stack"      = var.stack_name
        "happy_service"    = var.routing.service_name
        "deployment_stage" = var.deployment_stage
        "owner"            = var.tags.owner
        "project"          = var.tags.project
        "env"              = var.tags.env
        "service"          = var.tags.service
        "managedby"        = "happy"
        "happy_compute"    = "eks"
      })
    }
  }

  wait_for_rollout = var.wait_for_steady_state

  spec {
    replicas = var.desired_count

    strategy {
      type = "RollingUpdate"
      rolling_update {
        max_surge       = "25%"
        max_unavailable = "25%"
      }
    }

    progress_deadline_seconds = var.initial_delay_seconds + var.period_seconds

    selector {
      match_labels = {
        app = var.routing.service_name
      }
    }

    template {
      metadata {
        labels = local.labels

        annotations = {
          "ad.datadoghq.com/tags" = jsonencode({
            "happy_stack"      = var.stack_name
            "happy_service"    = var.routing.service_name
            "deployment_stage" = var.deployment_stage
            "owner"            = var.tags.owner
            "project"          = var.tags.project
            "env"              = var.tags.env
            "service"          = var.tags.service
            "managedby"        = "happy"
            "happy_compute"    = "eks"
          })
        }
      }

      spec {
        service_account_name = var.aws_iam_policy_json == "" ? "default" : module.iam_service_account[0].service_account_name

        container {
          image = "${module.ecr.repository_url}:${var.image_tag}"
          name  = var.container_name
          env {
            name  = "DEPLOYMENT_STAGE"
            value = var.deployment_stage
          }
          env {
            name  = "AWS_REGION"
            value = data.aws_region.current.name
          }
          env {
            name  = "AWS_DEFAULT_REGION"
            value = data.aws_region.current.name
          }

          dynamic "env" {
            for_each = var.service_endpoints
            content {
              name  = replace(env.key, "-", "_")
              value = env.value
            }
          }

          dynamic "env" {
            for_each = var.additional_env_vars
            content {
              name  = env.key
              value = env.value
            }
          }

          dynamic "env_from" {
            for_each = toset(var.additional_env_vars_from_config_maps.items)
            content {
              prefix = var.additional_env_vars_from_config_maps.prefix
              config_map_ref {
                name = env_from.value
              }
            }
          }

          dynamic "env_from" {
            for_each = toset(var.additional_env_vars_from_secrets.items)
            content {
              prefix = var.additional_env_vars_from_secrets.prefix
              secret_ref {
                name = env_from.value
              }
            }
          }

          port {
            name           = "http"
            container_port = var.routing.service_port
          }

          resources {
            limits = {
              cpu    = var.cpu
              memory = var.memory
            }
            requests = {
              cpu    = var.cpu
              memory = var.memory
            }
          }

          volume_mount {
            mount_path = "/var/happy"
            name       = "integration-secret"
            read_only  = true
          }

          liveness_probe {
            http_get {
              path = var.health_check_path
              port = var.routing.service_port
            }

            initial_delay_seconds = var.initial_delay_seconds
            period_seconds        = var.period_seconds
          }

          readiness_probe {
            http_get {
              path = var.health_check_path
              port = var.routing.service_port
            }

            initial_delay_seconds = var.initial_delay_seconds
            period_seconds        = var.period_seconds
          }
        }

        volume {
          name = "integration-secret"
          secret {
            secret_name = "integration-secret"
          }
        }
      }
    }
  }
}

resource "kubernetes_service_v1" "service" {
  metadata {
    name      = var.routing.service_name
    namespace = var.k8s_namespace
    labels    = local.labels
  }

  spec {
    selector = {
      app = var.routing.service_name
    }

    port {
      port        = var.routing.service_port
      target_port = var.routing.service_port
    }

    type = local.service_type
  }
}

module "ingress" {
  count = (var.routing.service_type == "EXTERNAL" || var.routing.service_type == "INTERNAL") ? 1 : 0

  source             = "../happy-ingress-eks"
  ingress_name       = var.routing.service_name
  cloud_env          = var.cloud_env
  k8s_namespace      = var.k8s_namespace
  certificate_arn    = var.certificate_arn
  tags_string        = local.tags_string
  routing            = var.routing
  labels             = local.labels
  regional_wafv2_arn = var.regional_wafv2_arn
}

resource "kubernetes_horizontal_pod_autoscaler_v1" "hpa" {
  metadata {
    name      = var.routing.service_name
    namespace = var.k8s_namespace
    labels    = local.labels
  }

  spec {
    max_replicas = var.max_count
    min_replicas = var.desired_count

    target_cpu_utilization_percentage = var.scaling_cpu_threshold_percentage

    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = kubernetes_deployment_v1.deployment.metadata[0].name
    }
  }
}
