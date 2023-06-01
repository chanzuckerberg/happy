locals {
  allow_ingress_controller = var.routing.service_type == "EXTERNAL" || var.routing.service_type == "INTERNAL"
  authorization_enabled    = var.routing.service_mesh && var.routing.allow_mesh_services != null
  needs_policy             = local.authorization_enabled && (local.allow_ingress_controller || length(var.routing.allow_mesh_services) > 0)
}

resource "kubernetes_manifest" "linkerd_server" {
  count = local.authorization_enabled ? 1 : 0
  manifest = {
    "apiVersion" = "policy.linkerd.io/v1alpha1"
    "kind"       = "Server"
    "metadata" = {
      "name"      = "${var.routing.service_name}-server"
      "namespace" = var.k8s_namespace
    }
    "spec" = {
      "port" = var.routing.service_port
      "podSelector" = {
        "matchLabels" = {
          "app" = var.routing.service_name
        }
      }
    }
  }
}

resource "kubernetes_manifest" "linkerd_mesh_tls_authentication" {
  count = local.needs_policy ? 1 : 0
  manifest = {
    "apiVersion" = "policy.linkerd.io/v1alpha1"
    "kind"       = "MeshTLSAuthentication"
    "metadata" = {
      "name"      = "${var.routing.service_name}-mesh-tls-auth"
      "namespace" = var.k8s_namespace
    }
    "spec" = {
      "identityRefs" = concat([for v in var.routing.allow_mesh_services : {
        "kind" = "ServiceAccount"
        "name" = "${v.stack}-${v.service}-${var.deployment_stage}-${v.stack}"
        }], local.allow_ingress_controller ? [{
        "kind" = "ServiceAccount"
        "name" = "nginx-ingress-ingress-nginx"
      }] : [])
    }
  }
}

resource "kubernetes_manifest" "linkerd_authorization_policy" {
  count = local.needs_policy ? 1 : 0
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
      "requiredAuthenticationRefs" = [{
        "group" = "policy.linkerd.io"
        "kind"  = "MeshTLSAuthentication"
        "name"  = "${var.routing.service_name}-mesh-tls-auth"
      }]
    }
  }
}
