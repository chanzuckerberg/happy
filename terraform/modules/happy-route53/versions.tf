terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.14"

      configuration_aliases = [aws.czi-si]
    }
  }
  required_version = ">= 1.3"
}
