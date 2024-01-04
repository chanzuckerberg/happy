variable "app_name" {
  type        = string
  description = "The happy application name"
}

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

variable "args" {
  type        = list(string)
  description = "Args to pass to the command"
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

variable "platform_architecture" {
  type        = string
  description = "Platform architecture"
  default     = "amd64"
}

variable "cpu" {
  type        = string
  description = "CPU shares (1cpu=1000m) per pod"
  default     = "100m"
}

variable "cpu_requests" {
  type        = string
  description = "CPU shares (1cpu=1000m) requested per pod"
  default     = "10m"
}

variable "memory" {
  type        = string
  description = "Memory in megabits per pod"
  default     = "100Mi"
}

variable "memory_requests" {
  type        = string
  description = "Memory requests per pod"
  default     = "10Mi"
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

variable "is_cron_job" {
  type        = bool
  description = "Indicates if this job should be run on a schedule or one-off. If true, set cron_schedule as well"
  default     = false
}

variable "cron_schedule" {
  type        = string
  description = "Cron schedule for this job"
  // default to every year so people have to actually set this to something
  default = "0 0 1 1 *"
}

variable "aws_iam" {
  type = object({
    service_account_name : optional(string, null),
    policy_json : optional(string, ""),
  })
  default     = {}
  description = "The AWS IAM service account or policy JSON to give to the pod. Only one of these should be set."

  validation {
    condition     = var.aws_iam.service_account_name == null || var.aws_iam.policy_json == ""
    error_message = "Only one of service_account_name or policy_json should be set."
  }
}

variable "additional_env_vars" {
  type        = map(string)
  description = "Additional environment variables to add to the task definition"
  default     = {}
}

variable "additional_env_vars_from_config_maps" {
  type = object({
    items : optional(list(string), []),
    prefix : optional(string, ""),
  })
  default = {
    items  = []
    prefix = ""
  }
  description = "Additional environment variables to add to the container from the following config maps"
}

variable "additional_env_vars_from_secrets" {
  type = object({
    items : optional(list(string), []),
    prefix : optional(string, ""),
  })
  default = {
    items  = []
    prefix = ""
  }
  description = "Additional environment variables to add to the container from the following secrets"
}

variable "additional_volumes_from_secrets" {
  type = object({
    items : optional(list(string), []),
    base_dir : optional(string, "/var"),
  })
  default = {
    items    = []
    base_dir = "/var"
  }
  description = "Additional volumes to add to the container from the following secrets"
}

variable "additional_volumes_from_config_maps" {
  type = object({
    items : optional(list(string), []),
  })
  default = {
    items = []
  }
  description = "Additional volumes to add to the container from the following config maps"
}


variable "eks_cluster" {
  type = object({
    cluster_id : string,
    cluster_arn : string,
    cluster_endpoint : string,
    cluster_ca : string,
    cluster_oidc_issuer_url : string,
    cluster_version : string,
    worker_iam_role_name : string,
    worker_security_group : string,
    oidc_provider_arn : string,
  })
  description = "eks-cluster module output"
}


variable "tags" {
  description = "Standard tags to attach to all happy services"
  type = object({
    env : string,
    owner : string,
    project : string,
    service : string,
    managedBy : string,
  })
  default = {
    env       = "ADDTAGS"
    managedBy = "ADDTAGS"
    owner     = "ADDTAGS"
    project   = "ADDTAGS"
    service   = "ADDTAGS"
  }
}
