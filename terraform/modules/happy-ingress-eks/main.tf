data "aws_region" "current" {}

locals {
  create_ingress    = (var.service_type == "EXTERNAL" || var.service_type == "INTERNAL")
  scheme            = var.service_type == "EXTERNAL" ? "internet-facing" : "internal"
  listen_ports_base = [{ "HTTP" : 80 }]
  listen_ports_tls  = [merge(local.listen_ports_base[0], { "HTTPS" : 443 })]
  ingress_base_annotations = {
    "alb.ingress.kubernetes.io/backend-protocol"        = "HTTP"
    "alb.ingress.kubernetes.io/healthcheck-path"        = var.health_check_path
    "alb.ingress.kubernetes.io/healthcheck-protocol"    = "HTTP"
    "alb.ingress.kubernetes.io/listen-ports"            = var.service_type == "EXTERNAL" ? jsonencode(local.listen_ports_tls) : jsonencode(local.listen_ports_base)
    "alb.ingress.kubernetes.io/scheme"                  = local.scheme
    "alb.ingress.kubernetes.io/subnets"                 = join(",", var.cloud_env.public_subnets)
    "alb.ingress.kubernetes.io/success-codes"           = var.success_codes
    "alb.ingress.kubernetes.io/tags"                    = var.tags_string
    "alb.ingress.kubernetes.io/target-group-attributes" = "deregistration_delay.timeout_seconds=60"
    "alb.ingress.kubernetes.io/target-type"             = "instance"
    "kubernetes.io/ingress.class"                       = "alb"
    "alb.ingress.kubernetes.io/group.name"              = var.routing.group_name
    "alb.ingress.kubernetes.io/group.order"             = var.routing.priority
  }
  ingress_tls_annotations = {
    "alb.ingress.kubernetes.io/actions.redirect" = <<EOT
        {"Type": "redirect", "RedirectConfig": {"Protocol": "HTTPS", "Port": "443", "StatusCode": "HTTP_301"}}
      EOT
    "alb.ingress.kubernetes.io/certificate-arn"  = var.certificate_arn
    "alb.ingress.kubernetes.io/ssl-policy"       = "ELBSecurityPolicy-TLS-1-2-2017-01"
  }
  ingress_annotations = var.service_type == "EXTERNAL" ? merge(local.ingress_tls_annotations, local.ingress_base_annotations) : local.ingress_base_annotations
}

resource "kubernetes_ingress_v1" "ingress" {
  count = local.create_ingress ? 1 : 0
  metadata {
    name      = var.ingress_name
    namespace = var.k8s_namespace
    labels = {
      app = var.ingress_name
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

      host = var.routing.host_match
      http {
        path {
          path = var.routing.path
          backend {
            service {
              name = var.routing.service_name
              port {
                number = var.routing.service_port
              }
            }
          }
        }
      }
    }
  }
}
