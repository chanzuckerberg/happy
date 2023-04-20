module "test_validate" {
  source = "../../happy-tfe-okta-service-account"

  tags = {
    env       = "test"
    managedBy = "teste"
    owner     = "test"
    project   = "test"
    service   = "test"
  }
  providers = {
    aws.czi-si = aws.czi-si
  }
}

provider "aws" {
  alias = "czi-si"
}