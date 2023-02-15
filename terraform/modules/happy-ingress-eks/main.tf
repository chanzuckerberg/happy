locals {
  ingress_base_annotations = {
    "kubernetes.io/ingress.class"                    = "alb"
    "alb.ingress.kubernetes.io/backend-protocol"     = "HTTP"
    "alb.ingress.kubernetes.io/healthcheck-path"     = var.health_check_path
    "alb.ingress.kubernetes.io/healthcheck-protocol" = "HTTP"
    # All ingresses are "internet-facing". If a service_type was marked "INTERNAL", it will be protected using OIDC.
    "alb.ingress.kubernetes.io/scheme"                  = "internet-facing"
    "alb.ingress.kubernetes.io/subnets"                 = join(",", var.cloud_env.public_subnets)
    "alb.ingress.kubernetes.io/success-codes"           = var.routing.success_codes
    "alb.ingress.kubernetes.io/tags"                    = var.tags_string
    "alb.ingress.kubernetes.io/target-group-attributes" = "deregistration_delay.timeout_seconds=60"
    "alb.ingress.kubernetes.io/target-type"             = "instance"
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

  ingress_annotations = (
    var.routing.service_type == "EXTERNAL" ?
    merge(local.ingress_tls_annotations, local.ingress_base_annotations) :
    merge(local.ingress_tls_annotations, local.ingress_auth_annotations, local.ingress_base_annotations)
  )

  // as long as priority is a positive number, the length of bypasses should never be bigger than the priority
  // due to how we do the spread priority in happy-stack-eks. TODO: add validation on this
  priority_spread = range(var.routing.priority - length(var.routing.bypasses), var.routing.priority)
  routing_keys    = keys(var.routing.bypasses)
  ingress_bypass_annotations = ({
    for k, v in zipmap(local.routing_keys, local.priority_spread) : k =>
    merge(
      local.ingress_base_annotations,
      // override the base group order with the priority spread values
      {
        "alb.ingress.kubernetes.io/group.order" = v
      },
      // the bypass rules should only be on the 443 listener
      {
        "alb.ingress.kubernetes.io/listen-ports"=jsonencode([{HTTPS=443}])
      },
      // add our bypass conditions
      {
        "alb.ingress.kubernetes.io/conditions.${var.routing.service_name}" = jsonencode([
          (length(var.routing.bypasses[k].methods) != 0 ? {
            field = "http-request-method"
            httpRequestMethodConfig = {
              Values = var.routing.bypasses[k].methods
            }
          } : null),
          (length(var.routing.bypasses[k].paths) != 0 ? {
            field = "path-pattern"
            pathPatternConfig = {
              Values = var.routing.bypasses[k].paths
            }
          } : null)
        ])
    })
  })
}

resource "kubernetes_ingress_v1" "ingress_bypasses" {
  for_each = local.ingress_bypass_annotations
  metadata {
    name      = replace("${var.ingress_name}-${each.key}-bypass", "_", "-")
    namespace = var.k8s_namespace
    labels = {
      app = var.ingress_name
    }

    annotations = each.value
  }

  spec {
    rule {
      host = var.routing.host_match
      http {
        path {
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
