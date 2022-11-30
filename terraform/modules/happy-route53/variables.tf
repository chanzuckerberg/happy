<<<<<<< HEAD
variable env {
=======
variable environment {
>>>>>>> eb4c92ac1d6268eef2fc5b55eb207728c5dad28d
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
