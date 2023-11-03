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

variable "enable_service_mesh" {
  type        = bool
  description = "Enable service mesh for this stack"
  default     = false
}

variable "services" {
  type = map(object({
    name : string,
    service_type : optional(string, "INTERNAL"),
    allow_mesh_services : optional(list(object({
      service : optional(string, null),
      stack : optional(string, null),
      service_account_name : optional(string, null)
    })), null),
    ingress_security_groups : optional(list(string), []), // Only used for VPC service_type
    alb : optional(object({
      name : string,
      listener_port : number,
    }), null), // Only used for TARGET_GROUP_ONLY
    desired_count : optional(number, 2),
    max_count : optional(number, 5),
    max_unavailable_count : optional(string, "1"),
    scaling_cpu_threshold_percentage : optional(number, 80),
    port : optional(number, 80),
    scheme : optional(string, "HTTP"),
    cmd : optional(list(string), []),
    args : optional(list(string), []),
    image_pull_policy : optional(string, "IfNotPresent"), // Supported values: IfNotPresent, Always, Never
    tag_mutability : optional(bool, true),
    scan_on_push : optional(bool, false),
    service_port : optional(number, null),
    service_scheme : optional(string, "HTTP"),
    memory : optional(string, "500Mi"),
    memory_requests : optional(string, "200Mi"),
    cpu : optional(string, "1"),
    cpu_requests : optional(string, "500m"),
    gpu : optional(number, null), // Whole number of GPUs to request, 0 will schedule all available GPUs. Requires GPU-enabled nodes in the cluster, `k8s-device-plugin` installed, platform_architecture = "amd64", and additional_node_selectors = { "nvidia.com/gpu.present" = "true" } present.
    health_check_path : optional(string, "/"),
    aws_iam : optional(object({
      policy_json : optional(string, ""),
      service_account_name : optional(string, null),
    }), {}),
    path : optional(string, "/*"),  // Only used for CONTEXT and TARGET_GROUP_ONLY routing
    priority : optional(number, 0), // Only used for CONTEXT and TARGET_GROUP_ONLY routing
    success_codes : optional(string, "200-499"),
    synthetics : optional(bool, false),
    initial_delay_seconds : optional(number, 30),
    alb_idle_timeout : optional(number, 60) // in seconds
    period_seconds : optional(number, 3),
    platform_architecture : optional(string, "amd64"),     // Supported values: amd64, arm64; GPU nodes are amd64 only.
    additional_node_selectors : optional(map(string), {}), // For GPU use: { "nvidia.com/gpu.present" = "true" }
    bypasses : optional(map(object({                       // Only used for INTERNAL service_type
      paths   = optional(set(string), [])
      methods = optional(set(string), [])
    })), {})
    sidecars : optional(map(object({
      image : string
      tag : string
      port : optional(number, 80),
      scheme : optional(string, "HTTP"),
      memory : optional(string, "200Mi")
      cpu : optional(string, "500m")
      image_pull_policy : optional(string, "IfNotPresent") // Supported values: IfNotPresent, Always, Never
      health_check_path : optional(string, "/")
      initial_delay_seconds : optional(number, 30),
      period_seconds : optional(number, 3),
    })), {})
  }))
  description = "The services you want to deploy as part of this stack."

  validation {
    condition = alltrue([for k, v in var.services : (
      v.scheme == "HTTP" ||
      v.scheme == "HTTPS"
    )])
    error_message = "The scheme argument needs to be 'HTTP' or 'HTTPS'."
  }

  validation {
    condition = alltrue([for k, v in var.services : (
      v.image_pull_policy == "IfNotPresent" ||
      v.image_pull_policy == "Always" ||
      v.image_pull_policy == "Never"
    )])
    error_message = "The image_pull_policy argument needs to be 'IfNotPresent', 'Always', or 'Never'."
  }

  validation {
    condition = alltrue([for k, v in var.services : (
      v.service_scheme == "HTTP" ||
      v.service_scheme == "HTTPS"
    )])
    error_message = "The service_scheme argument needs to be 'HTTP' or 'HTTPS'."
  }

  validation {
    condition = alltrue([for k, v in var.services : (
      v.service_type == "EXTERNAL" ||
      v.service_type == "INTERNAL" ||
      v.service_type == "PRIVATE" ||
      v.service_type == "IMAGE_TEMPLATE" ||
      v.service_type == "TARGET_GROUP_ONLY" ||
      v.service_type == "VPC"
    )])
    error_message = "The service_type argument needs to be 'EXTERNAL', 'INTERNAL', 'PRIVATE', 'VPC', or 'IMAGE_TEMPLATE'."
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
  validation {
    condition     = alltrue([for service in var.services : alltrue([for sidecar in service.sidecars : contains(["IfNotPresent", "Always", "Never"], sidecar.image_pull_policy)])])
    error_message = "Value of a sidecar image_pull_policy needs to be 'IfNotPresent', 'Always', or 'Never'."
  }
  validation {
    condition     = alltrue([for service in var.services : alltrue([for sidecar in service.sidecars : length(sidecar.health_check_path) > 0])])
    error_message = "Value of a sidecar health_check_path must be a non-empty string."
  }
  validation {
    condition     = alltrue([for service in var.services : alltrue([for sidecar in service.sidecars : sidecar.initial_delay_seconds > 0])])
    error_message = "Value of a sidecar initial_delay_seconds must be a positive number."
  }
  validation {
    condition     = alltrue([for service in var.services : alltrue([for sidecar in service.sidecars : sidecar.period_seconds > 0])])
    error_message = "Value of a sidecar period_seconds must be a positive number."
  }
}

variable "tasks" {
  type = map(object({
    image : string,
    memory : optional(string, "200Mi"),
    cpu : optional(string, "500m"),
    cmd : optional(list(string), []),
    args : optional(list(string), []),
    platform_architecture : optional(string, "amd64"), // Supported values: amd64, arm64
    is_cron_job : optional(bool, false),
    aws_iam : optional(object({
      policy_json : optional(string, ""),
      service_account_name : optional(string, null),
    }), {}),
    cron_schedule : optional(string, "0 0 1 1 *"),
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

variable "create_dashboard" {
  type        = bool
  description = "Create a dashboard for this stack"
  default     = false
}

variable "additional_pod_labels" {
  type        = map(string)
  description = "Additional labels to add to the pods."
  default     = {}
}

variable "skip_config_injection" {
  type        = bool
  description = "Skip injecting app configs into the services / tasks"
  default     = false
}

variable "features" {
  type = map(object({
    name : string,
    enabled : bool,
    config : optional(map(string), {}),
  }))
  description = "The features you want to enable for this stack."
  default     = {}
}