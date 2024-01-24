data "aws_region" "current" {}


resource "kubernetes_cron_job_v1" "task_definition" {
  metadata {
    name      = var.task_name
    namespace = var.k8s_namespace
    labels = {
      stack = var.stack_name
    }
  }
  spec {
    concurrency_policy            = "Forbid"
    failed_jobs_history_limit     = var.failed_jobs_history_limit
    schedule                      = var.cron_schedule
    suspend                       = !var.is_cron_job // This cron job is suspended by default to be used to create jobs on demand
    starting_deadline_seconds     = var.starting_deadline_seconds
    successful_jobs_history_limit = var.successful_jobs_history_limit
    job_template {
      metadata {}
      spec {
        backoff_limit              = var.backoff_limit
        ttl_seconds_after_finished = var.ttl_seconds_after_finished
        template {
          metadata {}
          spec {
            service_account_name = var.aws_iam.service_account_name == null ? module.iam_service_account.service_account_name : var.aws_iam.service_account_name
            node_selector = {
              "kubernetes.io/arch" = var.platform_architecture
            }
            container {
              name    = var.task_name
              image   = var.image
              command = var.cmd
              args    = var.args

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

              env {
                name  = "REMOTE_DEV_PREFIX"
                value = var.remote_dev_prefix
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

              // happy configs: add env-level configs first
              env_from {
                secret_ref {
                  name     = "happy-config.${var.app_name}.${var.deployment_stage}"
                  optional = true
                }
              }
              // happy configs: add stack-level configs second so they override env-level configs
              env_from {
                secret_ref {
                  name     = "happy-config.${var.app_name}.${var.deployment_stage}.${var.stack_name}"
                  optional = true
                }
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
                for_each = toset(var.emptydir_volumes)
                content {
                  # TODO FIXME do we want the mount path to be configurable????
                  mount_path = "/var/${volume_mount.value.name}"
                  name       = volume_mount.value.name
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

              resources {
                limits = {
                  cpu    = var.cpu
                  memory = var.memory
                }
                requests = {
                  cpu    = var.cpu_requests
                  memory = var.memory_requests
                }
              }
            }
          }
        }
      }
    }
  }
}
