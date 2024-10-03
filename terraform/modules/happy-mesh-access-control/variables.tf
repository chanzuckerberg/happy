variable "k8s_namespace" {
  type        = string
  description = "K8S namespace for this service being protected"
}

variable "service_port" {
  type        = number
  description = "Port of the service being protected"
}

variable "service_name" {
  type        = string
  description = "Name of the service being protected"
}

variable "service_type" {
  type        = string
  description = "Type of the service being protected"
}

variable "deployment_stage" {
  type        = string
  description = "The name of the deployment stage of the Application"
}

variable "allow_mesh_services" {
  type = list(object({
    service : optional(string, null),
    stack : optional(string, null),
    service_account_name : optional(string, null),
  }))
  description = "A list of service/stack that we want to allow access to this protected service"
}

variable "allow_k6_operator" {
  type        = bool
  description = "A flag to allow the grafana k6 operator to access this protected service"
  default     = false
}

variable "labels" {
  type        = map(string)
  description = "Labels to apply to Linkerd CRDs"
}
