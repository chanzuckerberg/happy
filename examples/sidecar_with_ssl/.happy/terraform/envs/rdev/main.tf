
 # update if different environment
locals {
  external_domain = "happy-playground-rdev.rdev.si.czi.technology"
  local_domain = "si-rdev-happy-eks-rdev-happy-env.svc.cluster.local"
}

resource "kubernetes_manifest" "issuer" {
  manifest = {
    "apiVersion" = "cert-manager.io/v1"
    "kind" = "ClusterIssuer"
    "metadata" = {
      "name" = "${var.stack_name}-cluster-issuer"
    }
    "spec" = {
      "selfSigned" = {}
    }
  }
}
resource "kubernetes_manifest" "ssl-cert" {
  manifest = {
    "apiVersion"  = "cert-manager.io/v1"
    "kind"        = "Certificate"
    "metadata"    = {
      "name"      = "${var.stack_name}-cert"
      "namespace" = "${var.k8s_namespace}"
    }
    "spec" = {
      "secretName" = "${var.stack_name}-tls-secret"
      "issuerRef"  = {
        "name"     = "${var.stack_name}-cluster-issuer"
        "kind"     = "ClusterIssuer"
      }
      "dnsNames"   = [
        "${var.stack_name}-api.${local.external_domain}",
        "${var.stack_name}-api.${local.local_domain}"
      ]
    }
  }
}

module "stack" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks?ref=main"

  image_tag        = var.image_tag
  image_tags       = jsondecode(var.image_tags)
  stack_name       = var.stack_name
  deployment_stage = "rdev"
  stack_prefix     = "/${var.stack_name}"
  k8s_namespace    = var.k8s_namespace

  services = {
    api = {
      name                  = "api"
      cpu                   = "100m"
      memory                = "100Mi"
      port                  = 3000
      service_type          = "INTERNAL"
      platform_architecture = "arm64"
      health_check_path     = "/"
      sidecars = {
        ssl-sidecar = {
          image  = "ACCOUNTID.dkr.ecr.us-west-2.amazonaws.com/ssl-sidecar"
          tag    = "0.0.6"
          port   = 8443
          scheme = "HTTPS"
          cpu    = "100m"
          memory = "128Mi"
          health_check_path = "/"
        }
      }
      service_port = 8443
      service_scheme = "HTTPS"
    }
  }

  additional_volumes_from_secrets = {
    items = ["${var.stack_name}-tls-secret"]
  }

  additional_volumes_from_config_maps = {
    items = ["stacklist"]
  }
  additional_env_vars = {
    "CERT_FILE" = "/var/${var.stack_name}-tls-secret/tls.crt"
    "CA_CERT_FILE" = "/var/${var.stack_name}-tls-secret/ca.crt"
    "KEY_FILE" = "/var/${var.stack_name}-tls-secret/tls.key"
    "STACK_LIST_FILE" = "/var/stacklist/stacklist"
    "API_URL" = "http://localhost:3000/stacklist"
  }
}
