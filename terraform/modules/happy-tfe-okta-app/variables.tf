variable "teams" {
  type        = set(string)
  description = "The set of teams to give access to the Okta app"
}

variable "app_name" {
  type        = string
  description = "The name of the happy application"
}

variable "env" {
  type        = string
  description = "The environment this happy application supports"
}

variable "service_name" {
  type        = string
  default     = "happy"
  description = "The component name this service is going to be deployed into"
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

variable "redirect_uris" {
  type    = list(string)
  default = []
}

variable "login_uri" {
  type    = string
  default = ""
}

variable "grant_types" {
  type        = list(string)
  default     = ["authorization_code"]
  description = "Additional grant types (authorization_code is offered by default)"
}

variable "app_type" {
  type        = string
  default     = "web"
  description = "The type of OAuth application. Valid values: `web`, `native`, `browser`, `service`. For SPA apps use `browser`."
}

variable "token_endpoint_auth_method" {
  type        = string
  default     = "client_secret_basic"
  description = "Requested authentication method for the token endpoint. It can be set to `none`, `client_secret_post`, `client_secret_basic`, `client_secret_jwt`, `private_key_jwt`. To enable PKCE, set this to `none`."
}

variable "omit_secret" {
  default     = false
  type        = bool
  description = "Whether the provider should persist the application's secret to state. Your app's client_secret will be recreated if this ever changes from true => false."
}

variable "base_domain" {
  default     = "si.czi.technology"
  type        = string
  description = "The base domain to use for all the happy stacks. The default is app_name.env.si.czi.technology"
}
