
locals {
  target_group_only_services = [for sd in local.service_definitions : sd if sd.service_type == "TARGET_GROUP_ONLY"]
  other_services             = [for sd in local.service_definitions : sd if sd.service_type != "TARGET_GROUP_ONLY"]

  updated_target_service_definitions = [for sd in local.service_definitions : merge(sd, {
    "targetGroupArn" = module.target_group_only.aws_lb_target_group_arn
    "securityGroups" = module.target_group_only.security_groups
  })]

  updated_other_service_definitions = [for sd in local.other_services : merge(sd, {
    "targetGroupArn" = ""
    "securityGroups" = []
  })]

  patched_service_definitions = concat(local.updated_other_service_definitions, local.updated_target_service_definitions)
}

module "target_group_only" {
  for_each          = local.target_group_only_services
  source            = "./target_group_only"
  routing           = each.value.routing
  cloud_env         = local.cloud_env
  health_check_path = each.value.health_check_path
}