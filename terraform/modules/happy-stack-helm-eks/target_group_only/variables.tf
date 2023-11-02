variable "health_check_path" {
  description = "The path to use for the health check"
  type        = string
  default     = "/health"
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

variable "routing" {
  type = object({
    method : optional(string, "DOMAIN")
    host_match : string
    additional_hostnames : optional(set(string), [])
    group_name : string
    alb : optional(object({
      name : string,
      listener_port : number,
    }), null)
    priority : number
    path : optional(string, "/*")
    service_name : string
    port : number
    service_port : number
    alb_idle_timeout : optional(number, 60) // in seconds
    service_scheme : optional(string, "HTTP")
    scheme : optional(string, "HTTP")
    success_codes : optional(string, "200-499")
    service_type : string
    service_mesh : bool
    allow_mesh_services : optional(list(object({
      service : optional(string, null),
      stack : optional(string, null),
      service_account_name : optional(string, null),
    })), null)
    oidc_config : optional(object({
      issuer : string
      authorizationEndpoint : string
      tokenEndpoint : string
      userInfoEndpoint : string
      secretName : string
      }), {
      issuer                = ""
      authorizationEndpoint = ""
      tokenEndpoint         = ""
      userInfoEndpoint      = ""
      secretName            = ""
    })
    bypasses : optional(map(object({
      paths   = optional(set(string), [])
      methods = optional(set(string), [])
    })))
  })
  description = "Routing configuration for the ingress"

  validation {
    condition     = var.routing.service_mesh == true || var.routing.allow_mesh_services == null
    error_message = "The allow_mesh_services option is only supported if service_mesh is enabled on the stack"
  }
}