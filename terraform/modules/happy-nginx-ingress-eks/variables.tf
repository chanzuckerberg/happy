variable "ingress_name" {
  type        = string
  description = "Name of the ingress resource"
}

variable "k8s_namespace" {
  type        = string
  description = "K8S namespace for this service"
}

variable "host_match" {
  type        = string
  description = "Host header value to match when routing to the service"
}

variable "host_path" {
  type        = string
  default     = "/"
  description = "Path value to match when routing to the service"
}


variable "target_service_name" {
  type        = string
  description = "Name of destination service that the ingress should route to"
}

variable "target_service_port" {
  type        = string
  description = "Port of destination service that the ingress should route to"
}

variable "labels" {
  type        = map(string)
  description = "Labels to apply to ingress resource"
}


variable "timeout" {
  type        = number
  description = "Timeout for the ingress resource"
  default     = 60
}

variable "sticky_sessions" {
  type = object({
    enabled          = optional(bool, true),
    duration_seconds = optional(number, 600),
    cookie_name      = optional(string, "happy_sticky_session"),
    cookie_samesite  = optional(string, "Lax"),
  })
  description = "Sticky session configuration"
  default     = {}
}
