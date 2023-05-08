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

variable "stack_name" {
  type        = string
  description = "Happy Path stack name"
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
    service_type : optional(string, "INTERNAL"),
    alb : optional(object({
      name : string,
      listener_port : number,
    }), null), // Only used for TARGET_GROUP_ONLY
    desired_count : optional(number, 2),
    max_count : optional(number, 2),
    scaling_cpu_threshold_percentage : optional(number, 80),
    port : optional(number, 80),
    memory : optional(string, "100Mi"),
    cpu : optional(string, "100m"),
    health_check_path : optional(string, "/"),
    aws_iam_policy_json : optional(string, ""),
    path : optional(string, "/*"),  // Only used for CONTEXT and TARGET_GROUP_ONLY routing
    priority : optional(number, 0), // Only used for CONTEXT and TARGET_GROUP_ONLY routing
    success_codes : optional(string, "200-499"),
    synthetics : optional(bool, false),
    initial_delay_seconds : optional(number, 30),
    period_seconds : optional(number, 3),
    platform_architecture : optional(string, "amd64"), // Supported values: amd64, arm64
    bypasses : optional(map(object({                   // Only used for INTERNAL service_type
      paths   = optional(set(string), [])
      methods = optional(set(string), [])
    })), {})
    sidecars : optional(map(object({
      name : string
      image : string
      tag : string
      port : optional(number, 80),
      memory : optional(string, "100Mi")
      cpu : optional(string, "100m")
    })), [])
  }))
  description = "The services you want to deploy as part of this stack."
  validation {
    condition = alltrue([for k, v in var.services : (
      v.service_type == "EXTERNAL" ||
      v.service_type == "INTERNAL" ||
      v.service_type == "PRIVATE" ||
      v.service_type == "IMAGE_TEMPLATE" ||
      v.service_type == "TARGET_GROUP_ONLY"
    )])
    error_message = "The service_type argument needs to be 'EXTERNAL', 'INTERNAL', 'PRIVATE', or 'IMAGE_TEMPLATE'."
  }
  validation {
    condition     = alltrue([for k, v in var.services : v.alb != null if v.service_type == "TARGET_GROUP_ONLY"])
    error_message = "The service_type 'TARGET_GROUP_ONLY' requires an alb"
  }
  validation {
    # The health check prefix needs to contain the service path for CONTEXT services, but not TARGET_GROUP_ONLY services.
    condition     = alltrue([for k, v in var.services : startswith(v.health_check_path, trimsuffix(v.path, "*")) if v.service_type != "TARGET_GROUP_ONLY"])
    error_message = "The health_check_path should start with the same prefix as the path argument."
  }
  validation {
    condition     = alltrue(flatten([for k, v in var.services : [for path in flatten([for x, y in v.bypasses : y.paths]) : startswith(path, trimsuffix(v.path, "*"))]]))
    error_message = "The bypasses.paths should all start with the same prefix as the path argument."
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
  default     = {}
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
