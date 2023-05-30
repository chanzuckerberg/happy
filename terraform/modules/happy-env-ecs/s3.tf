module "s3_bucket" {
  for_each          = var.s3_buckets
  source            = "github.com/chanzuckerberg/cztack//aws-s3-private-bucket?ref=v0.56.2"
  project           = var.tags["project"]
  env               = var.tags["env"]
  service           = "happy"
  owner             = var.tags["owner"]
  bucket_name       = each.value["name"]
  bucket_policy     = ""
  enable_versioning = true
}
