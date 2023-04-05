output "name" {
  value = module.ecr.repository_name
}

output "url" {
  value = module.ecr.repository_url
}

output "arn" {
  value = module.ecr.repository_arn
}

output "target_group_arn" {
  value = length(aws_lb_target_group.this) == 0 ? "" : aws_lb_target_group.this[0].arn
}