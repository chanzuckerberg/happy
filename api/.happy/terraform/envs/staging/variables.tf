variable aws_account_id {
  type        = string
  description = "AWS account ID to apply changes to"
}

variable aws_role {
  type        = string
  description = "Name of the AWS role to assume to apply changes"
}

variable image_tag {
  type        = string
  description = "Please provide an image tag"
}

variable image_tags {
  type        = string
  description = "Override the default image tags (json-encoded map)"
  default     = "{}"
}

variable priority {
  type        = number
  description = "Listener rule priority number within the given listener"
}

variable stack_name {
  type        = string
  description = "Happy Path stack name"
}

variable happy_config_secret {
  type        = string
  description = "Happy Path configuration secret name"
}

variable wait_for_steady_state {
  type        = bool
  description = "Should terraform block until ECS reaches a steady state?"
  default     = true
}
