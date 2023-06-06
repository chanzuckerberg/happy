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
    service : string,
    stack : string
  }))
  description = "A list of service/stack that we want to allow access to this protected service"
}

variable "labels" {
  type        = map(string)
  description = "Labels to apply to Linkerd CRDs"
}