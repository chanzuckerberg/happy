variable "ingress_name" {
    type        = string
    description = "Name of the ingress resource"
}

variable "host_match" {
  type        = string
  description = "Host header to match for target rule. Leave empty to match all requests"
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


variable "backends" {
  type        = list(object(
    {
      service_name : string
      service_port : number
      path         : string
    }
  ))
  description = "The backends to register with the Application Load Balancer"
  default     = []
}


variable "path" {
  type        = string
  description = "The path to register with the Application Load Balancer"
  default     = "/*"
}

variable "service_name" {
  type        = string
  description = "Service name to be deployed"
}

variable "service_port" {
  type        = number
  description = "What ports does this service run on?"
  default     = 80
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

variable "success_codes" {
  type        = string
  default     = "200-499"
  description = "The range of success codes that are used by the ALB ingress controller."
}