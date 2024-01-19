
locals {
  sticky_annotations = {
    "nginx.ingress.kubernetes.io/affinity"               = "cookie"
    "nginx.ingress.kubernetes.io/session-cookie-name"    = var.sticky_sessions.cookie_name
    "nginx.ingress.kubernetes.io/session-cookie-max-age" = var.sticky_sessions.duration_seconds
  }

  base_annotations = {
    "cert-manager.io/cluster-issuer"                    = "nginx-issuer"
    "linkerd.io/inject"                                 = "enabled"
    "external-dns.alpha.kubernetes.io/exclude"          = "true"
    "nginx.ingress.kubernetes.io/proxy-connect-timeout" = var.timeout
    "nginx.ingress.kubernetes.io/proxy-send-timeout"    = var.timeout
    "nginx.ingress.kubernetes.io/proxy-read-timeout"    = var.timeout
  }

  annotations = merge(local.base_annotations, var.sticky_sessions.enabled ? local.sticky_annotations : {})
}
resource "kubernetes_ingress_v1" "ingress" {
  metadata {
    name        = var.ingress_name
    namespace   = var.k8s_namespace
    annotations = local.annotations
    labels      = var.labels
  }

  spec {
    ingress_class_name = "nginx"
    tls {
      hosts = [
        var.host_match
      ]
      secret_name = "${var.ingress_name}-tls-secret"
    }
    rule {
      host = var.host_match
      http {
        path {
          path = var.host_path
          backend {
            service {
              name = var.target_service_name
              port {
                number = var.target_service_port
              }
            }
          }
        }
      }
    }
  }
}