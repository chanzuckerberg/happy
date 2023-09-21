module "cloudfront" {
  count  = var.routing.frontend.cloudfront_enabled ? 1 : 0
  source = "../happy-cloudfront"
  frontend = {
    domain_name = var.routing.frontend.domain_name
    zone_id     = var.routing.frontend.zone_id
  }
  origin = {
    domain_name = var.routing.host_match
  }
  tags = var.tags
}