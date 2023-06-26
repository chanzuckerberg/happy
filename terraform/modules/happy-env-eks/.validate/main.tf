module "test_validate" {
  source = "../../happy-env-eks"
  eks-cluster = {
    cluster_id              = "test",
    cluster_arn             = "test",
    cluster_endpoint        = "test",
    cluster_ca              = "test",
    cluster_oidc_issuer_url = "test",
    cluster_version         = "test",
    worker_iam_role_name    = "test",
    worker_security_group   = "test",
    oidc_provider_arn       = "test",
  }
  okta_teams   = []
  base_zone_id = "test"
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