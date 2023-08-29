provider "aws" {
}

terraform {
  required_version = ">=1.3"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.14.0"
    }
  }
}

module "test" {
  source = "../"

  stack_resource_prefix = "stack_resource_prefix"
  execution_role        = "ecs_execution_role"
  memory                = 1024
  cpu                   = 512
  custom_stack_name     = "my-stack"
  app_name              = "testapp"
  vpc                   = "test_vpc_id"
  image                 = "latest"
  cluster               = "test-cluster"
  desired_count         = 2
  listener              = "test-listener-arn"
  subnets               = ["test-subnet"]
  security_groups       = ["test-security-group"]
  task_role             = { arn : "test_ecs_role_arn", name : "test_ecs_role_name" }
  service_port          = 3001
  deployment_stage      = "rdev"
  host_match            = "my-stack-test.test.rdev.si.czi.technology"
  priority              = 1
  wait_for_steady_state = true
  launch_type           = "FARGATE"
  additional_env_vars   = { db_host = "some-url", db_password = "pa$$w0rd" }
  chamber_service       = "happy-rdev-testapp"
  tags = {
    happy_env : "string",
    happy_stack_name : "string",
    happy_service_name : "string",
    happy_region : "string",
    happy_image : "string",
    happy_service_type : "string",
    happy_last_applied : "string",
  }
  datadog_api_key = "dd_api_key"
}
