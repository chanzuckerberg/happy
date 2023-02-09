variable "dns_prefix" {
  type        = string
  description = "Stack-specific prefix for DNS records"
}

variable "zone" {
  type        = string
  description = "Route53 zone name. Trailing . must be OMITTED!"
}

variable "alb_dns" {
  type        = string
  description = "DNS name for the shared ALB"
}

variable "canonical_hosted_zone" {
  type        = string
  description = "Route53 zone for the shared ALB"
}
