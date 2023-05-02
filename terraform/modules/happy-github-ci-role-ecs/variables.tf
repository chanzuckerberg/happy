# variable "ecs_cluster_arn" {
#   description = "The ARN of the ECS cluster that the role should have permissions to"
#   type        = string
# }

variable "happy_app_name" {
  description = "The name of the happy environment"
  type        = string
}

variable "env" {
  description = "The environment this CI role has access to"
  type        = string
}

variable "gh_actions_role_name" {
  description = "The name of the role that was created for the Github Action."
  type        = string
}