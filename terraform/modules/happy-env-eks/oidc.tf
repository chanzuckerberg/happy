
module "happy_okta_app" {
  source = "../happy-tfe-okta-app"

  app_name = "${var.tags.project}-${var.tags.env}-${var.tags.service}"
  env     = var.tags.env
  teams    = var.okta_teams
}
