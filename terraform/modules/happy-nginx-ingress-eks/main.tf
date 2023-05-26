resource "kubernetes_ingress_v1" "ingress" {
  metadata {
    name      = var.ingress_name
    namespace = var.k8s_namespace
    annotations = {
      "cert-manager.io/issuer"                       = "nginx-issuer"
      "nginx.ingress.kubernetes.io/service-upstream" = "true"
      "linkerd.io/inject"                            = "enabled"

    }
    labels = var.labels
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
          path = "/"
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