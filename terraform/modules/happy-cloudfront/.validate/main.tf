module "test_validate" {
  source = "../../happy-cloudfront"

  frontend = {
    domain_name = "example.com"
    zone_id     = "1234567890"
  }
  origins = [
    {
      domain_name  = "example1.com"
      path_pattern = "/api/oauth/*"
    },
    {
      domain_name  = "example2.com"
      path_pattern = "/"
    }
  ]
  tags = {
    owner     = ""
    service   = ""
    project   = ""
    env       = ""
    managedBy = ""
  }
  providers = {
    aws.useast1 = aws.useast1
  }
}

provider "aws" {
  alias = "useast1"
}