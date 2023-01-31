module "happy_okta_app" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-tfe-okta-app?ref=happy-tfe-okta-app-v2.0.0"

  app_name = var.tags.project
  env      = var.tags.env
  teams    = var.okta_teams
}
