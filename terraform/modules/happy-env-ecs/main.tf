locals {
  env       = var.tags["env"]
  owner     = var.tags["owner"]
  project   = var.tags["project"]
  component = "happy-env"
  service   = var.tags["service"]
}
