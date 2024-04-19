variable "env" {
  type        = string
  description = "The deployment environment for the app"
}

variable "subdomain" {
  type        = string
  description = "The hosted zone subdomain to create"
}

variable "tags" {
  type        = map(string)
  description = "Tags to associate with env resources"
}


variable "base_domain" {
  type        = string
  description = "The base domain to create the hosted zone in"
  default     = "si.czi.technology"
}