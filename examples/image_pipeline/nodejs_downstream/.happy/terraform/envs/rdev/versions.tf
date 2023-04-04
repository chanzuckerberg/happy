terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.45"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.16"
    }
    datadog = {
      source = "datadog/datadog"
      version = ">= 3.20.0"
    }
    happy = {
      source = "chanzuckerberg/happy"
      version = ">= 0.53.5"
    }
  }
  required_version = ">= 1.3"
}

