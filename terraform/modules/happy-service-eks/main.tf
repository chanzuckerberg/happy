data "aws_region" "current" {}

locals {
  tags_string  = join(",", [for key, val in local.routing_tags : "${key}=${val}"])
  service_type = (var.routing.service_type == "PRIVATE" || var.routing.service_mesh) ? "ClusterIP" : "NodePort"
  match_labels = {
    app = var.routing.service_name
  }
  labels = merge({
    app                            = var.routing.service_name
    "app.kubernetes.io/name"       = var.stack_name
    "app.kubernetes.io/component"  = var.routing.service_name
    "app.kubernetes.io/part-of"    = var.stack_name
    "app.kubernetes.io/managed-by" = "happy"
  }, var.additional_pod_labels)

  external_dns_exclude_annotation = {
    "external-dns.alpha.kubernetes.io/exclude" = "true"
  }

  base_ingress_variables = {
    ingress_name            = var.routing.service_name
    target_service_port     = var.routing.service_mesh ? 443 : var.routing.service_port
    target_service_name     = var.routing.service_mesh ? "nginx-ingress-ingress-nginx-controller" : var.routing.service_name
    target_service_scheme   = var.routing.service_mesh ? "HTTPS" : var.routing.service_scheme
    cloud_env               = var.cloud_env
    k8s_namespace           = var.routing.service_mesh ? "nginx-encrypted-ingress" : var.k8s_namespace
    certificate_arn         = var.certificate_arn
    tags_string             = local.tags_string
    labels                  = local.labels
    regional_wafv2_arn      = var.regional_wafv2_arn
    ingress_security_groups = var.ingress_security_groups
  }
}

