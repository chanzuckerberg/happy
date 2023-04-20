module "test_validate" {
  source    = "../../happy-env-ecs"
  name      = "test"
  base_zone = "test"
  cloud-env = {
    database_subnet_group = "test"
    database_subnets      = ["test"]
    private_subnets       = ["test"]
    public_subnets        = ["test"]
    vpc_cidr_block        = "test"
    vpc_id                = "test"
  }
  tags = {
    env       = "test"
    managedBy = "teste"
    owner     = "test"
    project   = "test"
    service   = "test"
  }
  providers = {
    aws.czi-si = aws.czi-si
  }
}

provider "aws" {
  alias = "czi-si"
}