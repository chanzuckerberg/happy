terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.45"
    }
    hapi = {
      source  = "github.com/chanzuckerberg/happy/terraform/provider"
      version = ">= 0.45"
    }
  }
  required_version = ">= 1.0"
}