resource "kubernetes_deployment_v1" "deployment" {
  count = var.routing.service_type == "IMAGE_TEMPLATE" ? 0 : 1
  metadata {
    name      = var.routing.service_name
    namespace = var.k8s_namespace
    labels    = local.labels
    annotations = {
      "ad.datadoghq.com/tags" = jsonencode({
        "happy_stack"      = var.stack_name
        "happy_service"    = var.routing.service_name
        "deployment_stage" = var.deployment_stage
        "owner"            = var.tags.owner
        "project"          = var.tags.project
        "env"              = var.tags.env
        "service"          = var.tags.service
        "managedby"        = "happy"
        "happy_compute"    = "eks"
      })
      "linkerd.io/inject" = var.routing.service_mesh ? "enabled" : "disabled"
    }
  }

  wait_for_rollout = var.wait_for_steady_state

  spec {
    replicas = var.desired_count

    strategy {
      type = "RollingUpdate"
      rolling_update {
        max_surge       = "25%"
        max_unavailable = "25%"
      }
    }

    selector {
      match_labels = local.match_labels
    }

    template {
      metadata {
        labels = local.labels

        annotations = merge({
          "ad.datadoghq.com/tags" = jsonencode({
            "happy_stack"      = var.stack_name
            "happy_service"    = var.routing.service_name
            "deployment_stage" = var.deployment_stage
            "owner"            = var.tags.owner
            "project"          = var.tags.project
            "env"              = var.tags.env
            "service"          = var.tags.service
            "managedby"        = "happy"
            "happy_compute"    = "eks"
          })
          }, var.routing.service_mesh ? {
          "linkerd.io/inject"                        = "enabled"
          "config.linkerd.io/default-inbound-policy" = "all-authenticated"
          //Skipping all ports listed here https://linkerd.io/2.13/features/protocol-detection/
          "config.linkerd.io/skip-outbound-ports" = "25,587,3306,4444,4567,4568,5432,6379,9300,11211"
        } : {})
      }

      spec {
        service_account_name = var.aws_iam.service_account_name == null ? module.iam_service_account.service_account_name : var.aws_iam.service_account_name

        dynamic "toleration" {
          for_each = var.gpu != null ? [1] : []
          content {
            key      = "nvidia.com/gpu.present"
            operator = "Exists"
            effect   = "NoSchedule"
          }
        }
        dynamic "toleration" {
          for_each = var.gpu != null ? [1] : []
          content {
            key      = "nvidia.com/gpu"
            operator = "Exists"
            effect   = "NoSchedule"
          }
        }

        topology_spread_constraint {
          max_skew = 3
          #TODO: Once min_domains are supported, uncomment line below. https://github.com/hashicorp/terraform-provider-kubernetes/issues/2292
          #min_domains        = 3
          topology_key       = "topology.kubernetes.io/zone"
          when_unsatisfiable = "DoNotSchedule"
          label_selector {
            match_labels = local.match_labels
          }
        }

        affinity {
          node_affinity {
            required_during_scheduling_ignored_during_execution {
              node_selector_term {
                match_expressions {
                  key      = "kubernetes.io/arch"
                  operator = "In"
                  values   = [var.platform_architecture]
                }
              }
            }
          }
          pod_anti_affinity {
            preferred_during_scheduling_ignored_during_execution {
              weight = 100
              pod_affinity_term {
                topology_key = "kubernetes.io/hostname"
                label_selector {
                  match_labels = local.match_labels
                }
              }
            }
            preferred_during_scheduling_ignored_during_execution {
              weight = 100
              pod_affinity_term {
                topology_key = "topology.kubernetes.io/zone"
                label_selector {
                  match_labels = local.match_labels
                }
              }
            }
          }
        }

        restart_policy = "Always"

        container {
          name              = var.container_name
          image             = "${module.ecr.repository_url}:${var.image_tag}"
          command           = var.cmd
          args              = var.args
          image_pull_policy = var.image_pull_policy

          env {
            name  = "DEPLOYMENT_STAGE"
            value = var.deployment_stage
          }
          env {
            name  = "AWS_REGION"
            value = data.aws_region.current.name
          }
          env {
            name  = "AWS_DEFAULT_REGION"
            value = data.aws_region.current.name
          }

          env {
            name  = "HAPPY_STACK"
            value = var.stack_name
          }

          env {
            name  = "HAPPY_SERVICE"
            value = var.container_name
          }

          env {
            name  = "HAPPY_CONTAINER"
            value = var.container_name
          }

          dynamic "env" {
            for_each = var.service_endpoints
            content {
              name  = replace(env.key, "-", "_")
              value = env.value
            }
          }

          dynamic "env" {
            for_each = var.additional_env_vars
            content {
              name  = env.key
              value = env.value
            }
          }

          dynamic "env_from" {
            for_each = toset(var.additional_env_vars_from_config_maps.items)
            content {
              prefix = var.additional_env_vars_from_config_maps.prefix
              config_map_ref {
                name = env_from.value
              }
            }
          }

          dynamic "env_from" {
            for_each = toset(var.additional_env_vars_from_secrets.items)
            content {
              prefix = var.additional_env_vars_from_secrets.prefix
              secret_ref {
                name = env_from.value
              }
            }
          }

          port {
            name           = "http"
            container_port = var.routing.port
          }

          resources {
            limits = {
              cpu              = var.cpu
              memory           = var.memory
              "nvidia.com/gpu" = var.gpu
            }
            requests = {
              cpu              = var.cpu_requests
              memory           = var.memory_requests
              "nvidia.com/gpu" = var.gpu_requests
            }
          }

          volume_mount {
            mount_path = "/var/happy"
            name       = "integration-secret"
            read_only  = true
          }

          dynamic "volume_mount" {
            for_each = toset(var.additional_volumes_from_secrets.items)
            content {
              mount_path = "${var.additional_volumes_from_secrets.base_dir}/${volume_mount.value}"
              name       = volume_mount.value
              read_only  = true
            }
          }

          dynamic "volume_mount" {
            for_each = toset(var.additional_volumes_from_config_maps.items)
            content {
              mount_path = "/var/${volume_mount.value}"
              name       = volume_mount.value
              read_only  = true
            }
          }

          liveness_probe {
            http_get {
              path   = var.health_check_path
              port   = var.routing.health_check_port == null ? var.routing.port : var.routing.healthcheck_port
              scheme = var.routing.scheme
            }

            initial_delay_seconds = var.initial_delay_seconds
            period_seconds        = var.period_seconds
            timeout_seconds       = var.liveness_timeout_seconds
          }

          readiness_probe {
            http_get {
              path   = var.health_check_path
              port   = var.routing.port
              scheme = var.routing.scheme
            }

            initial_delay_seconds = var.initial_delay_seconds
            period_seconds        = var.period_seconds
            timeout_seconds       = var.readiness_timeout_seconds
          }
        }

        dynamic "container" {
          for_each = var.sidecars
          content {
            image             = "${container.value.image}:${container.value.tag}"
            name              = container.key
            image_pull_policy = container.value.image_pull_policy

            port {
              name           = "http"
              container_port = container.value.port
            }

            resources {
              limits = {
                cpu    = container.value.cpu
                memory = container.value.memory
              }
              requests = {
                cpu    = container.value.cpu
                memory = container.value.memory
              }
            }

            liveness_probe {
              http_get {
                path   = container.value.health_check_path
                port   = container.value.port
                scheme = container.value.scheme
              }

              initial_delay_seconds = container.value.initial_delay_seconds
              period_seconds        = container.value.period_seconds
              timeout_seconds       = container.value.liveness_timeout_seconds
            }

            readiness_probe {
              http_get {
                path   = container.value.health_check_path
                port   = container.value.port
                scheme = container.value.scheme
              }

              initial_delay_seconds = container.value.initial_delay_seconds
              period_seconds        = container.value.period_seconds
              timeout_seconds       = container.value.readiness_timeout_seconds
            }

            dynamic "volume_mount" {
              for_each = toset(var.additional_volumes_from_secrets.items)
              content {
                mount_path = "${var.additional_volumes_from_secrets.base_dir}/${volume_mount.value}"
                name       = volume_mount.value
                read_only  = true
              }
            }

            dynamic "volume_mount" {
              for_each = toset(var.additional_volumes_from_config_maps.items)
              content {
                mount_path = "/var/${volume_mount.value}"
                name       = volume_mount.value
                read_only  = true
              }
            }

            env {
              name  = "DEPLOYMENT_STAGE"
              value = var.deployment_stage
            }
            env {
              name  = "AWS_REGION"
              value = data.aws_region.current.name
            }
            env {
              name  = "AWS_DEFAULT_REGION"
              value = data.aws_region.current.name
            }

            env {
              name  = "HAPPY_STACK"
              value = var.stack_name
            }

            env {
              name  = "HAPPY_SERVICE"
              value = var.container_name
            }

            env {
              name  = "HAPPY_CONTAINER"
              value = container.key
            }

            dynamic "env_from" {
              for_each = toset(var.additional_env_vars_from_secrets.items)
              content {
                prefix = var.additional_env_vars_from_secrets.prefix
                secret_ref {
                  name = env_from.value
                }
              }
            }

            dynamic "env_from" {
              for_each = toset(var.additional_env_vars_from_config_maps.items)
              content {
                prefix = var.additional_env_vars_from_config_maps.prefix
                config_map_ref {
                  name = env_from.value
                }
              }
            }

            dynamic "env" {
              for_each = var.additional_env_vars
              content {
                name  = env.key
                value = env.value
              }
            }
          }
        }

        volume {
          name = "integration-secret"
          secret {
            secret_name = "integration-secret"
          }
        }

        dynamic "volume" {
          for_each = toset(var.additional_volumes_from_secrets.items)
          content {
            secret {
              secret_name = volume.value
            }
            name = volume.value
          }
        }

        dynamic "volume" {
          for_each = toset(var.additional_volumes_from_config_maps.items)
          content {
            config_map {
              name = volume.value
            }
            name = volume.value
          }
        }
      }
    }
  }
}

