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
    schedule                      = "* * * * *"
    suspend                       = true // This cron job is suspended to be used to create jobs on demand
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
            container {
              name    = var.task_name
              image   = var.image
              command = var.cmd

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

              resources {
                limits = {
                  cpu    = var.cpu
                  memory = var.memory
                }
                requests = {
                  cpu    = var.cpu
                  memory = var.memory
                }
              }
            }
          }
        }
      }
    }
  }
}
