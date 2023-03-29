module "happy_service_account" {
  source = "../happy-tfe-okta-service-account"
  tags   = var.tags

  providers = {
    aws = aws.czi-si
  }
}
