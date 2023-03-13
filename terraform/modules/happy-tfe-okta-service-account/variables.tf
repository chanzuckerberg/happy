variable "tags" {
  description = "Standard tags"
  type = object({
    env : string,
    owner : string,
    project : string,
    service : string,
    managedBy : string,
  })
}

variable "okta_tenant" {
  type        = string
  description = "The Okta tenant to create the authorization server."
  default     = "czi-prod"
}

variable "happy_namespace" {
  description = "Happy Chamber service values"
  type = object({
    env : string,
    project : string,
    service : string,
  })
}