resource "kubernetes_service_v1" "service" {
  count = var.routing.service_type == "IMAGE_TEMPLATE" ? 0 : 1
  metadata {
    name      = var.routing.service_name
    namespace = var.k8s_namespace
    labels    = local.labels
  }

  spec {
    selector = {
      app = var.routing.service_name
    }

    port {
      port        = var.routing.service_port
      target_port = var.routing.service_port
    }

    type = local.service_type
  }
}

module "ingress" {
  count  = (var.routing.service_type == "EXTERNAL" || var.routing.service_type == "INTERNAL" || var.routing.service_type == "VPC") ? 1 : 0
  source = "../happy-ingress-eks"

  ingress_name            = local.base_ingress_variables.ingress_name
  target_service_port     = local.base_ingress_variables.target_service_port
  target_service_name     = local.base_ingress_variables.target_service_name
  target_service_scheme   = local.base_ingress_variables.target_service_scheme
  cloud_env               = local.base_ingress_variables.cloud_env
  k8s_namespace           = local.base_ingress_variables.k8s_namespace
  certificate_arn         = local.base_ingress_variables.certificate_arn
  tags_string             = local.base_ingress_variables.tags_string
  labels                  = local.base_ingress_variables.labels
  regional_wafv2_arn      = local.base_ingress_variables.regional_wafv2_arn
  ingress_security_groups = local.base_ingress_variables.ingress_security_groups

