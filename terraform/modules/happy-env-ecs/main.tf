locals {
  env       = var.tags["env"]
  owner     = var.tags["owner"]
  project   = var.tags["project"]
  component = "happy_env"
  service   = var.tags["service"]
}
