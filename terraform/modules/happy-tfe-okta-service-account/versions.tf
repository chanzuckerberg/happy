terraform {
  required_providers {
    okta = {
      source  = "okta/okta"
      version = "~> 3.41"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.14"

      configuration_aliases = [aws.czi-si]
    }
    jwks = {
      source  = "iwarapter/jwks"
      version = "0.0.3"
    }
  }
  required_version = ">= 1.3"
}
