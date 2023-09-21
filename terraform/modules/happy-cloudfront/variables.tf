variable "frontend" {
  type = object({
    domain_name = string
    zone_id     = string
  })
  description = "The domain name and zone ID the user will see."
}

variable "origin" {
  type = object({
    domain_name = string
  })
  description = "The domain name of the origin."
}

variable "cache" {
  type = object({
    min_ttl     = optional(number, 0)
    default_ttl = optional(number, 300)
    max_ttl     = optional(number, 300)
    compress    = optional(bool, true)
  })
  description = "The cache settings for the CloudFront distribution."
  default     = {}
}

variable "tags" {
  type        = map(string)
  description = "Tags to associate with env resources"
}