  routing = var.routing
}

module "ingress_exclude_external_dns" {
  for_each = (var.routing.service_type == "EXTERNAL" || var.routing.service_type == "INTERNAL" || var.routing.service_type == "VPC") ? var.routing.additional_hostnames : []
  source   = "../happy-ingress-eks"

  ingress_name            = replace("${local.base_ingress_variables.ingress_name}-${each.value}", ".", "-")
  target_service_port     = local.base_ingress_variables.target_service_port
  target_service_name     = local.base_ingress_variables.target_service_name
  target_service_scheme   = local.base_ingress_variables.target_service_scheme
  cloud_env               = local.base_ingress_variables.cloud_env
  k8s_namespace           = local.base_ingress_variables.k8s_namespace
  certificate_arn         = local.base_ingress_variables.certificate_arn
  tags_string             = local.base_ingress_variables.tags_string
  labels                  = local.base_ingress_variables.labels
  regional_wafv2_arn      = local.base_ingress_variables.regional_wafv2_arn
  ingress_security_groups = local.base_ingress_variables.ingress_security_groups

  routing                = merge(var.routing, { host_match : "" })
  additional_annotations = local.external_dns_exclude_annotation
}

module "nginx-ingress" {
  for_each = ((var.routing.service_type == "EXTERNAL" || var.routing.service_type == "INTERNAL" || var.routing.service_type == "VPC") && var.routing.service_mesh) ? setunion([var.routing.host_match], var.routing.additional_hostnames) : []
  source   = "../happy-nginx-ingress-eks"

  ingress_name        = replace("${var.routing.service_name}-${each.value}-nginx", ".", "-")
  k8s_namespace       = var.k8s_namespace
  host_match          = each.value
  host_path           = replace(var.routing.path, "/\\*$/", "") //NGINX does not support paths that end with *
  target_service_name = var.routing.service_name
  target_service_port = var.routing.service_port
  timeout             = var.routing.alb_idle_timeout
  labels              = local.labels
}

module "mesh-access-control" {
  count               = var.routing.service_mesh && var.routing.allow_mesh_services != null ? 1 : 0
  source              = "../happy-mesh-access-control"
  k8s_namespace       = var.k8s_namespace
  service_port        = var.routing.service_port
  service_name        = var.routing.service_name
  service_type        = var.routing.service_type
  deployment_stage    = var.deployment_stage
  allow_mesh_services = var.routing.allow_mesh_services
  labels              = local.labels
}

resource "kubernetes_horizontal_pod_autoscaler_v1" "hpa" {
  count = var.routing.service_type == "IMAGE_TEMPLATE" ? 0 : 1
  metadata {
    name      = var.routing.service_name
    namespace = var.k8s_namespace
    labels    = local.labels
  }

  spec {
    max_replicas = var.max_count
    min_replicas = var.desired_count

    target_cpu_utilization_percentage = var.scaling_cpu_threshold_percentage

    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = kubernetes_deployment_v1.deployment[0].metadata[0].name
    }
  }
}

resource "kubernetes_pod_disruption_budget_v1" "pdb" {
  count = var.routing.service_type == "IMAGE_TEMPLATE" || var.max_unavailable_count >= var.desired_count ? 0 : 1
  metadata {
    name      = var.routing.service_name
    namespace = var.k8s_namespace
    labels    = local.labels
  }

  spec {
    max_unavailable = var.max_unavailable_count
    selector {
      match_labels = local.match_labels
    }
  }
}
