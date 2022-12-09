provider "aws" {

  region = "us-west-2"


  assume_role {
    role_arn = "arn:aws:iam::626314663667:role/tfe-si"
  }

  allowed_account_ids = ["626314663667"]
}

# Aliased Providers (for doing things in every region).
provider "aws" {
  alias  = "czi-si"
  region = "us-west-2"

  assume_role {
    role_arn = "arn:aws:iam::626314663667:role/tfe-si"
  }

  allowed_account_ids = ["626314663667"]
}

terraform {
  required_version = "=1.2.6"
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "4.45.0"
    }
  }
}

module "test" {
  source = "../"

  subdomain = "test"
  env       = "test"
  tags      = {}

  providers = {
    aws.czi-si = aws.czi-si
  }
}
