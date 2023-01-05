variable "aws_account_id" {
  type        = string
  description = "AWS account ID to apply changes to"
  default     = ""
  sensitive = false
}

variable "image_tags" {
  type        = map(string)
  description = "Override image tag for each docker image"
  default     = {}
  sensitive = false
}

variable "image_tag" {
  type        = string
  description = "Please provide a default image tag"
  sensitive = false
}

variable "happymeta_" {
  type        = string
  description = "Happy Path metadata. Ignored by actual terraform."
  sensitive = false
}

variable "stack_name" {
  type        = string
  description = "Happy Path stack name"
  sensitive = false
}

variable "happy_config_secret" {
  type        = string
  description = "Happy Path configuration secret name"
  sensitive = false
}

variable "deployment_stage" {
  type        = string
  description = "Deployment stage for the app"
  sensitive = false
}

variable "backend_url" {
  type        = string
  description = "For non-proxied stacks, send in the canonical front/backend URL's"
  default     = ""
  sensitive = false
}

variable "frontend_url" {
  type        = string
  description = "For non-proxied stacks, send in the canonical front/backend URL's"
  default     = ""
  sensitive = false
}

variable "stack_prefix" {
  type        = string
  description = "Do bucket storage paths and db schemas need to be prefixed with the stack name? (Usually '/{stack_name}' for dev stacks, and '' for staging/prod stacks)"
  default     = ""
  sensitive = false
}

variable "k8s_namespace" {
  type        = string
  description = "K8S namespace for this stack"
  sensitive = false
}

variable "services" {
  type = map(object({
    name : string,
    service_type : string,
    desired_count : number,
    port : number,
    memory : string,
    cpu : string,
    health_check_path : optional(string, "/"),
    aws_iam_policy_json : optional(string, ""),
  }))
  description = "The services you want to deploy as part of this stack."
  sensitive = false
}

variable "tasks" {
  type = map(object({
    image : string,
    memory : string,
    cpu : string,
    cmd : set(string),
  }))
  description = "The deletion/migration tasks you want to run when a stack comes up and down."
  sensitive = false
}
