variable environment {
  type        = string
  description = "The deployment environment for the app"
}

variable subdomain {
  type        = string
  description = "The hosted zone subdomain to create"
}

variable tags {
  type        = map(string)
  description = "Tags to associate with env resources"
}
