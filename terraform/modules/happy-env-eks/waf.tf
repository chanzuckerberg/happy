module "regional-waf" {
  count  = var.include_waf
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/web-acl-regional?ref=web-acl-regional-v1.1.0"
  tags   = var.tags
}
