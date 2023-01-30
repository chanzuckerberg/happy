terraform {
  required_providers {
    okta = {
      source  = "okta/okta"
      version = "~> 3.41"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.45"
    }
  }
  required_version = ">= 1.3"
}
