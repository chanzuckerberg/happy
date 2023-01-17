terraform {
  experiments = [module_variable_optional_attrs]
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.45"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.16"
    }
    hapi = {
      source  = "github.com/chanzuckerberg/happy"
      version = ">= 0.45"
    }
  }
  required_version = ">= 1.0"
}
