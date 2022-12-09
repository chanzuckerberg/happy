module "happy_apps" {
  for_each = var.envs
  source   = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/okta-app-oauth?ref=v0.202.0"

  okta = {
    label         = "${var.service_name}-${var.app_name}-${each.value}"
    redirect_uris = ["https://oauth.${var.app_name}.${each.value}.si.czi.technology/oauth2/callback"]
    login_uri     = "https://oauth.${var.app_name}.${each.value}.si.czi.technology"
    tenant        = "czi-prod"
  }
  tags = {
    owner   = "infra-eng@chanzuckerberg.com"
    service = "${var.service_name}-oauth"
    project = var.app_name
    env     = each.value
  }
  aws_ssm_paths = var.aws_ssm_paths
}

resource "okta_app_group_assignments" "happy_app" {
    for_each  = module.happy_apps
  app_id    = each.value.app.id
  group_ids = var.teams
}
