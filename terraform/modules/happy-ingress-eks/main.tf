locals {
  ingress_base_annotations = {
    "kubernetes.io/ingress.class"                            = "alb"
    "alb.ingress.kubernetes.io/healthcheck-interval-seconds" = var.aws_alb_healthcheck_interval_seconds
    "alb.ingress.kubernetes.io/backend-protocol"             = var.target_service_scheme
    "alb.ingress.kubernetes.io/healthcheck-path"             = var.health_check_path
    "alb.ingress.kubernetes.io/healthcheck-protocol"         = var.target_service_scheme
    "alb.ingress.kubernetes.io/listen-ports"                 = jsonencode([{ HTTPS = 443 }, { HTTP = 80 }])
    # All ingresses are "internet-facing". If a service_type was marked "INTERNAL", it will be protected using OIDC.
    "alb.ingress.kubernetes.io/scheme"                  = var.routing.service_type == "VPC" ? "internal" : "internet-facing"
    "alb.ingress.kubernetes.io/subnets"                 = join(",", var.cloud_env.public_subnets)
    "alb.ingress.kubernetes.io/success-codes"           = var.routing.success_codes
    "alb.ingress.kubernetes.io/tags"                    = var.tags_string
    "alb.ingress.kubernetes.io/target-group-attributes" = "deregistration_delay.timeout_seconds=60"
    #    "alb.ingress.kubernetes.io/target-type"             = "instance"
    "alb.ingress.kubernetes.io/target-type" = "ip"
    "alb.ingress.kubernetes.io/group.name"  = var.routing.group_name
    "alb.ingress.kubernetes.io/group.order" = var.routing.priority
    "alb.ingress.kubernetes.io/load-balancer-attributes" = join(",", [ // Add any additional load-balancer-attributes here
      "idle_timeout.timeout_seconds=${var.routing.alb_idle_timeout}",
    ])
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

  # Attach a WAF if provided. Otherwise, ignore
  ingress_wafv2_annotations = var.regional_wafv2_arn != null ? {
    "alb.ingress.kubernetes.io/wafv2-acl-arn" = var.regional_wafv2_arn
  } : {}

  ingress_sg_annotations = (var.routing.service_type == "VPC") ? {
    "alb.ingress.kubernetes.io/security-groups"                     = aws_security_group.alb_sg[0].id
    "alb.ingress.kubernetes.io/manage-backend-security-group-rules" = "true"
  } : {}

  # All the annotations you want by default
  default_ingress_annotations = merge(
    local.ingress_tls_annotations,
    local.ingress_base_annotations,
    local.ingress_wafv2_annotations,
    var.additional_annotations,
  )

  additional_ingress_annotations = {
    // For VPC only routing, we want to configure the security group ourselves.
    "VPC" = local.ingress_sg_annotations
    // For internal routing, add auth annotations.
    "INTERNAL" = local.ingress_auth_annotations
  }
  ingress_annotations = merge(
    local.default_ingress_annotations,
    lookup(local.additional_ingress_annotations, var.routing.service_type, {}),
  )

  // the length of bypasses should never be bigger than the priority due to how we do the spread priority
  // in happy-stack-eks. We also have a validation on the input variable to ensure this.
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
        "alb.ingress.kubernetes.io/listen-ports" = jsonencode([{ HTTPS = 443 }])
      },
      // add our bypass conditions
      {
        "alb.ingress.kubernetes.io/conditions.${var.target_service_name}" = jsonencode([
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

// ALB's security group
resource "aws_security_group" "alb_sg" {
  count       = var.routing.service_type == "VPC" ? 1 : 0
  name        = "${var.ingress_name}-sg"
  description = "Security group for the ${var.ingress_name} alb."

  vpc_id = var.cloud_env.vpc_id

  // ingress from other security groups at the listen ports
  dynamic "ingress" {
    for_each = var.ingress_security_groups
    content {
      from_port       = 443
      to_port         = 443
      protocol        = "tcp"
      security_groups = [ingress.value]
    }
  }

  dynamic "ingress" {
    for_each = var.ingress_security_groups
    content {
      from_port       = 80
      to_port         = 80
      protocol        = "tcp"
      security_groups = [ingress.value]
    }
  }

  lifecycle {
    create_before_destroy = true
    ignore_changes        = [name, description]
  }
}

resource "kubernetes_ingress_v1" "ingress_bypasses" {
  for_each = local.ingress_bypass_annotations
  metadata {
    name        = replace("${var.ingress_name}-${each.key}-bypass", "_", "-")
    namespace   = var.k8s_namespace
    annotations = each.value
    labels      = var.labels
  }

  spec {
    rule {
      host = var.routing.host_match
      http {
        path {
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

resource "kubernetes_ingress_v1" "ingress" {
  metadata {
    name        = var.ingress_name
    namespace   = var.k8s_namespace
    annotations = local.ingress_annotations
    labels      = var.labels
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
