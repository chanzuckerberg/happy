# Batch with idseq's swipe configurations.
module "swipe" {
  for_each = var.swipe_module_envs
  source   = "git@github.com:chanzuckerberg/swipe?ref=v1.2.1"

  app_name               = each.value.name
  batch_ami_id           = each.value.ami_id
  job_policy_arns        = each.value.job_policy_arns
  extra_env_vars         = each.value.extra_env_vars
  sqs_queues             = each.value.sqs_queues
  workspace_s3_prefixes  = [each.value.workspace_s3_prefix]
  wdl_workflow_s3_prefix = each.value.wdl_workflow_s3_prefix

  network_info = {
    vpc_id           = var.cloud-env.vpc_id
    batch_subnet_ids = var.cloud-env.private_subnets
  }

  spot_min_vcpus           = each.value.spot_min_vcpus
  spot_max_vcpus           = each.value.spot_max_vcpus
  on_demand_min_vcpus      = each.value.on_demand_min_vcpus
  on_demand_max_vcpus      = each.value.on_demand_max_vcpus
  batch_ec2_instance_types = each.value.instance_types
  imdsv2_policy            = "required"
  sfn_template_files       = {}

  tags            = var.tags
  user_data_parts = module.instance-cloud-init-script.parts
}

# Batch with idseq's swipe configurations.
module "batch-swipe" {
  for_each           = var.swipe_envs
  source             = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-batch-env-swipe?ref=v0.227.0"
  cloud-env          = var.cloud-env
  tags               = var.tags
  name               = each.value.name
  ami_id             = each.value.ami_id
  job_policy_arns    = each.value.job_policy_arns
  min_vcpus          = each.value.min_vcpus
  max_vcpus          = each.value.max_vcpus
  spot_desired_vcpus = each.value.spot_desired_vcpus
  ec2_desired_vcpus  = each.value.ec2_desired_vcpus
  instance_type      = each.value.instance_type
  conf_version       = each.value.version
  user_data_parts    = module.instance-cloud-init-script.parts
}

module "instance-cloud-init-script" {
  source = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/instance-cloud-init-script?ref=v0.227.0"

  project = local.project
  owner   = local.owner
  env     = local.env
  service = local.service

  users = var.ssh_users

  base64_encode = "false"
  gzip          = "false"
}

# Batch "classic" -- we're hoping to deprecate this in favor of swipe!
module "batch" {
  for_each        = var.batch_envs
  source          = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-batch-env?ref=aws-batch-env-v0.1.0"
  cloud-env       = var.cloud-env
  tags            = var.tags
  name            = each.value.name
  job_policy_arns = each.value.job_policy_arns
  min_vcpus       = each.value.min_vcpus
  max_vcpus       = each.value.max_vcpus
  desired_vcpus   = each.value.desired_vcpus
  instance_type   = each.value.instance_type
  conf_version    = each.value.version
  volume_size     = each.value.volume_size
  ssh_users       = var.ssh_users
  init_script     = each.value.init_script
}
