terraform {
  required_providers {
    okta = {
      source  = "chanzuckerberg/okta"
      version = "~> 3.10"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.45"
    }
    jwks = {
      source  = "iwarapter/jwks"
      version = "0.0.3"
    }
  }
  required_version = ">= 1.3"
}
