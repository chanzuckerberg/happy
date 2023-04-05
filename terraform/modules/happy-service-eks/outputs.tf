output "ecr" {
  value = {
    name = module.ecr.repository_name
    url  = module.ecr.repository_url
    arn  = module.ecr.repository_arn
  }
}

# output "target_groups" {
#   value = aws_lb_target_group.this
# }