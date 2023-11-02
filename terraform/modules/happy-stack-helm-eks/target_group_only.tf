
locals {
  target_group_only_services = [for sd in local.service_definitions : sd if sd.service_type == "TARGET_GROUP_ONLY"]
}

module "target_group_only" {
  for_each          = local.target_group_only_services
  source            = "./target_group_only"
  routing           = each.value.routing
  cloud_env         = local.cloud_env
  health_check_path = each.value.health_check_path
}