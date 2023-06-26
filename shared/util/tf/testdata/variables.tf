variable "aws_account_id" {
  type        = string
  description = "AWS account ID to apply changes to"
}
variable "k8s_cluster_id" {
  type        = string
  description = "EKS K8S Cluster ID"
}
variable "k8s_namespace" {
  type        = string
  description = "K8S namespace for this stack"
}
variable "aws_role" {
  type        = string
  description = "Name of the AWS role to assume to apply changes"
}
variable "image_tag" {
  type        = string
  description = "Please provide an image tag"
}
variable "image_tags" {
  type        = string
  description = "Override the default image tags (json-encoded map)"
  default     = "{}"
}
variable "stack_name" {
  type        = string
  description = "Happy Path stack name"
}