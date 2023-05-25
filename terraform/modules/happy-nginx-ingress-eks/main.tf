locals {
  ingress_annotations = {}
}

resource "kubernetes_ingress_v1" "ingress" {
  metadata {
    name        = var.ingress_name
    namespace   = var.k8s_namespace
    annotations = local.ingress_annotations
    labels      = var.labels
  }

  spec {
    ingress_class_name = "nginx"
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