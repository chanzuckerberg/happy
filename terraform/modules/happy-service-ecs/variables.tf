variable "stack_resource_prefix" {
  type        = string
  description = "Prefix for account-level resources"
}

variable "cpu" {
  type        = number
  description = "CPU shares (1cpu=1024) per task"
  default     = 256
}

variable "memory" {
  type        = number
  description = "Memory in megabytes per task"
  default     = 1024
}

variable "vpc" {
  type        = string
  description = "The VPC that the ECS cluster is deployed to"
}

variable "custom_stack_name" {
  type        = string
  description = "Please provide the stack name"
}

variable "remote_dev_prefix" {
  type        = string
  description = "S3 storage path / db schema prefix"
  default     = ""
}

variable "app_name" {
  type        = string
  description = "Please provide the ECS service name"
}

variable "cluster" {
  type        = string
  description = "Please provide the ECS Cluster ID that this service should run on"
}

variable "image" {
  type        = string
  description = "Image name"
}

variable "service_port" {
  type        = number
  description = "What ports does this service run on?"
  default     = 80
}

variable "desired_count" {
  type        = number
  description = "How many instances of this task should we run across our cluster?"
  default     = 2
}

variable "listener" {
  type        = string
  description = "The Application Load Balancer listener to register with"
}

variable "host_match" {
  type        = string
  description = "Host header to match for target rule. Leave empty to match all requests"
}

variable "security_groups" {
  type        = list(string)
  description = "Security groups for ECS tasks"
}

variable "subnets" {
  type        = list(string)
  description = "Subnets for ecs tasks"
}

variable "task_role" {
  type        = object({ arn = string, name = string })
  description = "ARN and name for the role assumed by tasks"
}

variable "deployment_stage" {
  type        = string
  description = "The name of the deployment stage of the Application"
  default     = "dev"
}

variable "chamber_service" {
  type        = string
  description = "The name of the chamber service from which to load env vars"
  default     = ""
}

variable "priority" {
  type        = number
  description = "Listener rule priority number within the given listener"
}

variable "health_check_path" {
  type        = string
  description = "path to use for health checks"
  default     = "/"
}

variable "wait_for_steady_state" {
  type        = bool
  description = "Whether Terraform should block until the service is in a steady state before exiting"
  default     = false
}

variable "execution_role" {
  type        = string
  description = "Execution role to use for fargate tasks - required for fargate services!"
  default     = ""
}

variable "launch_type" {
  type        = string
  description = "Launch type on which to run your service. The valid values are EC2, FARGATE, and EXTERNAL"
  default     = "FARGATE"
}

variable "additional_env_vars" {
  type        = map(string)
  description = "Additional environment variables to add to the task definition"
  default     = {}
}

variable "tags" {
  type = object({
    happy_env : string,
    happy_stack_name : string,
    happy_service_name : string,
    happy_region : string,
    happy_image : string,
    happy_service_type : string,
    happy_last_applied : string,
  })
  description = "The happy conventional tags."
}

variable "datadog_api_key" {
  type        = string
  default     = ""
  description = "DataDog API Key"
}

variable "datadog_agent" {
  type = object({
    registry : optional(string, "public.ecr.aws/datadog/agent"),
    tag : optional(string, "latest"),
    memory : optional(number, 512),
    cpu : optional(number, 256),
    enabled : optional(bool, false),
  })
  default = {
    registry = "public.ecr.aws/datadog/agent"
    tag      = "latest"
    memory   = 512
    cpu      = 256
    enabled  = false
  }
  description = "DataDog agent image to use"
}

variable "fail_fast" {
  type        = bool
  description = "Should containers fail fast if any errors are encountered?"
  default     = false
}