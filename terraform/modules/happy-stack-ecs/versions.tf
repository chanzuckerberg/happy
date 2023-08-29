terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.14"
    }
    datadog = {
      source  = "datadog/datadog"
      version = ">= 3.20.0"
    }
    happy = {
      source  = "chanzuckerberg/happy"
      version = ">= 0.97.1"
    }
  }
  required_version = ">= 1.3"
}
