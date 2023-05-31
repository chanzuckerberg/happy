locals {
  allow_ingress_controller = var.routing.service_type == "EXTERNAL" || var.routing.service_type == "INTERNAL"
}

resource "kubernetes_manifest" "linkerd_server" {
  count = var.routing.service_mesh == false || var.routing.allow_mesh_services == null ? 0 : 1
  manifest = {
    "apiVersion" = "policy.linkerd.io/v1alpha1"
    "kind"       = "Server"
    "metadata" = {
      "name"      = "${var.routing.service_name}-server"
      "namespace" = var.k8s_namespace
    }
    "spec" = {
      "port" = var.routing.service_port
      "proxyProtocol" = "HTTP/2" //Can make this configurable later
      "podSelector" = {
        "matchLabels" = {
          "app" = var.routing.service_name
        }
      }
    }
  }
}

resource "kubernetes_manifest" "linkerd_authorization_policy" {
  count = var.routing.service_mesh == false || var.routing.allow_mesh_services == null ? 0 : 1
  manifest = {
    "apiVersion" = "policy.linkerd.io/v1alpha1"
    "kind"       = "AuthorizationPolicy"
    "metadata" = {
      "name"      = "${var.routing.service_name}-policy"
      "namespace" = var.k8s_namespace
    }
    "spec" = {
      "targetRef" = {
        "group" = "policy.linkerd.io"
        "kind"  = "Server"
        "name"  = "${var.routing.service_name}-server"
      }
      "requiredAuthenticationRefs" = concat([for v in var.routing.allow_mesh_services: {
        "kind"  = "ServiceAccount"
        "name"  = "${v.service}-${var.deployment_stage}-${v.stack}"
      }], local.allow_ingress_controller ? [{
        "kind"  = "ServiceAccount"
        "name"  = "nginx-ingress-ingress-nginx"
      }] : [])
      "proxyProtocol" = "HTTP/2" //Can make this configurable later
      "podSelector" = {
        "matchLabels" = {
          "app" = var.routing.service_name
        }
      }
    }
  }
}
