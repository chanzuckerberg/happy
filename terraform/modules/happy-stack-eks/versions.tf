terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.14"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.16"
    }
    validation = {
      source  = "tlkamp/validation"
      version = "1.0.0"
    }
    happy = {
      source  = "chanzuckerberg/happy"
      version = ">= 0.108.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.4.3"
    }
  }
  required_version = ">= 1.3"
}
