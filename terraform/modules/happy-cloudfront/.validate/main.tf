module "test_validate" {
  source = "../../happy-cloudfront"

  frontend = {
    domain_name = "example.com"
    zone_id     = "1234567890"
  }
  origin = {
    domain_name = "example.com"
  }
  tags = {}
  providers = {
    aws.useast1 = aws.useast1
  }
}

provider "aws" {
  alias = "useast1"
}