provider "aws" {
}

terraform {
  required_version = ">=1.3"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.14.0"
    }
  }
}

module "test" {
  source = "../"

  app_name              = "testapp"
  happy_config_secret   = "happy_config_secret"
  image_tag             = "latest"
  priority              = 1
  stack_name            = "my-stack"
  deployment_stage      = "rdev"
  require_okta          = true
  stack_prefix          = "/my-stack"
  wait_for_steady_state = true
  chamber_service       = "happy-rdev-testapp"
  service_port          = 3001
}
