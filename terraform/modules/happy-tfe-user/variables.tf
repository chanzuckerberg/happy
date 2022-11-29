variable "happy_app_name" {
  type        = string
  description = "The name of the happy application"
}

variable "aws_accounts_can_assume" {
  type        = set(string)
  description = "The set of AWS account names the TFE user should be allowed to assume into"
}
