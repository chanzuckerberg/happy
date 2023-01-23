data "aws_region" "current" {}

locals {
  tags_string = join(",", [for key, val in local.tags : "${key}=${val}"])
}

resource "kubernetes_deployment" "deployment" {
  metadata {
    name      = var.service_name
    namespace = var.k8s_namespace
    labels = {
      app = var.service_name
    }
  }

  wait_for_rollout = var.wait_for_steady_state

  spec {
    replicas = var.desired_count

    selector {
      match_labels = {
        app = var.service_name
      }
    }

    template {
      metadata {
        labels = {
          app = var.service_name
        }
      }

      spec {
        service_account_name = var.aws_iam_policy_json == "" ? "default" : module.iam_service_account[0].service_account_name

        container {
          image = var.image
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

          port {
            name           = "http"
            container_port = var.service_port
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
              port = var.service_port
            }

            initial_delay_seconds = var.initial_delay_seconds
            period_seconds        = var.period_seconds
          }

          readiness_probe {
            http_get {
              path = var.health_check_path
              port = var.service_port
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

resource "kubernetes_service" "service" {
  metadata {
    name      = var.service_name
    namespace = var.k8s_namespace
    labels = {
      app = var.service_name
    }
  }

  spec {
    selector = {
      app = var.service_name
    }

    port {
      port        = var.service_port
      target_port = var.service_port
    }

    type = local.service_type
  }
}

resource "kubernetes_ingress_v1" "ingress" {
  count = local.create_ingress ? 1 : 0
  metadata {
    name      = var.service_name
    namespace = var.k8s_namespace
    labels = {
      app = var.service_name
    }

    annotations = local.ingress_annotations
  }

  spec {
    dynamic "rule" {
      for_each = var.service_type == "EXTERNAL" ? [0] : []
      content {
        http {
          path {
            backend {
              service {
                name = "redirect"
                port {
                  name = "use-annotation"
                }
              }
            }

            path = "/*"
          }
        }
      }
    }

    rule {
      host = var.host_match

      http {
        path {
          path = "/*"
          backend {
            service {
              name = kubernetes_service.service.metadata.0.name
              port {
                number = var.service_port
              }
            }
          }
        }
      }
    }
  }
}


module "ingress" {
  count           = var.create_ingress ? 1 : 0
  source          = "../happy-ingress-eks"
  ingress_name    = var.service_name
  cloud_env       = var.cloud_env
  k8s_namespace   = var.k8s_namespace
  host_match      = var.host_match
  service_type    = var.service_type
  certificate_arn = var.certificate_arn
  tags_string     = var.tags_string
  backends = [
    {
      service_name = var.service_name
      service_port = var.service_port
      path         = "/*"
    }
  ]
}
