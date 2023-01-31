variable "task_name" {
  type        = string
  description = "Happy Path task name"
}

variable "stack_name" {
  type        = string
  description = "Happy Path stack name"
}

variable "image" {
  type        = string
  description = "Image name"
}

variable "cmd" {
  type        = list(string)
  description = "Command to run"
  default     = []
}

variable "remote_dev_prefix" {
  type        = string
  description = "S3 storage path / db schema prefix"
  default     = ""
}

variable "deployment_stage" {
  type        = string
  description = "The name of the deployment stage of the Application"
}

variable "k8s_namespace" {
  type        = string
  description = "K8S namespace for this task"
}

variable "cpu" {
  type        = string
  description = "CPU shares (1cpu=1000m) per pod"
  default     = "100m"
}

variable "memory" {
  type        = string
  description = "Memory in megabits per pod"
  default     = "100Mi"
}

variable "failed_jobs_history_limit" {
  type        = number
  default     = 5
  description = "kubernetes_cron_job failed jobs history limit"
}

variable "starting_deadline_seconds" {
  type        = number
  default     = 30
  description = "kubernetes_cron_job starting_deadline_seconds"
}

variable "successful_jobs_history_limit" {
  type        = number
  default     = 5
  description = "kubernetes_cron_job successful_jobs_history_limit"
}

variable "backoff_limit" {
  type        = number
  default     = 2
  description = "kubernetes_cron_job backoff_limit"
}

variable "ttl_seconds_after_finished" {
  type        = number
  default     = 10
  description = "kubernetes_cron_job ttl_seconds_after_finished"
}
