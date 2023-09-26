module "service" {
  source = "../../happy-service-eks"

  image_tag  = "balh"
  stack_name = "blah"
  cloud_env = {
    database_subnet_group = "test"
    database_subnets      = ["test"]
    private_subnets       = ["test"]
    public_subnets        = ["test"]
    vpc_cidr_block        = "test"
    vpc_id                = "test"
  }
  routing = {
    group_name   = "blah"
    host_match   = "blah"
    port         = 80
    priority     = 1
    service_mesh = false
    service_name = "blah"
    service_port = 80
    service_type = "blah"
  }
  certificate_arn = "certificate_arn"
  k8s_namespace   = "blah"
  eks_cluster = {
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
  container_name = "blah"
  providers = {
    aws.useast1 = aws.useast1
  }
}

provider "aws" {
  alias = "useast1"
}

provider "aws" {

}