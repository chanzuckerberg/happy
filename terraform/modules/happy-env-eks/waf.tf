# Only create a WAF if we decide to via variables
module "regional-waf" {
  count  = var.include_waf ? 1 : 0
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/web-acl-regional?ref=web-acl-regional-v1.1.0"
  tags   = var.tags
}
