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

variable "stack_name" {
  type        = string
  description = "Please provide the stack name"
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

variable "cloud_env" {
  type = object({
    public_subnets : list(string),
    private_subnets : list(string),
    database_subnets : list(string),
    database_subnet_group : string,
    vpc_id : string,
    vpc_cidr_block : string,
  })
  description = "Typically data.terraform_remote_state.cloud-env.outputs"
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
  type        = list(object({ name : string, value : string }))
  description = "Additional environment variables to add to the task definition"
  default     = []
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

variable "service_type" {
  type        = string
  description = "The type of the service to deploy. Supported types include 'EXTERNAL', 'INTERNAL', and 'PRIVATE'"
}
