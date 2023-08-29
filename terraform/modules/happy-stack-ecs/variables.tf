variable "app_name" {
  type        = string
  description = "The happy application name"
  default     = ""
}

variable "image_tags" {
  type        = map(string)
  description = "Override image tag for each docker image"
  default     = {}
}

variable "image_tag" {
  type        = string
  description = "Please provide a default image tag"
}

variable "priority" {
  type        = number
  description = "Listener rule priority number within the given listener"
}

variable "stack_name" {
  type        = string
  description = "Happy Path stack name"
}

variable "happy_config_secret" {
  type        = string
  description = "Happy Path configuration secret name"
}

variable "deployment_stage" {
  type        = string
  description = "Deployment stage for the app"
}

variable "chamber_service" {
  type        = string
  description = "The name of the chamber service from which to load env vars"
  default     = ""
}

variable "require_okta" {
  type        = bool
  description = "Whether the ALB's should be on private subnets"
  default     = true
}

variable "stack_prefix" {
  type        = string
  description = "Do bucket storage paths and db schemas need to be prefixed with the stack name? (Usually '/{stack_name}' for dev stacks, and '' for staging/prod stacks)"
  default     = ""
}

variable "wait_for_steady_state" {
  type        = bool
  description = "Should terraform block until ECS services reach a steady state?"
  default     = false
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

variable "desired_count" {
  type        = number
  description = "How many instances of this task should we run across our cluster?"
  default     = 2
}

variable "service_port" {
  type        = number
  description = "What ports does this service run on?"
  default     = 80
}

variable "launch_type" {
  type        = string
  description = "Launch type on which to run your service. The valid values are EC2, FARGATE, and EXTERNAL"
  default     = "FARGATE"
}

variable "fail_fast" {
  type        = bool
  description = "Should containers fail fast if any errors are encountered?"
  default     = false
}