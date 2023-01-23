locals {
  create_ingress = var.service_type == "EXTERNAL" || var.service_type == "INTERNAL"
  ingress_base_annotations = {
    "alb.ingress.kubernetes.io/backend-protocol"        = "HTTP"
    "alb.ingress.kubernetes.io/healthcheck-path"        = var.health_check_path
    "alb.ingress.kubernetes.io/healthcheck-protocol"    = "HTTP"
    "alb.ingress.kubernetes.io/listen-ports"            = var.service_type == "EXTERNAL" ? jsonencode(local.listen_ports_tls) : jsonencode(local.listen_ports_base)
    "alb.ingress.kubernetes.io/scheme"                  = local.scheme
    "alb.ingress.kubernetes.io/subnets"                 = join(",", var.cloud_env.public_subnets)
    "alb.ingress.kubernetes.io/success-codes"           = var.success_codes
    "alb.ingress.kubernetes.io/tags"                    = local.tags_string
    "alb.ingress.kubernetes.io/target-group-attributes" = "deregistration_delay.timeout_seconds=60"
    "alb.ingress.kubernetes.io/target-type"             = "instance"
    "kubernetes.io/ingress.class"                       = "alb"
  }
  ingress_tls_annotations = {
    "alb.ingress.kubernetes.io/actions.redirect" = jsonencode({
      Type = "redirect",
      RedirectConfig = {
        Protocol   = "HTTPS",
        Port       = "443",
        StatusCode = "HTTP_301",
      }
    })
    "alb.ingress.kubernetes.io/certificate-arn" = var.certificate_arn
    "alb.ingress.kubernetes.io/ssl-policy"      = "ELBSecurityPolicy-TLS-1-2-2017-01"
  }
  ingress_annotations      = var.service_type == "EXTERNAL" ? merge(local.ingress_tls_annotations, local.ingress_base_annotations) : local.ingress_base_annotations
  ingress_group_name       = "${var.stack_name}-${var.deployment_stage}"
  ingress_group_annotation = { alb.ingress.kubernetes.io / group.name : local.ingress_group_name }

}
// https://repost.aws/questions/QUEyFKpZCBR_OTFMQlJNypaQ/ingress-annotations-only-for-a-specific-path
// https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.2/guide/ingress/annotations/#ingressgroup
// The following ingresses great an IngressGroup. AWS ALB will only make 1 ALB and merge the rules. This allows
// us to do both authenticated requests and unauthenticated requests in the same ALB.

// ingress to handle unauthenticated requests
resource "kubernetes_ingress_v1" "health_check_unauthenticated" {
  count = local.create_ingress ? 1 : 0
  metadata {
    name      = var.service_name
    namespace = var.k8s_namespace
    labels = {
      app = var.service_name
    }

    annotations = merge(local.ingress_annotations, local.ingress_group_annotation)
  }

  spec {
    rule {
      host = var.host_match

      http {
        path {
          path = "/${var.health_check_path}"
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

# ingress to handle all authenticated requests
resource "kubernetes_ingress_v1" "authenticated" {
  count = local.create_ingress ? 1 : 0
  metadata {
    name      = var.service_name
    namespace = var.k8s_namespace
    labels = {
      app = var.service_name
    }

    annotations = merge(local.ingress_annotations, local.ingress_group_annotation)
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
