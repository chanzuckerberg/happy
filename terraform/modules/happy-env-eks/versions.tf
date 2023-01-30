terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.45"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.4"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.16"
    }
    okta = {
      source  = "chanzuckerberg/okta"
      version = "~> 3.10"
    }
    opsgenie = {
      source  = "opsgenie/opsgenie"
      version = "= 0.6.14"
    }
  }
  required_version = ">= 1.3"
}
