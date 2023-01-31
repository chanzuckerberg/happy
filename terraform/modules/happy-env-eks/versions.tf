terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.45"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.4"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.16"
    }
    opsgenie = {
      source = "opsgenie/opsgenie"
      # [JH]: opsgenie terraform provider is full of bugs, this is the only version right now that doesn't throw stack traces
      # be careful when updating this version
      version = "= 0.6.14"
    }
  }
  required_version = ">= 1.3"
}
