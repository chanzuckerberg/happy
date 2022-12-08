terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.45"

      configuration_aliases = [aws.czi-si]
    }
  }
  required_version = ">= 1.0"
}
