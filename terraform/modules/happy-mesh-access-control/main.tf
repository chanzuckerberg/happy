locals {
  allow_ingress_controller = var.service_type == "EXTERNAL" || var.service_type == "INTERNAL" || var.service_type == "VPC"
  needs_policy             = local.allow_ingress_controller || length(var.allow_mesh_services) > 0
}

resource "kubernetes_manifest" "linkerd_server" {
  manifest = {
    "apiVersion" = "policy.linkerd.io/v1alpha1"
    "kind"       = "Server"
    "metadata" = {
      "name"      = "${var.service_name}-server"
      "namespace" = var.k8s_namespace
      "labels"    = var.labels
    }
    "spec" = {
      "port" = var.service_port
      "podSelector" = {
        "matchLabels" = {
          "app" = var.service_name
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
      "name"      = "${var.service_name}-mesh-tls-auth"
      "namespace" = var.k8s_namespace
      "labels"    = var.labels
    }
    "spec" = {
      "identityRefs" = concat([for v in var.allow_mesh_services : {
        "kind"      = "ServiceAccount"
        "name"      = v.service_account_name != null ? v.service_account_name : "${v.stack}-${v.service}-${var.deployment_stage}-${v.stack}"
        "namespace" = var.k8s_namespace
        }], local.allow_ingress_controller ? [{
        "kind"      = "ServiceAccount"
        "name"      = "nginx-ingress-ingress-nginx"
        "namespace" = "nginx-encrypted-ingress"
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
      "name"      = "${var.service_name}-policy"
      "namespace" = var.k8s_namespace
      "labels"    = var.labels
    }
    "spec" = {
      "targetRef" = {
        "group" = "policy.linkerd.io"
        "kind"  = "Server"
        "name"  = "${var.service_name}-server"
      }
      "requiredAuthenticationRefs" = [{
        "group" = "policy.linkerd.io"
        "kind"  = "MeshTLSAuthentication"
        "name"  = "${var.service_name}-mesh-tls-auth"
      }]
    }
  }
}
