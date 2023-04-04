module "stack" {
  source = "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-eks?ref=main"

  image_tag        = var.image_tag
  image_tags       = jsondecode(var.image_tags)
  stack_name       = var.stack_name
  deployment_stage = "staging"
  stack_prefix     = "/${var.stack_name}"
  k8s_namespace    = var.k8s_namespace

  # don't deploy any services for your base image
  services = {
    base = {
      name         = "base"
      service_type = "IMAGE_TEMPLATE"
    }
  }
  tasks = {
  }
}
