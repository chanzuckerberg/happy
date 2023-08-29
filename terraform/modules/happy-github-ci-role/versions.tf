terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.14"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.5.1"
    }
  }
  required_version = ">= 1.3"
}
