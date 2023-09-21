variable "source" {
  type = object({
    domain_name = string
    zone_id     = string
  })
  description = "The domain name and zone id of the domain to be used as the source for the redirect."
}

variable "origin" {

}