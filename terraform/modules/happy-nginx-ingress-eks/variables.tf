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

