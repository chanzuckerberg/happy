variable "custom_stack_name" {
  type        = string
  description = "Please provide the stack name"
}

variable "app_name" {
  type        = string
  description = "Please provide the ECS service name"
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


variable "tags" {
  type = object({
    happy_env : string,
    happy_stack_name : string,
    happy_service_name : string,
    happy_region : string,
    happy_image : string,
    happy_service_type : string,
    happy_last_applied : string,
  })
  description = "The happy conventional tags."
}
