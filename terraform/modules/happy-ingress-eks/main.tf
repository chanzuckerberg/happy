locals {
  ingress_base_annotations = {
    "alb.ingress.kubernetes.io/backend-protocol"     = "HTTP"
    "alb.ingress.kubernetes.io/healthcheck-path"     = var.health_check_path
    "alb.ingress.kubernetes.io/healthcheck-protocol" = "HTTP"
    # All ingresses are "internet-facing" so we need them all to listen on TLS
    "alb.ingress.kubernetes.io/listen-ports" = jsonencode([{ "HTTP" : 80 }, { "HTTPS" : 443 }])
    # All ingresses are "internet-facing". If a service_type was marked "INTERNAL", it will be protected using OIDC.
    "alb.ingress.kubernetes.io/scheme"                  = "internet-facing"
    "alb.ingress.kubernetes.io/subnets"                 = join(",", var.cloud_env.public_subnets)
    "alb.ingress.kubernetes.io/success-codes"           = var.routing.success_codes
    "alb.ingress.kubernetes.io/tags"                    = var.tags_string
    "alb.ingress.kubernetes.io/target-group-attributes" = "deregistration_delay.timeout_seconds=60"
    "alb.ingress.kubernetes.io/target-type"             = "instance"
    "kubernetes.io/ingress.class"                       = "alb"
    "alb.ingress.kubernetes.io/group.name"              = var.routing.group_name
    "alb.ingress.kubernetes.io/group.order"             = var.routing.priority
  }

  redirect_action_name = "redirect"
  ingress_tls_annotations = {
    "alb.ingress.kubernetes.io/actions.${local.redirect_action_name}" = jsonencode({
      Type = local.redirect_action_name
      RedirectConfig = {
        Protocol   = "HTTPS"
        Port       = "443"
        StatusCode = "HTTP_301"
      }
    })
    "alb.ingress.kubernetes.io/certificate-arn" = var.certificate_arn
    "alb.ingress.kubernetes.io/ssl-policy"      = "ELBSecurityPolicy-TLS-1-2-2017-01"
  }

  ingress_auth_annotations = {
    "alb.ingress.kubernetes.io/auth-type"                       = "oidc"
    "alb.ingress.kubernetes.io/auth-on-unauthenticated-request" = "authenticate"
    "alb.ingress.kubernetes.io/auth-idp-oidc"                   = jsonencode(var.routing.oidc_config)
  }

  bypass_annotations = { for k, v in var.routing.bypasses : "alb.ingress.kubernetes.io/conditions.${k}" => [
    {
      field = "http-request-method"
      httpRequestMethodConfig = {
        Values = v.methods
      }
    },
    {
      field = "path-pattern"
      pathPatternConfig = {
        Values = v.paths
      }
    }]
  }

  ingress_annotations = (
    var.routing.service_type == "EXTERNAL" ?
    merge(local.ingress_tls_annotations, local.ingress_base_annotations) :
    merge(local.ingress_tls_annotations, local.ingress_auth_annotations, local.ingress_base_annotations)
  )

  ingress_bypass_annotations = merge(local.ingress_tls_annotations, local.ingress_base_annotations, local.bypass_annotations)
}

resource "kubernetes_ingress_v1" "ingress_options_bypass" {
  for_each = var.routing.bypasses
  metadata {
    name      = var.ingress_name
    namespace = var.k8s_namespace
    labels = {
      app = var.ingress_name
    }

    annotations = local.ingress_bypass_annotations
  }

  spec {
    rule {
      http {
        path {
          backend {
            service {
              name = local.redirect_action_name
              port {
                name = "use-annotation"
              }
            }
          }

          path = "/*"
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
              name = each.value.key
              port {
                number = "use-annotation"
              }
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_ingress_v1" "ingress" {
  metadata {
    name      = var.ingress_name
    namespace = var.k8s_namespace
    labels = {
      app = var.ingress_name
    }

    annotations = local.ingress_annotations
  }

  spec {
    rule {
      http {
        path {
          backend {
            service {
              name = local.redirect_action_name
              port {
                name = "use-annotation"
              }
            }
          }

          path = "/*"
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
