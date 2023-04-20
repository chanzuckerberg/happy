terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.45"

      configuration_aliases = [aws.czi-si]
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.4"
    }
  }
  required_version = ">= 1.3"
}
