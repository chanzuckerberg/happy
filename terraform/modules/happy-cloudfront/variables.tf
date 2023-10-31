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

variable "geo_restriction_locations" {
  type        = set(string)
  default     = ["US"]
  description = "The countries to whitelist for the CloudFront distribution."
}

variable "cache_allowed_methods" {
  type        = set(string)
  default     = ["GET", "HEAD"]
  description = "The allowed cache methods for the CloudFront distribution."
}

variable "allowed_methods" {
  type        = set(string)
  default     = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
  description = "The allowed methods for the CloudFront distribution."
}

variable "origin_request_policy_id" {
  type = string
  #  managed_caching_disabled_policy_id : https://us-east-1.console.aws.amazon.com/cloudfront/v3/home?region=us-west-2#/policies/cache/b689b0a8-53d0-40ab-baf2-68738e2966ac
  default     = "b689b0a8-53d0-40ab-baf2-68738e2966ac"
  description = "The origin request policy ID for the CloudFront distribution."
}

variable "cache_policy_id" {
  type = string
  #  managed_all_viewer_except_host_policy_id : https://us-east-1.console.aws.amazon.com/cloudfront/v3/home?region=us-west-2#/policies/origin/4135ea2d-6df8-44a3-9df3-4b5a84be39ad
  default     = "4135ea2d-6df8-44a3-9df3-4b5a84be39ad"
  description = "The cache policy ID for the CloudFront distribution."
}

variable "price_class" {
  type        = string
  default     = "PriceClass_100"
  description = "The price class for the CloudFront distribution."

  validation {
    condition     = contains(["PriceClass_100", "PriceClass_200", "PriceClass_All"], var.price_class)
    error_message = "Price class must be one of PriceClass_100, PriceClass_200, or PriceClass_All."
  }
}

variable "tags" {
  type        = map(string)
  description = "Tags to associate with env resources"
}