terraform {
  required_providers {
    okta = {
      source  = "okta/okta"
      version = "~> 3.41"
    }
  }
  required_version = ">= 1.3"
}
