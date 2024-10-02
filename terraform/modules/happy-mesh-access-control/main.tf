locals {
  allow_ingress_controller     = var.service_type == "EXTERNAL" || var.service_type == "INTERNAL" || var.service_type == "VPC"
  allow_k6_operator_controller = var.deployment_stage == "rdev" || var.deployment_stage == "staging"
  needs_policy                 = local.allow_ingress_controller || length(var.allow_mesh_services) > 0
  # Service accounts that we want to allow access to this protected service
  mesh_services_service_accounts = [for v in var.allow_mesh_services : {
    "kind"      = "ServiceAccount"
    "name"      = v.service_account_name != null && v.service_account_name != "" ? v.service_account_name : "${v.stack}-${v.service}-${var.deployment_stage}-${v.stack}"
    "namespace" = var.k8s_namespace
  }]
  optional_ingress_controller_service_account = local.allow_ingress_controller ? [{
    "kind"      = "ServiceAccount"
    "name"      = "nginx-ingress-ingress-nginx"
    "namespace" = "nginx-encrypted-ingress"
  }] : []
  status_page_service_account = [{
    "kind"      = "ServiceAccount"
    "name"      = "edu-platform-${var.deployment_stage}-status-page"
    "namespace" = "status-page"
  }]
  k6_operator_service_account = local.allow_k6_operator_controller ? [{
    "kind"      = "ServiceAccount"
    "name"      = "k6-operator-controller"
    "namespace" = "k6-operator-system"
  }] : []
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
      "identityRefs" = concat(
        local.mesh_services_service_accounts,
        local.optional_ingress_controller_service_account,
        local.status_page_service_account,
        local.k6_operator_service_account
      )
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
