variable "ingress_name" {
  type        = string
  description = "Name of the ingress resource"
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

variable "health_check_path" {
  type        = string
  description = "path to use for health checks"
  default     = "/"
}

variable "k8s_namespace" {
  type        = string
  description = "K8S namespace for this service"
}

variable "certificate_arn" {
  type        = string
  description = "ACM certificate ARN to attach to the load balancer listener"
}


variable "service_type" {
  type        = string
  description = "The type of the service to deploy. Supported types include 'EXTERNAL', 'INTERNAL', and 'PRIVATE'"
}

variable "tags_string" {
  type        = string
  description = "Tags to apply to ingress resource, comma delimited key=value pairs"
  default     = ""
}

variable "routing" {
  type = object({
    method : optional(string, "CONTEXT")
    host_match : string
    group_name : string
    priority : number
    path : optional(string, "/*")
    service_name : string
    service_port : number
    success_codes : optional(string, "200-499")
  })
  description = "Routing configuration for the ingress"
}