output "ecr" {
  value = {
    name = module.ecr.repository_name
    url  = module.ecr.repository_url
    arn  = module.ecr.repository_arn
  }
}

output "target_group_arn" {
  value = length(aws_lb_target_group.this) == 0 ? "" : aws_lb_target_group.this[0].arn
}