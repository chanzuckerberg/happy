terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.14"

      configuration_aliases = [aws.useast1]
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.16"
    }
    datadog = {
      source  = "datadog/datadog"
      version = ">= 3.20.0"
    }
    validation = {
      source  = "tlkamp/validation"
      version = "1.0.0"
    }
    happy = {
      source  = "chanzuckerberg/happy"
      version = ">= 0.53.5"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.4.3"
    }
  }
  required_version = ">= 1.3"
}
