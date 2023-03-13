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
