terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.14"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.4"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.16"
    }
    datadog = {
      source  = "datadog/datadog"
      version = ">= 3.20.0"
    }
  }
  required_version = ">= 1.3"
}
