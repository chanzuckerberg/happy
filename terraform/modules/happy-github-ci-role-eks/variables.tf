variable "eks_cluster_arn" {
  description = "The ARN of the EKS cluster that the role should have permissions to"
  type        = string
}

variable "gh_actions_role_name" {
  description = "The name of the role that was created for the Github Action."
  type        = string
}