variable "scope_name" {
  type        = string
  default     = "service_account"
  description = "The name of the custom scope that allows the service account to authenticate with Client Credentials flow."
}

variable "aws_ssm_paths" {
  type = object({
    client_id     = string
    client_secret = string
    okta_idp_url  = string
    config_uri    = string
  })
  default = {
    client_id     = "oauth2_proxy_client_id"
    client_secret = "oauth2_proxy_client_secret"
    okta_idp_url  = "oauth2_proxy_oidc_issuer_url"
    config_uri    = "oauth2_proxy_config_uri"
  }
  description = "The name of the SSM paths for the client ID, secret, and other values produced from this app."
}

variable "app_name" {
  type        = string
  description = "The name of the happy application"
}

variable "service_name" {
  type        = string
  default     = "happy"
  description = "The component name this service is going to be deployed into"
}

variable "teams" {
  type        = set(string)
  description = "The set of teams to give access to the Okta app"
}
variable "rbac_role_mapping" {
  type = map(list(string))

  default = {}
}

variable "redirect_uris" {
  default = []
  type    = list(string)
}

variable "login_uri" {
  default = ""
  type    = string
}
