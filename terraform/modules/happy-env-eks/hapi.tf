module "happy_service_account" {
  source = "../happy-tfe-okta-service-account"

  app_name    = var.tags.project
  env         = var.tags.env
}
