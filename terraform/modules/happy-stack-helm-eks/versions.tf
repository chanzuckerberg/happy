terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.23"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.23"
    }
    datadog = {
      source  = "datadog/datadog"
      version = ">= 3.31"
    }
    validation = {
      source  = "tlkamp/validation"
      version = "1.0.0"
    }
    happy = {
      source  = "chanzuckerberg/happy"
      version = ">= 0.108"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.5"
    }
    random = {
      source  = "hashicorp/helm"
      version = ">= 2.11"
    }
  }
  required_version = ">= 1.3"
}
