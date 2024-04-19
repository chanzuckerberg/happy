

variable "service_name" {
  type        = string
  description = "Service name"
}

variable "synthetic_url" {
  type        = string
  description = "URL to run synthetic tests against"
}

variable "stack_name" {
  type        = string
  description = "Stack name"
}

variable "deployment_stage" {
  type        = string
  description = "Deployment stage"
}

variable "opsgenie_owner" {
  type        = string
  description = "Opsgenie Owner"
}

variable "tags" {
  type        = list(string)
  description = "tags"
}
