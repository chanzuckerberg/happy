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

variable "happymeta_" { # tflint-ignore: terraform_unused_declarations
  type        = string
  description = "Happy Path metadata. Ignored by actual terraform."
}

variable "stack_name" {
  type        = string
  description = "Happy Path stack name"
}

variable "happy_config_secret" { # tflint-ignore: terraform_unused_declarations
  type        = string
  description = "Happy Path configuration secret name"
}

variable "deployment_stage" {
  type        = string
  description = "Deployment stage for the app"
}

variable "stack_prefix" {
  type        = string
  description = "Do bucket storage paths and db schemas need to be prefixed with the stack name? (Usually '/{stack_name}' for dev stacks, and '' for staging/prod stacks)"
  default     = ""
}

variable "k8s_namespace" {
  type        = string
  description = "K8S namespace for this stack"
}

variable "services" {
  type = map(object({
    name : string,
    service_type : string, // oneof: EXTERNAL, INTERNAL, PRIVATE
    desired_count : number,
    port : number,
    memory : string,
    cpu : string,
    health_check_path : optional(string, "/"),
    aws_iam_policy_json : optional(string, ""),
    path : optional(string, "/*"),  // Only used for CONTEXT routing
    priority : optional(number, 0), // Only used for CONTEXT routing
    success_codes : optional(string, "200-499"),
    synthetics : optional(bool, false),
    initial_delay_seconds : optional(number, 30),
    period_seconds : optional(number, 3),
    bypasses : optional(map(object({ // Only used for INTERNAL service_type
      paths   = optional(set(string), [])
      methods = optional(set(string), [])
    })))
  }))
  description = "The services you want to deploy as part of this stack."
  validation {
    condition     = alltrue([for k, v in var.services : (v.service_type == "EXTERNAL" || v.service_type == "INTERNAL" || v.service_type == "PRIVATE")])
    error_message = "The service_type argument needs to be 'EXTERNAL', 'INTERNAL', or 'PRIVATE'"
  }
  validation {
    condition     = alltrue([for k, v in var.services : startswith(v.health_check_path, trimsuffix(v.path, "*"))])
    error_message = "The health_check_path should start with the same prefix as the path argument"
  }
}

variable "tasks" {
  type = map(object({
    image : string,
    memory : string,
    cpu : string,
    cmd : set(string),
  }))
  description = "The deletion/migration tasks you want to run when a stack comes up and down."
}

variable "routing_method" {
  type        = string
  description = "Traffic routing method for this stack. Valid options are 'DOMAIN', when every service gets a unique domain name, or a 'CONTEXT' when all services share the same domain name, and routing is done by request path."
  default     = "DOMAIN"

  validation {
    condition     = var.routing_method == "DOMAIN" || var.routing_method == "CONTEXT"
    error_message = "Only DOMAIN and CONTEXT routing methods are supported."
  }
}
variable "additional_env_vars" {
  type        = map(string)
  description = "Additional environment variables to add to the container"
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

variable "create_dashboard" {
  type        = bool
  description = "Create a dashboard for this stack"
  default     = false
}
